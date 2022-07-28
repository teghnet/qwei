package ethereum

import (
	"context"

	"github.com/defiweb/go-offchain/bn"
)

type FeeEstimator interface {
	// TipValue is the maximum tip value. If nil, the suggested tip
	// value will be used.
	TipValue(ctx context.Context) (*Value, error)
	// MaxPrice is the maximum gas price. If nil, then suggested gas price
	// multiplied by two will be used.
	MaxPrice(ctx context.Context) (*Value, error)
}

type staticFee struct {
	tipValue *Value
	maxPrice *Value
}

func NewStaticFee(tipValue, maxPrice *Value) FeeEstimator {
	if tipValue == nil {
		tipValue = Wei(bn.IntFromInt64(0))
	}
	if maxPrice == nil {
		maxPrice = Wei(bn.IntFromInt64(0))
	}
	return &staticFee{
		tipValue: tipValue,
		maxPrice: maxPrice,
	}
}

func (f *staticFee) TipValue(_ context.Context) (*Value, error) {
	return f.tipValue, nil
}

func (f *staticFee) MaxPrice(_ context.Context) (*Value, error) {
	return f.maxPrice, nil
}

type suggestedFee struct {
	client             Client
	tipValueMultiplier float64
	maxPriceMultiplier float64
}

func NewSuggestedFee(client Client, tipValueMultiplier, maxPriceMultiplier float64) FeeEstimator {
	return &suggestedFee{
		client:             client,
		tipValueMultiplier: tipValueMultiplier,
		maxPriceMultiplier: maxPriceMultiplier,
	}
}

func (f *suggestedFee) TipValue(ctx context.Context) (*Value, error) {
	v, err := f.client.TipValue(ctx)
	if err != nil {
		return nil, err
	}
	return Wei(v.Wei().Float().Mul(bn.FloatFromFloat64(f.tipValueMultiplier)).Int()), nil
}

func (f *suggestedFee) MaxPrice(ctx context.Context) (*Value, error) {
	v, err := f.client.GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return Wei(v.Wei().Float().Mul(bn.FloatFromFloat64(f.maxPriceMultiplier)).Int()), nil
}
