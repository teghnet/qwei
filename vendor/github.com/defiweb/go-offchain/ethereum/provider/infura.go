package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/defiweb/go-offchain/ethereum"
)

var _ ethereum.Client = (*Infura)(nil)
var _ ethereum.ClientProvider = (*Infura)(nil)

type Infura struct {
	ctx    context.Context
	urls   map[*big.Int]string
	apiKey string
	*baseClient
	*clientele
}

// NewInfura creates an Infura client implementing ethereum.Client interface
func NewInfura(apiKey string) *Infura {
	p := &Infura{
		ctx:       context.Background(),
		urls:      RPCUrls[InfuraV3],
		apiKey:    apiKey,
		clientele: newClientele(),
	}
	p.baseClient = newBaseClient(p.Client)
	return p
}

func (i *Infura) Client(chainID *big.Int) (ethereum.Client, error) {
	if chainID == nil {
		return nil, ethereum.ErrMissingChainID
	}
	if !i.hasClient(chainID) {
		return i.provisionClient(chainID)
	}
	return i.getClient(chainID)
}

func (i *Infura) provisionClient(chainID *big.Int) (ethereum.Client, error) {
	url, ok := i.urls[chainID]
	if !ok {
		return nil, ethereum.ErrUnsupportedChain
	}
	dial, err := rpc.DialContext(i.ctx, url+i.apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create dialing client: %w", err)
	}
	client := ethereum.NewClient(ethclient.NewClient(dial))
	i.addClient(chainID, client)
	return client, nil
}
