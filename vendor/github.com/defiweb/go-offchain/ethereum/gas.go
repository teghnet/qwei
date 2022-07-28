package ethereum

import (
	"context"
	"errors"
)

type GasLimiter interface {
	GasLimit(ctx context.Context, call Callable) (uint64, error)
}

type GasLimit uint64

func (c GasLimit) GasLimit(_ context.Context, _ Callable) (uint64, error) {
	return uint64(c), nil
}

type gasLimitEstimator struct {
	client Client
}

func NewGasLimitEstimator(client Client) GasLimiter {
	return &gasLimitEstimator{client: client}
}

func (g *gasLimitEstimator) GasLimit(ctx context.Context, call Callable) (uint64, error) {
	if call == nil {
		return 0, errors.New("call cannot be empty")
	}
	data, err := call.Data(ctx)
	if err != nil {
		return 0, err
	}
	if data == nil {
		return TransferGas, nil
	}
	v, err := g.client.EstimateGas(ctx, call)
	if err != nil {
		return 0, err
	}
	return uint64(float64(v) * defaultGasEstimationMultiplier), nil
}
