package client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/teghnet/qwei/chain"
)

var ErrUnsupportedChainID = errors.New("unsupported chainID")

func NewInfura(key string) *Infura {
	return &Infura{
		apiKey: key,
	}
}

type Infura struct {
	apiKey  string
	clients map[*big.Int]*ethclient.Client
}

func (c *Infura) Client(chainID *big.Int) (*ethclient.Client, error) {
	if cc, ok := c.clients[chainID]; ok {
		return cc, nil
	}
	prefix, ok := infuraPrefixes[chainID]
	if !ok {
		return nil, ErrUnsupportedChainID
	}
	dial, err := rpc.Dial(fmt.Sprintf("https://%s.infura.io/v3/%s", prefix, c.apiKey))
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(dial)
	c.clients[chainID] = client
	return client, nil
}

var infuraPrefixes = map[*big.Int]string{
	chain.MainnetChainID:        "mainnet",
	chain.KovanChainID:          "kovan",
	chain.RinkebyChainID:        "rinkeby",
	chain.GoerliChainID:         "goerli",
	chain.RopstenChainID:        "ropsten",
	chain.PolygonMainnetChainID: "polygon-mainnet",
	chain.PolygonMumbaiChainID:  "polygon-mumbai",
	chain.ArbitrumMainnet:       "arbitrum-mainnet",
	chain.ArbitrumRinkeby:       "arbitrum-rinkeby",
	chain.OptimismMainnet:       "optimism-mainnet",
	chain.OptimismKovan:         "optimism-kovan",
}
