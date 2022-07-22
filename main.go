package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/bn"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/abi"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/teghnet/qwei/vars"
)

func Client(acc ethereum.Account) (context.Context, ethereum.Client, func()) {
	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, ethereum.GoerliChainID)
	ctx = ethereum.WithAccount(ctx, acc)

	client := provider.NewAlchemy(vars.AlchemyKeys)

	return ctx, provider.NewAlchemy(vars.AlchemyKeys), func() {
		if err := client.Close(ctx); err != nil {
			log.Printf("error closing client: %s", err)
		}
	}
}

func main() {
	acc, err := ethereum.NewPrivateKeyAccount(vars.EthKeys[4])
	if err != nil {
		log.Println(err)
		return
	}
	ctx, client, closeFn := Client(acc)
	defer closeFn()

	i := new(bn.Int)
	{
		fn, err := abi.Parse("count()(uint256)")
		if err != nil {
			log.Println(err)
			return
		}

		mt := fn.Attach(common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"))

		c, err := mt()
		if err != nil {
			log.Println(err)
			return
		}

		if err := abi.ReadVars(ctx, client, nil)(c, i); err != nil {
			log.Println(err)
			return
		}
	}

	var s string
	var a common.Address

	{
		fn, err := abi.Parse("get(uint256)(bytes32,address)")
		if err != nil {
			log.Println(err)
			return
		}

		mt := fn.Attach(common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"))

		c, err := mt(i.Sub(bn.IntFromInt64(1)).BigInt())
		if err != nil {
			log.Println(err)
			return
		}

		if err := abi.ReadVars(ctx, client, nil)(c, &s, &a); err != nil {
			log.Println(err)
			return
		}
	}
	fmt.Println(s, a)
}
