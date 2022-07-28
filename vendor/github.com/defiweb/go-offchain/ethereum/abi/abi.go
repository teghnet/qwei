package abi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

// Method is something you call while it is already attached to a specific contract
type Method func(...interface{}) (CallableUnpacker, error)

func ParseWithAddr(a common.Address, s string) (Method, error) {
	f, err := Parse(s)
	if err != nil {
		return nil, err
	}
	return f.Attach(a), nil
}

// Function is something that has not been attached to a contract. The first arg is the receiver.
type Function func(common.Address, ...interface{}) (CallableUnpacker, error)

func (f Function) Attach(addr common.Address) Method {
	return func(args ...interface{}) (CallableUnpacker, error) {
		return f(addr, args...)
	}
}

func Parse(sig string) (Function, error) {
	parsed, err := parse(sig)
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal([]SignatureMarshaling{parsed})
	if err != nil {
		return nil, err
	}

	a, err := abi.JSON(bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	if len(a.Methods) != 1 {
		return nil, fmt.Errorf("abi needs to have exactly one method - %d found", len(a.Methods))
	}

	for _, m := range a.Methods {
		return function(m.ID, m.Inputs, m.Outputs), nil
	}

	return nil, errors.New("failed to create method abi")
}

var function = func(ID []byte, ins, outs abi.Arguments) Function {
	return func(addr common.Address, args ...any) (CallableUnpacker, error) {
		packed, err := ins.Pack(InputArgs(args...)...)
		if err != nil {
			return nil, err
		}
		return &readyToCall{
			outs: outs,
			Callable: ethereum.NewCall(
				ethereum.StaticAddress(addr),
				ethereum.CallData(append(ID, packed...)),
			),
		}, nil
	}
}
