package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/bn"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/teghnet/qwei/vars"
)

func do(ctx context.Context, client ethereum.Client, params ethereum.TXParams) error {
	address := common.HexToAddress("0x64DE91F5A373Cd4c28de3600cB34C7C6cE410C85")
	fmt.Println("  TO", etherscan.Address(ethereum.ChainIDFromContext(ctx), address))
	fmt.Println(address.String())

	hash := common.BytesToHash(common.FromHex("0x1"))
	fmt.Println("@", hash.String())

	storage, err := client.Storage(ctx, address, hash)
	if err != nil {
		return err
	}

	fmt.Println("=", bn.IntFromBytes(storage[16:32]).String())
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
