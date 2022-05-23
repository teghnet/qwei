package vars

import (
	"math/big"
)

// KnownEthKeys is a list of publicly known keys
// Never use this keys for anything important - there are bots sweeping these addresses all the time.
var KnownEthKeys = []string{
	"0000000000000000000000000000000000000000000000000000000000000001",
	"1111111111111111111111111111111111111111111111111111111111111111",
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE",
}
var EthKeys []string
var InfuraKeys []string
var AlchemyKeys map[*big.Int]string
var Address string
var Pass string

var AddrBook []string
