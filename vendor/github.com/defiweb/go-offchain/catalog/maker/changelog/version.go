package changelog

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type VersionCall struct {
	ethereum.AddressProvider
	contract common.Address
}

func Version() *VersionCall {
	return &VersionCall{
		AddressProvider: ethereum.AddressByChain(Contracts),
	}
}

func (c *VersionCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("version()")
}

func (c *VersionCall) Unpack(data []byte) (interface{}, error) {
	u, err := _ABI.Unpack("version()", data)
	if err != nil {
		return nil, err
	}
	return u[0], nil
}

func (c *VersionCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (string, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return "", err
	}
	return c.Values(r)
}

func (c *VersionCall) Values(u interface{}) (string, error) {
	return u.(string), nil
}
