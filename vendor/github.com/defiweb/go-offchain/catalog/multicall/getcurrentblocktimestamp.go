package multicall

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type GetCurrentBlockTimestampMethod struct {
	ethereum.AddressProvider
}

func (m *GetCurrentBlockTimestampMethod) GetCurrentBlockTimestamp() *GetCurrentBlockTimestampCall {
	return &GetCurrentBlockTimestampCall{
		AddressProvider: m.AddressProvider,
	}
}

type GetCurrentBlockTimestampCall struct {
	ethereum.AddressProvider
}

func (c *GetCurrentBlockTimestampCall) Data(ctx context.Context) (ethereum.CallData, error) {
	return multicallABI.Pack("getCurrentBlockTimestamp")
}

func (c *GetCurrentBlockTimestampCall) Unpack(data []byte) (interface{}, error) {
	u, err := multicallABI.Unpack("getCurrentBlockTimestamp", data)
	if err != nil {
		return nil, err
	}
	return time.Unix(u[0].(*big.Int).Int64(), 0), nil
}

func (c *GetCurrentBlockTimestampCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (time.Time, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return time.Time{}, err
	}
	return r.(time.Time), nil
}

func (c *GetCurrentBlockTimestampCall) Write(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*common.Hash, error) {
	return client.Write(ctx, c, txParams)
}
