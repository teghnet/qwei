package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/defiweb/go-offchain/bn"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/abi"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/teghnet/qwei/vars"
)

func main() {
	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, ethereum.GoerliChainID)

	sourceAccount, err := ethereum.NewPrivateKeyAccount(vars.EthKeys[4])
	if err != nil {
		log.Println(err)
		return
	}
	ctx = ethereum.WithAccount(ctx, sourceAccount)

	c := provider.NewAlchemy(vars.AlchemyKeys)
	defer func() {
		err := c.Close(ctx)
		if err != nil {
			log.Printf("error closing client: %s", err)
		}
	}()

	f, err := abi.FromStrings("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F", "count()(uint256)")
	if err != nil {
		log.Println(err)
		return
	}
	u, err := c.Read(ctx, f.MustCall(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%#v", bn.IntFromBigInt(u.(*big.Int)).Int64())
	v, ok := u.([]interface{})
	if !ok {
		log.Printf("%#v", u)
		return
	}
	fmt.Printf("%#v", bn.IntFromBigInt(v[0].(*big.Int)).Int64())
}
