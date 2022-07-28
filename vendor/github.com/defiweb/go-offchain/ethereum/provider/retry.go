package provider

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/defiweb/go-offchain/ethereum"
)

var _ ethereum.Client = (*retryClient)(nil)

type retryClient struct {
	client   ethereum.Client
	attempts int
	delay    time.Duration
}

func NewRetryClient(client ethereum.Client, attempts int, delay time.Duration) ethereum.Client {
	return &retryClient{
		client:   client,
		attempts: attempts,
		delay:    delay,
	}
}

func (c *retryClient) BalanceOf(ctx context.Context, owner common.Address) (*ethereum.Value, error) {
	var res *ethereum.Value
	err := c.retry(func() error {
		var err error
		res, err = c.client.BalanceOf(ctx, owner)
		return err
	})
	return res, err
}

func (c *retryClient) Transfer(ctx context.Context, to common.Address, txParams ethereum.TXParams) (*common.Hash, error) {
	var res *common.Hash
	err := c.retry(func() error {
		var err error
		res, err = c.client.Transfer(ctx, to, txParams)
		return err
	})
	return res, err
}

func (c *retryClient) Cancel(ctx context.Context, nonce ethereum.NonceProvider, fee ethereum.FeeEstimator) (*common.Hash, error) {
	var res *common.Hash
	err := c.retry(func() error {
		var err error
		res, err = c.client.Cancel(ctx, nonce, fee)
		return err
	})
	return res, err
}

func (c *retryClient) Receipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var res *types.Receipt
	err := c.retry(func() error {
		var err error
		res, err = c.client.Receipt(ctx, txHash)
		return err
	})
	return res, err
}

func (c *retryClient) Nonce(ctx context.Context, address common.Address) (uint64, error) {
	var res uint64
	err := c.retry(func() error {
		var err error
		res, err = c.client.Nonce(ctx, address)
		return err
	})
	return res, err
}

func (c *retryClient) PendingNonce(ctx context.Context, address common.Address) (uint64, error) {
	var res uint64
	err := c.retry(func() error {
		var err error
		res, err = c.client.PendingNonce(ctx, address)
		return err
	})
	return res, err
}

func (c *retryClient) GasPrice(ctx context.Context) (*ethereum.Value, error) {
	var res *ethereum.Value
	err := c.retry(func() error {
		var err error
		res, err = c.client.GasPrice(ctx)
		return err
	})
	return res, err
}

func (c *retryClient) TipValue(ctx context.Context) (*ethereum.Value, error) {
	var res *ethereum.Value
	err := c.retry(func() error {
		var err error
		res, err = c.client.TipValue(ctx)
		return err
	})
	return res, err
}

func (c *retryClient) Block(ctx context.Context) (*types.Block, error) {
	var res *types.Block
	err := c.retry(func() error {
		var err error
		res, err = c.client.Block(ctx)
		return err
	})
	return res, err
}

func (c *retryClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	var res *types.Block
	err := c.retry(func() error {
		var err error
		res, err = c.client.BlockByHash(ctx, hash)
		return err
	})
	return res, err
}

func (c *retryClient) BlockNumber(ctx context.Context) (uint64, error) {
	var res uint64
	err := c.retry(func() error {
		var err error
		res, err = c.client.BlockNumber(ctx)
		return err
	})
	return res, err
}

func (c *retryClient) EstimateGas(ctx context.Context, call ethereum.Callable) (uint64, error) {
	var res uint64
	err := c.retry(func() error {
		var err error
		res, err = c.client.EstimateGas(ctx, call)
		return err
	})
	return res, err
}

func (c *retryClient) Logs(ctx context.Context, filter ethereum.LogFilter, logParams ethereum.LogParams) ([]types.Log, error) {
	var res []types.Log
	err := c.retry(func() error {
		var err error
		res, err = c.client.Logs(ctx, filter, logParams)
		return err
	})
	return res, err
}

func (c *retryClient) Storage(ctx context.Context, address common.Address, key common.Hash) ([]byte, error) {
	var res []byte
	err := c.retry(func() error {
		var err error
		res, err = c.client.Storage(ctx, address, key)
		return err

	})
	return res, err
}

func (c *retryClient) NetworkID(ctx context.Context) (*big.Int, error) {
	var res *big.Int
	err := c.retry(func() error {
		var err error
		res, err = c.client.NetworkID(ctx)
		return err
	})
	return res, err
}

func (c *retryClient) Read(ctx context.Context, call ethereum.Callable, txParams ethereum.TXParams) (interface{}, error) {
	var res interface{}
	err := c.retry(func() error {
		var err error
		res, err = c.client.Read(ctx, call, txParams)
		return err
	})
	return res, err
}

func (c *retryClient) Write(ctx context.Context, call ethereum.Callable, txParams ethereum.TXParams) (*common.Hash, error) {
	var res *common.Hash
	err := c.retry(func() error {
		var err error
		res, err = c.client.Write(ctx, call, txParams)
		return err
	})
	return res, err
}

func (c *retryClient) Close(ctx context.Context) error {
	return c.retry(func() error {
		return c.client.Close(ctx)
	})
}

func (c *retryClient) retry(f func() error) error {
	var err error
	for i := 0; i < c.attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}
		if strings.Index(err.Error(), "abi:") == 0 {
			return err
		}
		if errors.As(err, &ethereum.ErrRevert{}) {
			return err
		}
		time.Sleep(c.delay)
	}
	return err
}
