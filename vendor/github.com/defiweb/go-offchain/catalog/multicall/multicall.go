package multicall

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

// ContractAddresses is a list of multicall contracts.
//
// https://github.com/makerdao/multicall
var ContractAddresses = map[*big.Int]common.Address{
	ethereum.MainnetChainID: common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"),
	ethereum.KovanChainID:   common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"),
	ethereum.RinkebyChainID: common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"),
	ethereum.GoerliChainID:  common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"),
	ethereum.RopstenChainID: common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"),
}

type Multicall struct {
	GetCurrentBlockTimestampMethod
	TryAggregateMethod
}

func Contract() *Multicall {
	return ContractWithAddressProvider(ethereum.AddressByChain(ContractAddresses))
}

func ContractWithAddressProvider(ap ethereum.AddressProvider) *Multicall {
	return &Multicall{
		GetCurrentBlockTimestampMethod: GetCurrentBlockTimestampMethod{AddressProvider: ap},
		TryAggregateMethod:             TryAggregateMethod{AddressProvider: ap},
	}
}

// Caller is a simple helper that checks if the Multicall contract is available on
// the current chain. if so, it uses it, otherwise it executes each call separately.
type Caller struct {
	client         ethereum.Client
	txParams       ethereum.TXParams
	requireSuccess bool
}

// NewCaller creates a new Caller instance.
func NewCaller(client ethereum.Client) *Caller {
	return &Caller{client: client}
}

// RequireSuccess sets the requireSuccess parameter to true.
func (c *Caller) RequireSuccess() *Caller {
	c.requireSuccess = true
	return c
}

// WithTXParams sets the transaction parameters for the caller.
func (c *Caller) WithTXParams(txParams ethereum.TXParams) *Caller {
	c.txParams = txParams
	return c
}

// Call executes calls using Multicall contract if available, otherwise it
// executes each call separately.
func (c *Caller) Call(ctx context.Context, calls ...ethereum.Callable) ([]AggregateResult, error) {
	if _, err := ethereum.AddressByChain(ContractAddresses).Address(ctx); err != nil {
		return Contract().TryAggregate(c.requireSuccess, calls...).Read(ctx, c.client, c.txParams)
	}
	var results []AggregateResult
	for _, call := range calls {
		result, err := c.client.Read(ctx, call, c.txParams)
		if c.requireSuccess && err != nil {
			return nil, err
		}
		results = append(results, AggregateResult{
			Success: err == nil,
			Call:    call,
			Result:  result,
		})
	}
	return results, nil
}
