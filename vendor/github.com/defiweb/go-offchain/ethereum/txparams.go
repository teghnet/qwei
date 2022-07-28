package ethereum

import (
	"context"
	"errors"

	"github.com/defiweb/go-offchain/bn"
)

const TransferGas = 21000
const defaultTipValueMultiplier = 1
const defaultMaxPriceMultiplier = 2
const defaultGasEstimationMultiplier = 1.2

type TXParams interface {
	NonceProvider
	GasLimiter
	FeeEstimator
	Amount(ctx context.Context) (*Value, error)
}

type txParams struct {
	NonceProvider
	GasLimiter
	FeeEstimator
	amount *Value
}

func NewTXParams(amount *Value, nonce NonceProvider, gasLimit GasLimiter, fee FeeEstimator) TXParams {
	if amount == nil {
		amount = Wei(bn.IntFromInt64(0))
	}
	if nonce == nil {
		nonce = Nonce(0)
	}
	if gasLimit == nil {
		gasLimit = GasLimit(0)
	}
	if fee == nil {
		fee = NewStaticFee(nil, nil)
	}
	return &txParams{
		NonceProvider: nonce,
		GasLimiter:    gasLimit,
		FeeEstimator:  fee,
		amount:        amount,
	}
}

func NewAutoTXParams(client Client, amount *Value) (TXParams, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if amount == nil {
		amount = Wei(bn.IntFromInt64(0))
	}
	return &txParams{
		NonceProvider: NewPendingNonce(client),
		GasLimiter:    NewGasLimitEstimator(client),
		FeeEstimator:  NewSuggestedFee(client, defaultTipValueMultiplier, defaultMaxPriceMultiplier),
		amount:        amount,
	}, nil
}

func (t *txParams) Amount(_ context.Context) (*Value, error) {
	return t.amount, nil
}
