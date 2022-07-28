package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const contextBlockNumber = "ethereum_block_number"
const contextChainID = "ethereum_chain_id"
const contextAccount = "ethereum_account"

// WithBlockNumber adds the block number to the context.
func WithBlockNumber(ctx context.Context, block uint64) context.Context {
	return context.WithValue(ctx, contextBlockNumber, block)
}

// WithClientChainID adds the chain ID to the context using the ID provided by
// the client.
func WithClientChainID(ctx context.Context, client Client) (context.Context, error) {
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, contextChainID, chainID), nil
}

// WithChainID adds the chain ID to the context.
func WithChainID(ctx context.Context, chainID *big.Int) context.Context {
	return context.WithValue(ctx, contextChainID, chainID)
}

// WithAccount adds the account to the context.
func WithAccount(ctx context.Context, account Account) context.Context {
	return context.WithValue(ctx, contextAccount, account)
}

// BlockNumberFromContext returns the block number from the context.
func BlockNumberFromContext(ctx context.Context) uint64 {
	n, ok := ctx.Value(contextBlockNumber).(uint64)
	if ok {
		return n
	}
	return 0
}

// ChainIDFromContext returns the chain ID from the context.
func ChainIDFromContext(ctx context.Context) *big.Int {
	id, ok := ctx.Value(contextChainID).(*big.Int)
	if ok {
		return id
	}
	return nil
}

// AccountFromContext returns the account from the context.
func AccountFromContext(ctx context.Context) Account {
	id, ok := ctx.Value(contextAccount).(Account)
	if ok {
		return id
	}
	return ZeroAccount
}

// AddressFromContext returns the account address from the context.
func AddressFromContext(ctx context.Context) common.Address {
	return mustAddress(AccountFromContext(ctx).Address(ctx))
}

func mustAddress(a common.Address, e error) common.Address {
	if e != nil {
		panic(e)
	}
	return a
}
