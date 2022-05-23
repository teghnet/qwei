package main

import (
	"context"
	"log"

	"github.com/defiweb/go-offchain/catalog/create2"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/defiweb/go-offchain/examples"
	"github.com/defiweb/go-offchain/examples/utils"
	"github.com/defiweb/go-offchain/infura"
)

func main() {
	chainID := ethereum.GoerliChainID
	account := utils.MustAccount(ethereum.NewPrivateKeyAccount(examples.AccountPrivKey))

	ctx := context.Background()
	ctx = ethereum.WithAccount(ctx, account)
	ctx = ethereum.WithChainID(ctx, chainID)

	address, err := account.Address(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("FROM:", etherscan.Address(chainID, address))
	log.Println("-----")
	client := infura.NewClient(examples.InfuraAPIKey)

	var txParams ethereum.TXParams
	// txParams ,_:= ethereum.NewAutoTXParams(client, 0, 0, 0, ethereum.NewStaticFee(
	// 	ethereum.GWei(bn.FloatFromUint64(1)),
	// 	ethereum.GWei(bn.FloatFromUint64(50)),
	// ))

	gas := uint64(0)

	gas += utils.Run(true, false, false)(create2.Constructor())(ctx, client, txParams)

	utils.PrintGas("total:", gas, ctx, txParams)
}
