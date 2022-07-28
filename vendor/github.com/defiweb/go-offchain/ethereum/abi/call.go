package abi

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/defiweb/go-offchain/ethereum"
)

type CallableUnpacker interface {
	ethereum.AddressProvider
	ethereum.CallDataProvider
	ethereum.Unpacker
	Bind(...any) CallableUnpacker
}
type Callback func([]any) ([]any, error)

type readyToCall struct {
	ethereum.Callable
	outs abi.Arguments
	fns  []any
	cb   Callback
}

func (cu *readyToCall) Unpack(data []byte) (interface{}, error) {
	rets, err := cu.outs.Unpack(data)
	if err != nil {
		return nil, err
	}
	if cu.cb != nil {
		rets, err = cu.cb(rets)
		if err != nil {
			return nil, err
		}
	}
	size := len(rets)
	if l := len(cu.fns); l < size {
		size = l
	}
	for i := 0; i < size; i++ {
		switch f := cu.fns[i].(type) {
		case ValFn:
			if err = f(rets[i]); err != nil {
				return nil, err
			}
		case ValFn2:
			if rets[i], err = f(rets[i]); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported function type: %T", f)
		}
	}
	return rets, nil
}

func (cu *readyToCall) Bind(args ...any) CallableUnpacker {
	for _, arg := range args {
		switch a := arg.(type) {
		case Callback:
			if cu.cb != nil {
				panic("callback already assigned")
			}
			cu.cb = a
		case ValFn:
		case ValFn2:
			cu.fns = append(cu.fns, a)
		default:
			if v := reflect.ValueOf(a); v.Kind() == reflect.Pointer {
				cu.fns = append(cu.fns, outputArg(arg))
			} else {
				panic("non-pointer not supported")
			}
		}
	}
	// TODO: What should we do if the number of bound values does not match num of output args?
	// if len(args) < len(cu.outs) {
	// 	for _, out := range cu.outs[len(args):] {
	// 		cu.fns = append(cu.fns, outputArg(out))
	// 	}
	// }
	return cu
}
