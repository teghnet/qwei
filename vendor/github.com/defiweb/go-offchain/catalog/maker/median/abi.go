package median

import (
	_ "embed"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

//go:embed internal/median.abi
var _JSON []byte
var _ABI = ethereum.MustReadABI(_JSON)

//go:embed internal/median.bin
var _Hex string
var _Bin = common.FromHex(_Hex)
