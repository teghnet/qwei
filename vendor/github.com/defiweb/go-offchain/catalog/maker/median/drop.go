package median

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type DropCall struct {
	contract common.Address
	addrs    []common.Address
}

func Drop(contract common.Address, addrs ...common.Address) *DropCall {
	return &DropCall{
		contract: contract,
		addrs:    addrs,
	}
}

func (c *DropCall) Address(ctx context.Context) (common.Address, error) {
	return c.contract, nil
}

func (c *DropCall) Data(_ context.Context) (ethereum.CallData, error) {
	if len(c.addrs) == 0 {
		return nil, errors.New("no addresses to do the operation for")
	}
	return _ABI.Pack("drop(address[])", c.addrs)
}

func (c *DropCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) error {
	_, err := client.Read(ctx, c, txParams)
	if err != nil {
		return err
	}
	return nil
}

func (c *DropCall) Write(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*common.Hash, error) {
	return client.Write(ctx, c, txParams)
}
