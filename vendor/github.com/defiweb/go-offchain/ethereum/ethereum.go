package ethereum

import (
	"errors"
	"math/big"
)

var ErrWrongInterface = errors.New("value does not contain the required interface")
var ErrExpectingAddress = errors.New("value does not contain Address")
var ErrExpectingSliceOfInterfaces = errors.New("value does not contain []interface{}")
var ErrExpecting32ByteArray = errors.New("value does not contain [32]byte")
var ErrUnsupportedChain = errors.New("unsupported chain")
var ErrMissingChainID = errors.New("missing chainID")

var (
	MainnetChainID           = big.NewInt(1)
	KovanChainID             = big.NewInt(42)
	RinkebyChainID           = big.NewInt(4)
	GoerliChainID            = big.NewInt(5)
	RopstenChainID           = big.NewInt(3)
	XDAIChainID              = big.NewInt(100)
	PolygonMainnetChainID    = big.NewInt(137)
	PolygonMumbaiChainID     = big.NewInt(80001)
	ArbitrumMainnet          = big.NewInt(42161)
	ArbitrumRinkeby          = big.NewInt(421611)
	OptimismMainnet          = big.NewInt(10)
	OptimismKovan            = big.NewInt(69)
	BinanceSmartChainMainnet = big.NewInt(56)
	BinanceSmartChainTestnet = big.NewInt(97)
)
