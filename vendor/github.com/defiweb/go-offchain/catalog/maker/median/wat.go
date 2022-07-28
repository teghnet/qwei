package median

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type WatCall struct {
	contract common.Address
}

func Wat(contract common.Address) *WatCall {
	return &WatCall{
		contract: contract,
	}
}

func (c *WatCall) Address(ctx context.Context) (common.Address, error) {
	return c.contract, nil
}

func (c *WatCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("wat()")
}

func (c *WatCall) Unpack(data []byte) (interface{}, error) {
	u, err := _ABI.Unpack("wat()", data)
	if err != nil {
		return nil, err
	}
	b := u[0].([32]byte)
	return string(b[:]), nil
}

func (c *WatCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (string, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return "", err
	}
	return r.(string), nil
}
