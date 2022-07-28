package ethereum

import (
	"github.com/defiweb/go-offchain/bn"
)

const (
	UnitWei   = 1
	UnitGWei  = 1e9
	UnitEther = 1e18
)

type Value bn.Int

func Wei(v *bn.Int) *Value {
	return (*Value)(v)
}

func GWei(v *bn.Float) *Value {
	return (*Value)(v.Mul(bn.FloatFromUint64(UnitGWei)).Int())
}

func Ether(v *bn.Float) *Value {
	return (*Value)(v.Mul(bn.FloatFromUint64(UnitEther)).Int())
}

func (v *Value) String() string {
	return v.Ether().String()
}

func (v *Value) Wei() *bn.Int {
	if v == nil {
		return bn.IntFromUint64(0)
	}
	return (*bn.Int)(v)
}

func (v *Value) GWei() *bn.Float {
	if v == nil {
		return bn.FloatFromFloat64(0)
	}
	return (*bn.Int)(v).Float().Div(bn.FloatFromUint64(UnitGWei))
}

func (v *Value) Ether() *bn.Float {
	if v == nil {
		return bn.FloatFromFloat64(0)
	}
	return (*bn.Int)(v).Float().Div(bn.FloatFromUint64(UnitEther))
}

func GasPrice(gas uint64, price *Value) *Value {
	return (*Value)(price.Wei().Mul(bn.IntFromUint64(gas)))
}
