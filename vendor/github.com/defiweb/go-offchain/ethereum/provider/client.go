package provider

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/defiweb/go-offchain/ethereum"
)

var ErrClientFunctionNotFound = errors.New("client function not found")

var _ ethereum.Client = (*baseClient)(nil)

type baseClient struct {
	client clientFn
}

func newBaseClient(client clientFn) *baseClient {
	return &baseClient{client: client}
}

type clientFn func(chainID *big.Int) (ethereum.Client, error)

func (cp *baseClient) clientByChainFromContext(ctx context.Context) (ethereum.Client, error) {
	if cp.client == nil {
		return nil, ErrClientFunctionNotFound
	}
	chainID := ethereum.ChainIDFromContext(ctx)
	if chainID == nil {
		chainID = ethereum.MainnetChainID
	}
	return cp.client(chainID)
}

func (cp *baseClient) BalanceOf(ctx context.Context, owner common.Address) (*ethereum.Value, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.BalanceOf(ctx, owner)
}

func (cp *baseClient) Transfer(ctx context.Context, to common.Address, txParams ethereum.TXParams) (*common.Hash, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Transfer(ctx, to, txParams)
}

func (cp *baseClient) Cancel(ctx context.Context, nonce ethereum.NonceProvider, fee ethereum.FeeEstimator) (*common.Hash, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Cancel(ctx, nonce, fee)
}

func (cp *baseClient) Receipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Receipt(ctx, txHash)
}

func (cp *baseClient) Nonce(ctx context.Context, address common.Address) (uint64, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return c.Nonce(ctx, address)
}

func (cp *baseClient) PendingNonce(ctx context.Context, address common.Address) (uint64, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return c.PendingNonce(ctx, address)
}

func (cp *baseClient) GasPrice(ctx context.Context) (*ethereum.Value, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.GasPrice(ctx)
}

func (cp *baseClient) TipValue(ctx context.Context) (*ethereum.Value, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.TipValue(ctx)
}

func (cp *baseClient) Block(ctx context.Context) (*types.Block, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Block(ctx)
}

func (cp *baseClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.BlockByHash(ctx, hash)
}

func (cp *baseClient) BlockNumber(ctx context.Context) (uint64, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return c.BlockNumber(ctx)
}

func (cp *baseClient) EstimateGas(ctx context.Context, call ethereum.Callable) (uint64, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return c.EstimateGas(ctx, call)
}

func (cp *baseClient) Logs(ctx context.Context, filter ethereum.LogFilter, logParams ethereum.LogParams) ([]types.Log, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Logs(ctx, filter, logParams)
}

func (cp *baseClient) Storage(ctx context.Context, address common.Address, key common.Hash) ([]byte, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Storage(ctx, address, key)
}

func (cp *baseClient) NetworkID(ctx context.Context) (*big.Int, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.NetworkID(ctx)
}

func (cp *baseClient) Read(ctx context.Context, call ethereum.Callable, txParams ethereum.TXParams) (interface{}, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Read(ctx, call, txParams)
}

func (cp *baseClient) Write(ctx context.Context, call ethereum.Callable, txParams ethereum.TXParams) (*common.Hash, error) {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return c.Write(ctx, call, txParams)
}

func (cp *baseClient) Close(ctx context.Context) error {
	c, err := cp.clientByChainFromContext(ctx)
	if err != nil {
		return err
	}
	return c.Close(ctx)
}
