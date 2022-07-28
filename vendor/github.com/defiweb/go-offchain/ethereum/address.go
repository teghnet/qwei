package ethereum

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var ErrUnknownChainID = errors.New("unknown chain ID, use WithChainID or WithClientChainID to set it")
var ErrNotSupportedOnCurrentChain = errors.New("contract is not supported on the current chain")

type AddressProvider interface {
	Address(ctx context.Context) (common.Address, error)
}

type AddressesProvider interface {
	Addresses(ctx context.Context) (AddressList, error)
}

type AddressFunc func(ctx context.Context) (common.Address, error)

func (a AddressFunc) Address(ctx context.Context) (common.Address, error) {
	return a(ctx)
}

type NilAddress struct{}

func (NilAddress) Address(_ context.Context) (common.Address, error) {
	return common.Address{}, nil
}

type StaticAddress common.Address

func (a StaticAddress) Address(_ context.Context) (common.Address, error) {
	return common.Address(a), nil
}

type StaticAddresses []common.Address

func (p StaticAddresses) Addresses(_ context.Context) (AddressList, error) {
	return AddressList(p), nil
}

func HexToStaticAddress(hex string) StaticAddress {
	return StaticAddress(common.HexToAddress(hex))
}

type addressByChain struct {
	contracts map[*big.Int]common.Address
}

func AddressByChain(contracts map[*big.Int]common.Address) AddressProvider {
	return &addressByChain{contracts: contracts}
}

func (a *addressByChain) Address(ctx context.Context) (common.Address, error) {
	ctxCID := ChainIDFromContext(ctx)
	if ctxCID == nil {
		return ZeroAddress, ErrUnknownChainID
	}
	for cid, addr := range a.contracts {
		if cid.Cmp(ctxCID) == 0 {
			return addr, nil
		}
	}
	return ZeroAddress, ErrNotSupportedOnCurrentChain
}

type AddressProviders []AddressProvider

func (p AddressProviders) Addresses(ctx context.Context) (AddressList, error) {
	var list AddressList
	for _, a := range p {
		addr, err := a.Address(ctx)
		if err != nil {
			return nil, err
		}
		list.Add(addr)
	}
	return list, nil
}

func NilWhenZero(a common.Address, err error) (*common.Address, error) {
	if err != nil {
		return nil, err
	}
	if a == ZeroAddress {
		return nil, nil
	}
	return &a, nil
}
