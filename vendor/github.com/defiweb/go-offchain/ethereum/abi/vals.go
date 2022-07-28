package abi

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/bn"
)

type ValFn func(any) error
type ValFn2 func(any) (any, error)

func outputArg(target any) ValFn2 {
	switch t := target.(type) {
	case *int:
		return func(src any) (any, error) {
			switch s := src.(type) {
			case *big.Int:
				*t = int(s.Int64())
			default:
				return nil, fmt.Errorf("unsupported source type: %T", src)
			}
			return *t, nil
		}
	case *bn.Int:
		return func(src any) (any, error) {
			switch s := src.(type) {
			case *big.Int:
				*t = *bn.IntFromBigInt(s)
			default:
				return nil, fmt.Errorf("unsupported source type: %T", src)
			}
			return *t, nil
		}
	case *common.Address:
		return func(src any) (any, error) {
			switch s := src.(type) {
			case common.Address:
				*t = s
			default:
				return nil, fmt.Errorf("unsupported source type: %T", src)
			}
			return *t, nil
		}
	case *string:
		return func(src any) (any, error) {
			switch s := src.(type) {
			case [32]byte:
				*t = strings.Trim(string(s[:]), "\x00")
			case string:
				*t = s
			default:
				return nil, fmt.Errorf("unsupported source type: %T", src)
			}
			return *t, nil
		}
	}
	panic(fmt.Errorf("unsupported target type: %T", target).Error())
}

func InputArgs(args ...any) []any {
	for i, arg := range args {
		switch a := arg.(type) {
		case int:
			args[i] = big.NewInt(int64(a))
		}
	}
	return args
}
