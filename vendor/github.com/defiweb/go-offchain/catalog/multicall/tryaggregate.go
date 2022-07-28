package multicall

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type TryAggregateMethod struct {
	ethereum.AddressProvider
}

func (m *TryAggregateMethod) TryAggregate(requireSuccess bool, calls ...ethereum.Callable) *TryAggregateCall {
	return &TryAggregateCall{
		AddressProvider: m.AddressProvider,
		Calls:           calls,
		RequireSuccess:  requireSuccess,
	}
}

type TryAggregateCall struct {
	ethereum.AddressProvider
	Calls          []ethereum.Callable
	RequireSuccess bool
}

type AggregateResult struct {
	Success bool
	Call    ethereum.Callable
	Result  interface{}
}

func (c *TryAggregateCall) Data(ctx context.Context) (ethereum.CallData, error) {
	type call struct {
		Address common.Address `abi:"target"`
		Data    []byte         `abi:"callData"`
	}
	var cs []call
	for _, c := range c.Calls {
		address, err := ethereum.NilWhenZero(c.Address(ctx))
		if err != nil {
			return nil, err
		}
		if address == nil {
			return nil, ErrEmptyAddress
		}
		data, err := c.Data(ctx)
		if err != nil {
			return nil, err
		}
		cs = append(cs, call{
			Address: *address,
			Data:    data,
		})
	}
	return multicallABI.Pack("tryAggregate", c.RequireSuccess, cs)
}

func (c *TryAggregateCall) Unpack(data []byte) (interface{}, error) {
	u, err := multicallABI.Unpack("tryAggregate", data)
	if err != nil {
		return nil, err
	}
	var rs []AggregateResult
	for i, r := range u[0].([]struct {
		Success    bool    `json:"success"`
		ReturnData []uint8 `json:"returnData"`
	}) {
		result := AggregateResult{}
		result.Call = c.Calls[i]
		if cr, ok := c.Calls[i].(ethereum.Unpacker); ok {
			if unpacked, err := cr.Unpack(r.ReturnData); err != nil {
				result.Success = false
				result.Result = err
			} else {
				result.Success = r.Success
				result.Result = unpacked
			}
		} else {
			result.Success = r.Success
			result.Result = r.ReturnData
		}
		rs = append(rs, result)
	}
	return rs, nil
}

func (c *TryAggregateCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) ([]AggregateResult, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return nil, err
	}
	return r.([]AggregateResult), nil
}

func (c *TryAggregateCall) Write(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (*common.Hash, error) {
	return client.Write(ctx, c, txParams)
}
