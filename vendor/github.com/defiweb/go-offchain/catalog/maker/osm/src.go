package osm

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type SrcCall struct {
	contract common.Address
}

func Src(contract common.Address) *SrcCall {
	return &SrcCall{
		contract: contract,
	}
}

func (c *SrcCall) Address(ctx context.Context) (common.Address, error) {
	return c.contract, nil
}

func (c *SrcCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("src()")
}

func (c *SrcCall) Unpack(data []byte) (interface{}, error) {
	u, err := _ABI.Unpack("src()", data)
	if err != nil {
		return nil, err
	}
	return u[0].(common.Address), nil
}

func (c *SrcCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (common.Address, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return common.Address{}, err
	}
	return c.Values(r)
}

func (c *SrcCall) Values(u interface{}) (common.Address, error) {
	address, ok := u.(common.Address)
	if !ok {
		return common.Address{}, ethereum.ErrWrongInterface
	}
	return address, nil
}
