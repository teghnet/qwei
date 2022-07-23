package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/catalog/multicall"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/abi"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/teghnet/qwei/vars"
)

func main() {
	account, err := ethereum.NewPrivateKeyAccount(vars.EthKeys[4])
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, ethereum.MainnetChainID)
	ctx = ethereum.WithAccount(ctx, account)

	client := provider.NewAlchemy(vars.AlchemyKeys)
	defer func() {
		if err := client.Close(ctx); err != nil {
			log.Printf("error closing client: %s", err)
		}
	}()

	if err := do(ctx, client, nil); err != nil {
		log.Fatalln(err)
	}
}

func do(ctx context.Context, client ethereum.Client, params ethereum.TXParams) error {
	addr := abi.Receiver(common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"))
	var count int
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
	var version string
	{
		c, err := getVersion(addr)
		if err != nil {
			return err
		}
		_, err = client.Read(ctx, c.Bind(&version), params)
		if err != nil {
			return err
		}
	}
	fmt.Println(version)
	var sum string
	{
		c, err := getSum(addr)
		if err != nil {
			return err
		}
		_, err = client.Read(ctx, c.Bind(&sum), params)
		if err != nil {
			return err
		}
	}
	fmt.Println(sum)
	var ipfs string
	{
		c, err := getIPFS(addr)
		if err != nil {
			return err
		}
		_, err = client.Read(ctx, c.Bind(&ipfs), params)
		if err != nil {
			return err
		}
	}
	fmt.Println(ipfs)
	{
		cs, err := getItems(addr, count)
		if err != nil {
			return err
		}

		read, err := multicall.Contract().
			TryAggregate(false, cs...).
			Read(ctx, client, params)
		if err != nil {
			return err
		}

		for i, x := range read {
			u, ok := x.Result.([]any)
			if !ok {
				fmt.Println(x.Result, i)
				continue
			}
			fmt.Println(u[1], i, u[0])
		}
	}
	return nil
}

func getItems(addr abi.Receiver, count int) ([]ethereum.Callable, error) {
	var s string
	method, err := addr.Parse("get(uint256)(bytes32,address)")
	// TODO: we could try creating a set of default mappers
	// method, err := addr.Parse("get(uint256)(bytes32=>string,address)")
	if err != nil {
		return nil, err
	}
	var cs []ethereum.Callable
	for i := 0; i < count; i++ {
		c, err := method(i)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c.Bind(&s))
	}
	return cs, nil
}
func getIPFS(addr abi.Receiver) (abi.CallableUnpacker, error) {
	c, err := addr.Parse("ipfs()(string)")
	if err != nil {
		return nil, err
	}
	return c()
}
func getVersion(addr abi.Receiver) (abi.CallableUnpacker, error) {
	c, err := addr.Parse("version()(string)")
	if err != nil {
		return nil, err
	}
	return c()
}
func getSum(addr abi.Receiver) (abi.CallableUnpacker, error) {
	c, err := addr.Parse("sha256sum()(string)")
	if err != nil {
		return nil, err
	}
	return c()
}
func countItems(addr abi.Receiver) (abi.CallableUnpacker, error) {
	c, err := addr.Parse("count()(uint256)")
	if err != nil {
		return nil, err
	}
	return c()
}
