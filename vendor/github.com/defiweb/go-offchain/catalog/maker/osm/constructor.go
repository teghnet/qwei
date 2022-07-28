package osm

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type ConstructorCall struct {
	src common.Address
	bin []byte
}

func Constructor(src common.Address) *ConstructorCall {
	return &ConstructorCall{
		src: src,
		bin: _Bin,
	}
}

func (c *ConstructorCall) Address(_ context.Context) (common.Address, error) {
	return common.Address{}, nil
}

func (c *ConstructorCall) Data(_ context.Context) (ethereum.CallData, error) {
	data, err := _ABI.Pack("", c.src)
	if err != nil {
		return nil, err
	}
	return append(c.bin, data...), nil
}

func (c *ConstructorCall) Init(_ context.Context) (ethereum.CallData, error) {
	return c.bin, nil
}

func (c *ConstructorCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) error {
	_, err := client.Read(ctx, c, txParams)
	return err
}

func (c *ConstructorCall) Write(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*common.Hash, error) {
	return client.Write(ctx, c, txParams)
}
