package ethereum

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type CallDataProvider interface {
	Data(ctx context.Context) (CallData, error)
}

type CallDataFunc func(ctx context.Context) (CallData, error)

func (c CallDataFunc) Data(ctx context.Context) (CallData, error) {
	return c(ctx)
}

type CallData []byte

func (c CallData) Data(_ context.Context) (CallData, error) {
	return c, nil
}

func (c CallData) String() string {
	return hex.EncodeToString(c)
}

func NewCallDataFromHex(s string) (CallData, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}

// Unpacker can unpack response from the EVM.
type Unpacker interface {
	Unpack([]byte) (interface{}, error)
}

type Callable interface {
	AddressProvider
	CallDataProvider
}

type call struct {
	ap AddressProvider
	dp CallDataProvider
}

func (c *call) Data(ctx context.Context) (CallData, error) {
	if c.dp == nil {
		return nil, nil
	}
	return c.dp.Data(ctx)
}

func (c *call) Address(ctx context.Context) (common.Address, error) {
	if c.ap == nil {
		return common.Address{}, nil
	}
	return c.ap.Address(ctx)
}

func NewCall(address AddressProvider, callData CallDataProvider) Callable {
	return &call{
		ap: address,
		dp: callData,
	}
}

func NewTransferCall(address AddressProvider) Callable {
	return NewCall(address, nil)
}

func NewConstructorCall(callData CallData) Callable {
	return NewCall(nil, callData)
}
