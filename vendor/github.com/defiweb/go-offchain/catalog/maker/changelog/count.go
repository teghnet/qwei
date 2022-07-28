package changelog

import (
	"context"
	"math/big"

	"github.com/defiweb/go-offchain/bn"
	"github.com/defiweb/go-offchain/ethereum"
)

type CountCall struct {
	ethereum.AddressProvider
}

func Count() *CountCall {
	return &CountCall{
		AddressProvider: ethereum.AddressByChain(Contracts),
	}
}

func (c *CountCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("count()")
}

func (c *CountCall) Unpack(data []byte) (interface{}, error) {
	u, err := _ABI.Unpack("count()", data)
	if err != nil {
		return nil, err
	}
	return bn.IntFromBigInt(u[0].(*big.Int)), nil
}

func (c *CountCall) Values(u interface{}) (*bn.Int, error) {
	return u.(*bn.Int), nil
}

func (c *CountCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*bn.Int, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return nil, err
	}
	return c.Values(r)
}
