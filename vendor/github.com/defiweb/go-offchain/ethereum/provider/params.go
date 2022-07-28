package provider

import (
	"math/big"

	"github.com/defiweb/go-offchain/ethereum"
)

const AlchemyV2 = "alchemyV2"
const InfuraV3 = "infuraV3"

var RPCUrls = map[string]map[*big.Int]string{
	InfuraV3: {
		ethereum.MainnetChainID:        "https://mainnet.infura.io/v3/",
		ethereum.KovanChainID:          "https://kovan.infura.io/v3/",
		ethereum.RinkebyChainID:        "https://rinkeby.infura.io/v3/",
		ethereum.GoerliChainID:         "https://goerli.infura.io/v3/",
		ethereum.RopstenChainID:        "https://ropsten.infura.io/v3/",
		ethereum.PolygonMainnetChainID: "https://polygon-mainnet.infura.io/v3/",
		ethereum.PolygonMumbaiChainID:  "https://polygon-mumbai.infura.io/v3/",
		ethereum.ArbitrumMainnet:       "https://arbitrum-mainnet.infura.io/v3/",
		ethereum.ArbitrumRinkeby:       "https://arbitrum-rinkeby.infura.io/v3/",
		ethereum.OptimismMainnet:       "https://optimism-mainnet.infura.io/v3/",
		ethereum.OptimismKovan:         "https://optimism-kovan.infura.io/v3/",
	},
	AlchemyV2: {
		ethereum.MainnetChainID:        "https://eth-mainnet.alchemyapi.io/v2/",
		ethereum.KovanChainID:          "https://eth-kovan.alchemyapi.io/v2/",
		ethereum.RinkebyChainID:        "https://eth-rinkeby.alchemyapi.io/v2/",
		ethereum.GoerliChainID:         "https://eth-goerli.alchemyapi.io/v2/",
		ethereum.RopstenChainID:        "https://eth-ropsten.alchemyapi.io/v2/",
		ethereum.PolygonMainnetChainID: "https://polygon-mainnet.g.alchemy.com/v2/",
		ethereum.PolygonMumbaiChainID:  "https://polygon-mainnet.g.alchemy.com/v2/",
		ethereum.ArbitrumMainnet:       "https://arb-mainnet.g.alchemy.com/v2/",
		ethereum.ArbitrumRinkeby:       "https://arb-rinkeby.g.alchemy.com/v2/",
		ethereum.OptimismMainnet:       "https://opt-mainnet.g.alchemy.com/v2/",
		ethereum.OptimismKovan:         "https://opt-kovan.g.alchemy.com/v2/",
	},
}
