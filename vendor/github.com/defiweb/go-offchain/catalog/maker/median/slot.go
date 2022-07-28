package median

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type SlotCall struct {
	contract common.Address
	idx      uint8
}

func Slot(contract common.Address, idx int) *SlotCall {
	return &SlotCall{
		contract: contract,
		idx:      uint8(idx),
	}
}

func (c *SlotCall) Address(ctx context.Context) (common.Address, error) {
	return c.contract, nil
}

func (c *SlotCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("slot(uint8)", c.idx)
}

func (c *SlotCall) Unpack(data []byte) (interface{}, error) {
	u, err := _ABI.Unpack("slot(uint8)", data)
	if err != nil {
		return nil, err
	}
	return u[0].(common.Address), nil
}

func (c *SlotCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (common.Address, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return common.Address{}, err
	}
	return c.Values(r)
}

func (c *SlotCall) Values(u interface{}) (common.Address, error) {
	address, ok := u.(common.Address)
	if !ok {
		return common.Address{}, ethereum.ErrExpectingAddress
	}
	return address, nil
}
