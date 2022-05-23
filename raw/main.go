package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/defiweb/go-offchain/examples"
	"github.com/defiweb/go-offchain/examples/utils"
)

func main() {
	chainID := ethereum.GoerliChainID

	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, chainID)

	account := utils.MustAccount(ethereum.NewPrivateKeyAccount(examples.AccountPrivKey))
	address, err := account.Address(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("from:", etherscan.Address(chainID, address))
	ctx = ethereum.WithAccount(ctx, account)

	client := provider.NewInfura(examples.InfuraAPIKey)

	s := // "6080 6040 52" + // PUSH1 0x80 (free memory pointer) PUSH1 0x40 MSTORE (if uncommented change 61000b to 610010 - JUMPDEST)
		"34 80 15 61000b 57 6000 80 fd 5b 50" + // make call NOT payable
			"32 610100 52" + // ORIGIN
			"33 610120 52" + // CALLER
			"5A 610140 52" + // GAS
			"36 610160 52" + // CALLDATASIZE
			"34 610180 52" + // CALLVALUE
			"610100 51 6000 55" + // SSTORE
			"610120 51 6020 55" +
			"610140 51 6040 55" +
			"610160 51 6060 55" +
			"610180 51 6080 55" +
			"6000 54 610200 52" + // SLOAD
			"6020 54 610220 52" +
			"6040 54 610240 52" +
			"6060 54 610260 52" +
			"6080 54 610280 52" +
			"610300 6000 F3" + // RETURN
			"00" // STOP
	callData, err := ethereum.NewCallDataFromHex(strings.ReplaceAll(s, " ", ""))
	if err != nil {
		panic(err)
	}

	call := ethereum.NewConstructorCall(callData)
	gas, err := client.EstimateGas(ctx, call)
	if err != nil {
		panic(err)
	}
	fmt.Println(gas)

	// r, err := client.Read(ctx, call, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// var x int
	// const wrd = 32
	// for x = 0; x < (len(r)-len(r)%wrd)/wrd; x++ {
	// 	fmt.Println(hexutil.Encode(r[wrd*x:wrd*x+wrd]), x)
	// }
	// if len(r) > x*wrd {
	// 	fmt.Println(hexutil.Encode(r[wrd*x:]))
	// }
}
