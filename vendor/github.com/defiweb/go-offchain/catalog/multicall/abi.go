package multicall

import (
	_ "embed"
	"errors"

	"github.com/defiweb/go-offchain/ethereum"
)

var ErrEmptyAddress = errors.New("cannot pack a call without an address")

//go:embed abi.json
var multicallJSON []byte
var multicallABI = ethereum.MustReadABI(multicallJSON)
