package median

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type SetBarCall struct {
	contract common.Address
	bar      *big.Int
}

func SetBar(contract common.Address, bar int64) *SetBarCall {
	return &SetBarCall{
		contract: contract,
		bar:      big.NewInt(bar),
	}
}

func (c *SetBarCall) Address(ctx context.Context) (common.Address, error) {
	return c.contract, nil
}

func (c *SetBarCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("setBar(uint256)", c.bar)
}

func (c *SetBarCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) error {
	_, err := client.Read(ctx, c, txParams)
	if err != nil {
		return err
	}
	return nil
}

func (c *SetBarCall) Write(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*common.Hash, error) {
	return client.Write(ctx, c, txParams)
}
