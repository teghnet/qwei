package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/catalog/multicall"
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

	if err := do(ctx, client, nil); err != nil {
		log.Println(err)
		return
	}
}

func countItems(addr common.Address) (abi.CallableUnpacker, error) {
	fn, err := abi.Parse("count()(uint256)")
	if err != nil {
		return nil, err
	}
	c, err := fn.Attach(addr)()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getItems(addr common.Address, count int) ([]ethereum.Callable, error) {
	var s string
	var a common.Address
	function, err := abi.Parse("get(uint256)(bytes32,address)")
	if err != nil {
		return nil, err
	}
	method := function.Attach(addr)
	var cs []ethereum.Callable
	for i := 0; i < count; i++ {
		c, err := method(i)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c.Bind(&s, &a))
	}
	return cs, nil
}

func do(ctx context.Context, client ethereum.Client, params ethereum.TXParams) error {
	var count int
	addr := common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F")
	{
		c, err := countItems(addr)
		if err != nil {
			return err
		}
		_, err = client.Read(ctx, c.Bind(&count), params)
		if err != nil {
			return err
		}
	}
	fmt.Println(count)

	cs, err := getItems(addr, count)
	if err != nil {
		return err
	}

	{
		read, err := multicall.Contract().
			TryAggregate(false, cs...).
			Read(ctx, client, params)
		if err != nil {
			return err
		}
		for _, x := range read {
			fmt.Println(x.Success)
			u, ok := x.Result.([]any)
			if !ok {
				fmt.Println(x.Result)
				return errors.New("not array")
			}
			fmt.Println(u)
		}
	}
	return nil
}
