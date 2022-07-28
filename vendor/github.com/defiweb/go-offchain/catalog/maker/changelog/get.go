package changelog

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

type GetCall struct {
	ethereum.AddressProvider
	idx *big.Int
}

func Get(idx uint64) *GetCall {
	return &GetCall{
		AddressProvider: ethereum.AddressByChain(Contracts),
		idx:             big.NewInt(0).SetUint64(idx),
	}
}

func (c *GetCall) Data(_ context.Context) (ethereum.CallData, error) {
	return _ABI.Pack("get(uint256)", c.idx)
}

func (c *GetCall) Unpack(data []byte) (interface{}, error) {
	u, err := _ABI.Unpack("get(uint256)", data)
	if err != nil {
		return nil, err
	}
	b, ok := u[0].([32]byte)
	if !ok {
		return nil, ethereum.ErrExpecting32ByteArray
	}
	return []interface{}{strings.Trim(string(b[:]), "\x00"), u[1].(common.Address)}, nil
}

func (c *GetCall) Read(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (string, common.Address, error) {
	r, err := client.Read(ctx, c, txParams)
	if err != nil {
		return "", common.Address{}, err
	}
	v, ok := r.([]interface{})
	if !ok {
		return "", common.Address{}, ethereum.ErrExpectingSliceOfInterfaces
	}
	return v[0].(string), v[1].(common.Address), nil
}
