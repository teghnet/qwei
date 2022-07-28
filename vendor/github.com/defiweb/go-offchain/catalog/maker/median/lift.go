package median

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type LiftCall struct {
	contract common.Address
	addrs    []common.Address
}

func Lift(contract common.Address, addrs ...common.Address) *LiftCall {
	return &LiftCall{
		contract: contract,
		addrs:    addrs,
	}
}

func (c *LiftCall) Address(ctx context.Context) (common.Address, error) {
	return c.contract, nil
}

func (c *LiftCall) Data(_ context.Context) (ethereum.CallData, error) {
	if len(c.addrs) == 0 {
		return nil, errors.New("no addresses to do the operation for")
	}
	return _ABI.Pack("lift(address[])", c.addrs)
}

func (c *LiftCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) error {
	_, err := client.Read(ctx, c, txParams)
	if err != nil {
		return err
	}
	return nil
}

func (c *LiftCall) Write(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*common.Hash, error) {
	return client.Write(ctx, c, txParams)
}
