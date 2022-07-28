package median

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/catalog/multicall"

	"github.com/defiweb/go-offchain/ethereum"
)

type OrclsCall struct {
	contract  common.Address
	multicall *multicall.TryAggregateCall
}

func Orcls(contract common.Address) *OrclsCall {
	var cs []ethereum.Callable
	for i := 0; i < 256; i++ {
		cs = append(cs, Slot(contract, i))
	}
	return &OrclsCall{
		contract:  contract,
		multicall: multicall.Contract().TryAggregate(true, cs...),
	}
}

func (c *OrclsCall) Address(ctx context.Context) (common.Address, error) {
	return c.multicall.Address(ctx)
}

func (c *OrclsCall) Data(ctx context.Context) (ethereum.CallData, error) {
	return c.multicall.Data(ctx)
}

func (c *OrclsCall) Unpack(data []byte) (interface{}, error) {
	u, err := c.multicall.Unpack(data)
	if err != nil {
		return nil, err
	}
	var addrs []common.Address
	for _, r := range u.([]multicall.AggregateResult) {
		addr := r.Result.(common.Address)
		if addr != ethereum.ZeroAddress {
			addrs = append(addrs, addr)
		}
	}
	return addrs, nil
}

func (c *OrclsCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) ([]common.Address, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return nil, err
	}
	return r.([]common.Address), nil
}
