package bn

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
)

type Fixed struct {
	n uint16
	x *big.Int
}

func FixedFromString(n uint16, s string) *Fixed {
	x, ok := new(big.Float).SetString(s)
	if !ok {
		return nil
	}
	return &Fixed{n: n, x: floatToRawInt(n, x)}
}

func FixedFromFloat64(n uint16, x float64) *Fixed {
	return &Fixed{n: n, x: floatToRawInt(n, big.NewFloat(x))}
}

func FixedFromInt64(n uint16, x int64) *Fixed {
	return &Fixed{n: n, x: floatToRawInt(n, IntFromInt64(x).BigFloat())}
}

func FixedFromUint64(n uint16, x uint64) *Fixed {
	return &Fixed{n: n, x: floatToRawInt(n, IntFromUint64(x).BigFloat())}
}

func FixedFromBigFloat(n uint16, x *big.Float) *Fixed {
	return &Fixed{n: n, x: floatToRawInt(n, x)}
}

func FixedFromBigInt(n uint16, x *big.Int) *Fixed {
	return &Fixed{n: n, x: floatToRawInt(n, IntFromBigInt(x).BigFloat())}
}

func FixedFromRawBigInt(n uint16, x *big.Int) *Fixed {
	return &Fixed{n: n, x: x}
}

func FixedFromRawInt(n uint16, x *Int) *Fixed {
	return &Fixed{n: n, x: x.x}
}

func FixedFromBytes(n uint16, x []byte) *Fixed {
	return &Fixed{n: n, x: new(big.Int).SetBytes(x)}
}

func FixedFromFloat(n uint16, x *Float) *Fixed {
	return &Fixed{n: n, x: floatToRawInt(n, x.BigFloat())}
}

func FixedFromInt(n uint16, x *Int) *Fixed {
	return &Fixed{n: n, x: x.Lsh(uint(n)).BigInt()}
}

func (f *Fixed) String() string {
	return f.BigFloat().String()
}

func (f *Fixed) Text(base int) string {
	return f.x.Text(base)
}

func (f *Fixed) Float() *Float {
	return FloatFromBigFloat(f.BigFloat())
}

func (f *Fixed) Int() *Int {
	return IntFromBigInt(f.BigInt())
}

func (f *Fixed) RawInt() *Int {
	return IntFromBigInt(f.x)
}

func (f *Fixed) BigFloat() *big.Float {
	return rawIntToFloat(f.n, f.x)
}

func (f *Fixed) BigInt() *big.Int {
	return new(big.Int).Rsh(f.x, uint(f.n))
}

func (f *Fixed) RawBigInt() *big.Int {
	return f.x
}

func (f *Fixed) Float64() float64 {
	f64, _ := rawIntToFloat(f.n, f.x).Float64()
	return f64
}

func (f *Fixed) Bytes() []byte {
	return f.x.Bytes()
}

func (f *Fixed) FractionBits() uint16 {
	return f.n
}

func (f *Fixed) SetFractionBits(n uint16) *Fixed {
	if f.n > n {
		return &Fixed{n: n, x: new(big.Int).Rsh(f.x, uint(f.n-n))}
	}
	if f.n < n {
		return &Fixed{n: n, x: new(big.Int).Lsh(f.x, uint(n-f.n))}
	}
	return f
}

func (f *Fixed) Sign() int {
	return f.x.Sign()
}

func (f *Fixed) Add(x *Fixed) *Fixed {
	x = x.SetFractionBits(f.n)
	return &Fixed{n: f.n, x: new(big.Int).Add(f.x, x.x)}
}

func (f *Fixed) Sub(x *Fixed) *Fixed {
	x = x.SetFractionBits(f.n)
	return &Fixed{n: f.n, x: new(big.Int).Sub(f.x, x.x)}
}

func (f *Fixed) Mul(x *Fixed) *Fixed {
	x = x.SetFractionBits(f.n)
	return &Fixed{n: f.n, x: new(big.Int).Rsh(new(big.Int).Mul(f.x, x.x), uint(f.n))}
}

