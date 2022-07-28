package provider

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/defiweb/go-offchain/ethereum"
)

var _ ethereum.Client = (*poolClient)(nil)

type poolClient struct {
	mu sync.Mutex

	clients []ethereum.Client
	index   int
}

func NewPoolClient(clients ...ethereum.Client) ethereum.Client {
	return &poolClient{clients: clients}
}

func (c *poolClient) BalanceOf(ctx context.Context, owner common.Address) (*ethereum.Value, error) {
	return c.client().BalanceOf(ctx, owner)
}

func (c *poolClient) Transfer(ctx context.Context, to common.Address, txParams ethereum.TXParams) (*common.Hash, error) {
	return c.client().Transfer(ctx, to, txParams)
}

func (c *poolClient) Cancel(ctx context.Context, nonce ethereum.NonceProvider, fee ethereum.FeeEstimator) (*common.Hash, error) {
	return c.client().Cancel(ctx, nonce, fee)
}

func (c *poolClient) Receipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.client().Receipt(ctx, txHash)
}

func (c *poolClient) Nonce(ctx context.Context, address common.Address) (uint64, error) {
	return c.client().Nonce(ctx, address)
}

func (c *poolClient) PendingNonce(ctx context.Context, address common.Address) (uint64, error) {
	return c.client().PendingNonce(ctx, address)
}

func (c *poolClient) GasPrice(ctx context.Context) (*ethereum.Value, error) {
	return c.client().GasPrice(ctx)
}

func (c *poolClient) TipValue(ctx context.Context) (*ethereum.Value, error) {
	return c.client().TipValue(ctx)
}

func (c *poolClient) Block(ctx context.Context) (*types.Block, error) {
	return c.client().Block(ctx)
}

func (c *poolClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return c.client().BlockByHash(ctx, hash)
}

func (c *poolClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.client().BlockNumber(ctx)
}

func (c *poolClient) EstimateGas(ctx context.Context, call ethereum.Callable) (uint64, error) {
	return c.client().EstimateGas(ctx, call)
}

func (c *poolClient) Logs(ctx context.Context, filter ethereum.LogFilter, logParams ethereum.LogParams) ([]types.Log, error) {
	return c.client().Logs(ctx, filter, logParams)
}

func (c *poolClient) Storage(ctx context.Context, address common.Address, key common.Hash) ([]byte, error) {
	return c.client().Storage(ctx, address, key)
}

func (c *poolClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return c.client().NetworkID(ctx)
}

func (c *poolClient) Read(ctx context.Context, call ethereum.Callable, txParams ethereum.TXParams) (interface{}, error) {
	return c.client().Read(ctx, call, txParams)
}

func (c *poolClient) Write(ctx context.Context, call ethereum.Callable, txParams ethereum.TXParams) (*common.Hash, error) {
	return c.client().Write(ctx, call, txParams)
}

func (c *poolClient) Close(ctx context.Context) error {
	return c.client().Close(ctx)
}

func (c *poolClient) client() ethereum.Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	cli := c.clients[c.index]
	c.index = (c.index + 1) % len(c.clients)

	return cli
}
