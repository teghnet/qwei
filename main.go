package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/bn"
	"github.com/defiweb/go-offchain/catalog/maker/changelog"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/abi"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/teghnet/qwei/vars"
)

func main() {
	sourceAccount, err := ethereum.NewPrivateKeyAccount(vars.EthKeys[4])
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, ethereum.GoerliChainID)
	ctx = ethereum.WithAccount(ctx, sourceAccount)

	client := provider.NewAlchemy(vars.AlchemyKeys)
	defer func() {
		if err := client.Close(ctx); err != nil {
			log.Printf("error closing client: %s", err)
		}
	}()

	{
		f, err := abi.MethodCall(
			"0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F",
			"count()(uint256)",
		)()
		if err != nil {
			log.Println(err)
			return
		}

		u, err := client.Read(ctx, f, nil)
		if err != nil {
			log.Println(err)
			return
		}

		v, ok := u.([]interface{})
		if !ok {
			log.Printf("%#v", u)
			return
		}
		fmt.Println(bn.IntFromBigInt(v[0].(*big.Int)).Int64())
	}

	{
		fmt.Println(changelog.Get(1).Read(ctx, client, nil))
	}

	{
		f, err := abi.MethodCall(
			"0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F",
			"get(uint256)(bytes32,address)",
		)(big.NewInt(1))
		if err != nil {
			log.Println(err)
			return
		}

		u, err := client.Read(ctx, f, nil)
		if err != nil {
			log.Println(err)
			return
		}

		v, ok := u.([]interface{})
		if !ok {
			log.Printf("non array %#v", u)
			return
		}
		b := v[0].([32]byte)
		fmt.Println(strings.Trim(string(b[:]), "\x00"), v[1].(common.Address).Hex())
	}
}