func (f *Fixed) Div(x *Fixed) *Fixed {
	x = x.SetFractionBits(f.n)
	return &Fixed{n: f.n, x: new(big.Int).Div(new(big.Int).Lsh(f.x, uint(f.n)), x.x)}
}

// TODO: pow
// TODO: sqrt

func (f *Fixed) Cmp(x *Fixed) int {
	x = x.SetFractionBits(f.n)
	return f.x.Cmp(x.x)
}

func (f *Fixed) Lsh(n uint) *Fixed {
	return &Fixed{n: f.n, x: new(big.Int).Lsh(f.x, n)}
}

func (f *Fixed) Rsh(n uint) *Fixed {
	return &Fixed{n: f.n, x: new(big.Int).Rsh(f.x, n)}
}

func (f *Fixed) Floor() *Fixed {
	modf := f.modf()
	if modf.Cmp(intZero) == 0 {
		return f
	}
	absf := new(big.Int).Abs(f.x)
	if f.Sign() > 0 {
		return &Fixed{n: f.n, x: new(big.Int).AndNot(absf, modf)}
	}
	return &Fixed{n: f.n, x: new(big.Int).Sub(new(big.Int).Neg(new(big.Int).AndNot(absf, modf)), pow2nInt(f.n))}
}

func (f *Fixed) Ceil() *Fixed {
	modf := f.modf()
	if modf.Cmp(intZero) == 0 {
		return f
	}
	absf := new(big.Int).Abs(f.x)
	if f.Sign() > 0 {
		return &Fixed{n: f.n, x: new(big.Int).Add(new(big.Int).AndNot(absf, modf), pow2nInt(f.n))}
	}
	return &Fixed{n: f.n, x: new(big.Int).Neg(new(big.Int).AndNot(absf, modf))}
}

func (f *Fixed) Abs() *Fixed {
	return &Fixed{n: f.n, x: new(big.Int).Abs(f.x)}
}

func (f *Fixed) Neg() *Fixed {
	return &Fixed{n: f.n, x: new(big.Int).Neg(f.x)}
}

// GobEncode implements the gob.GobEncoder interface.
func (f *Fixed) GobEncode() ([]byte, error) {
	if f == nil {
		return nil, nil
	}
	n, err := f.x.GobEncode()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, len(n)+2)
	binary.BigEndian.PutUint16(buf, f.n)
	copy(buf[2:], n)
	return buf, nil
}

// GobDecode implements the gob.GobDecoder interface.
func (f *Fixed) GobDecode(b []byte) error {
	if len(b) < 2 {
		return errors.New("gob decode: too short")
	}
	f.n = binary.BigEndian.Uint16(b[:2])
	f.x = new(big.Int)
	return f.x.GobDecode(b[2:])
}

// Value implements the driver.Valuer interface.
func (f *Fixed) Value() (driver.Value, error) {
	return f.GobEncode()
}

// Scan implements the sql.Scanner interface.
func (f *Fixed) Scan(v interface{}) error {
	switch v := v.(type) {
	case nil:
		return fmt.Errorf("nil")
	case []byte:
		return f.GobDecode(v)
	case string:
		return f.GobDecode([]byte(v))
	}
	return fmt.Errorf("unsupported type: %T", v)
}

// modf returns a fraction part of the number.
func (f *Fixed) modf() *big.Int {
	return new(big.Int).And(new(big.Int).Abs(f.x), new(big.Int).Sub(pow2nInt(f.n), intOne))
}

func floatToRawInt(n uint16, x *big.Float) *big.Int {
	bi, _ := new(big.Float).Mul(x, pow2nFloat(n)).Int(nil)
	return bi
}

func rawIntToFloat(n uint16, x *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(x), pow2nFloat(n))
}

// pow2nFloat returns a big.Float number which is equal to 2^n.
func pow2nFloat(n uint16) *big.Float {
	return new(big.Float).SetMantExp(floatOne, int(n))
}

// pow2nFloat returns a big.Int number which is equal to 2^n.
func pow2nInt(n uint16) *big.Int {
	return new(big.Int).Lsh(intOne, uint(n))
}
