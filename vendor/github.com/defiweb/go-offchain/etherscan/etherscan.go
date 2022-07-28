package etherscan

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/ethereum"
)

func Context(ctx context.Context) string {
	address, err := ethereum.AccountFromContext(ctx).Address(ctx)
	if err != nil {
		panic(err)
	}
	return url(ethereum.ChainIDFromContext(ctx), "address", address.String())
}
func Account(ctx context.Context, acc ethereum.AddressProvider) string {
	address, err := acc.Address(ctx)
	if err != nil {
		panic(err)
	}
	return url(ethereum.ChainIDFromContext(ctx), "address", address.String())
}

func Address(chain *big.Int, address common.Address) string {
	return url(chain, "address", address.String())
}

func Txx(ctx context.Context, hash common.Hash) string {
	return url(ethereum.ChainIDFromContext(ctx), "tx", hash.String())
}
func Tx(chain *big.Int, hash common.Hash) string {
	return url(chain, "tx", hash.String())
}

func url(chain *big.Int, r string, s string) string {
	if chain.Cmp(ethereum.MainnetChainID) == 0 {
		return fmt.Sprintf("https://etherscan.io/%s/%s", r, s)
	}
	p, ok := prefix[chain]
	if !ok {
		return s
	}
	return fmt.Sprintf("https://%s.etherscan.io/%s/%s", p, r, s)
}

var prefix = map[*big.Int]string{
	ethereum.KovanChainID:   "kovan",
	ethereum.RinkebyChainID: "rinkeby",
	ethereum.GoerliChainID:  "goerli",
	ethereum.RopstenChainID: "ropsten",
}
