package ethereum

import (
	"bytes"
	"fmt"

	gethABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ABI interface {
	Pack(name string, args ...interface{}) ([]byte, error)
	Unpack(name string, data []byte) ([]interface{}, error)
	UnpackEvent(name string, topics []common.Hash, data []byte) ([]interface{}, []interface{}, error)
}

type abi struct {
	abi gethABI.ABI
}

func (a *abi) Pack(name string, args ...interface{}) ([]byte, error) {
	return a.abi.Pack(name, args...)
}

func (a *abi) Unpack(name string, data []byte) ([]interface{}, error) {
	return a.abi.Unpack(name, data)
}

func (a *abi) UnpackEvent(name string, topics []common.Hash, data []byte) ([]interface{}, []interface{}, error) {
	if topics[0] != a.abi.Events[name].ID {
		return nil, nil, fmt.Errorf("event signature mismatch")
	}
	var indexed, nonIndexed gethABI.Arguments
	for _, arg := range a.abi.Events[name].Inputs {
		if arg.Indexed {
			arg.Indexed = false
			indexed = append(indexed, arg)
		} else {
			nonIndexed = append(nonIndexed, arg)
		}
	}
	unpackedTopics, err := indexed.Unpack(hashesToBytes(topics[1:]))
	if err != nil {
		return nil, nil, err
	}
	unpackedData, err := nonIndexed.Unpack(data)
	if err != nil {
		return nil, nil, err
	}
	return unpackedTopics, unpackedData, nil
}

func ReadABI(json []byte) (ABI, error) {
	a, err := gethABI.JSON(bytes.NewReader(json))
	if err != nil {
		return nil, err
	}
	for _, m := range a.Methods {
		a.Methods[m.Sig] = m
	}
	return &abi{abi: a}, nil
}

func MustReadABI(json []byte) ABI {
	a, err := ReadABI(json)
	if err != nil {
		panic(err)
	}
	return a
}

func hashesToBytes(hashes []common.Hash) []byte {
	var b []byte
	for _, h := range hashes {
		b = append(b, h[:]...)
	}
	return b
}
