package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/callable"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/teghnet/qwei/vars"
)

func do(ctx context.Context, client ethereum.Client, params ethereum.TXParams) error {
	addr := common.HexToAddress("0x1E19CF2D73a72Ef1332C882F20534B6519Be0276")
	fmt.Println("  TO", etherscan.Address(ethereum.ChainIDFromContext(ctx), addr))
	fmt.Println(addr.String())

	m, err := callable.ParseMethod(addr, "getPriceRateCache(address)(uint256,uint256,uint256)")
	if err != nil {
		return err
	}
	c, err := m(common.HexToAddress("0xae78736Cd615f374D3085123A210448E74Fc6393"))
	if err != nil {
		return err
	}
	var rate, dur, exp int
	_, err = client.Read(ctx, c.Bind(&rate, &dur, &exp), params)
	if err != nil {
		return err
	}
	fmt.Println(rate, dur, exp)
	return nil
}

func main() {
	account, err := ethereum.NewPrivateKeyAccount(vars.DefaultEthKey)
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, ethereum.MainnetChainID)
	ctx = ethereum.WithAccount(ctx, account)

	fmt.Println("FROM", etherscan.Account(ctx, account))

	client := provider.NewAlchemy(vars.AlchemyKeys)
	defer func() {
		if err := client.Close(ctx); err != nil {
			log.Printf("error closing client: %s", err)
		}
	}()

	if err := do(ctx, client, nil); err != nil {
		log.Println(err)
	}
}
