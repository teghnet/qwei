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
	addr := common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F")
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
	var callables []ethereum.Callable
	var version string
	{
		c, err := getVersion(addr)
		if err != nil {
			return err
		}
		callables = append(callables, c.Bind(&version))
	}
	// var sum string
	// {
	// 	c, err := getSum(addr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	callables = append(callables, c.Bind(&sum))
	// }
	// var ipfs string
	// {
	// 	c, err := getIPFS(addr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	callables = append(callables, c.Bind(&ipfs))
	// }
	{
		cs, err := getItems(addr, count)
		if err != nil {
			return err
		}
		callables = append(callables, cs...)

		parsed, err := abi.Parse("tryAggregate(bool,(address,bytes)[])((bool,bytes)[])")
		if err != nil {
			return err
		}
		unpacker, err := parsed.
			Attach(ethereum.AddressByChain(multicall.ContractAddresses))(false, callables)
		if err != nil {
			return err
		}

		var fn abi.ValFn = func(arg any) error {
			for i, r := range arg.([]struct {
				Name0 bool    `json:"name0"`
				Name1 []uint8 `json:"name1"`
			}) {
				if cr, ok := callables[i].(ethereum.Unpacker); ok {
					if unpacked, err := cr.Unpack(r.Name1); err != nil {
						log.Printf("unable to unpack %d: %s", i, err)
					} else {
						fmt.Println(i, unpacked)
					}
				} else {
					fmt.Println(i, r)
				}
			}
			return nil
		}
		_, err = client.Read(ctx, unpacker.Bind(fn), params)
		if err != nil {
			return err
		}
	}
	fmt.Println("version:", version)
	fmt.Println("count:", count)
	// fmt.Println("sum:", sum)
	// fmt.Println("ipfs:", ipfs)
	return nil
}

func getItems(addr common.Address, count int) ([]ethereum.Callable, error) {
	method, err := abi.ParseWithAddr(addr, "get(uint256)(bytes32,address)")
	// TODO: we could try creating a set of default mappers
	// method, err := addr.Parse("get(uint256)(bytes32=>string,address)")
	if err != nil {
		return nil, err
	}
	var cs []ethereum.Callable
	var s string
	for i := 0; i < count; i++ {
		c, err := method(i)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c.Bind(&s))
	}
	return cs, nil
}
func getIPFS(addr common.Address) (abi.CallableUnpacker, error) {
	c, err := abi.ParseWithAddr(addr, "ipfs()(string)")
	if err != nil {
		return nil, err
	}
	return c()
}
func getVersion(addr common.Address) (abi.CallableUnpacker, error) {
	c, err := abi.ParseWithAddr(addr, "version()(string)")
	if err != nil {
		return nil, err
	}
	return c()
}
func getSum(addr common.Address) (abi.CallableUnpacker, error) {
	c, err := abi.ParseWithAddr(addr, "sha256sum()(string)")
	if err != nil {
		return nil, err
	}
	return c()
}
func countItems(addr common.Address) (abi.CallableUnpacker, error) {
	c, err := abi.ParseWithAddr(addr, "count()(uint256)")
	if err != nil {
		return nil, err
	}
	return c()
}
