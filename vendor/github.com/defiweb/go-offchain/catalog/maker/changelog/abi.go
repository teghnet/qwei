package changelog

import (
	_ "embed"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

//go:embed changelog.abi
var _JSON []byte
var _ABI = ethereum.MustReadABI(_JSON)

// //go:embed changelog.bin
// var _Hex string
// var _Bin = common.FromHex(_Hex)

var Contracts = map[*big.Int]common.Address{
	ethereum.MainnetChainID: common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"),
	ethereum.KovanChainID:   common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"),
	ethereum.GoerliChainID:  common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"),
	ethereum.RinkebyChainID: common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"),
	ethereum.RopstenChainID: common.HexToAddress("0xdA0Ab1e0017DEbCd72Be8599041a2aa3bA7e740F"),
}
