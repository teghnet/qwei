package provider

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/defiweb/go-offchain/ethereum"
)

var _ ethereum.Client = (*Alchemy)(nil)
var _ ethereum.ClientProvider = (*Alchemy)(nil)

type Alchemy struct {
	ctx  context.Context
	urls map[*big.Int]string
	keys map[*big.Int]string
	*baseClient
	*clientele
}

// NewAlchemy creates an Alchemy client implementing ethereum.Client interface
func NewAlchemy(apiKeys map[*big.Int]string) *Alchemy {
	p := &Alchemy{
		ctx:       context.Background(),
		urls:      RPCUrls[AlchemyV2],
		keys:      apiKeys,
		clientele: newClientele(),
	}
	p.baseClient = newBaseClient(p.Client)
	return p
}

func (i *Alchemy) Client(chainID *big.Int) (ethereum.Client, error) {
	if chainID == nil {
		return nil, ethereum.ErrMissingChainID
	}
	if !i.hasClient(chainID) {
		return i.provisionClient(chainID)
	}
	return i.getClient(chainID)
}

func (i *Alchemy) provisionClient(chainID *big.Int) (ethereum.Client, error) {
	url, ok := i.urls[chainID]
	if !ok {
		return nil, ethereum.ErrUnsupportedChain
	}
	key, ok := i.keys[chainID]
	if !ok {
		return nil, errors.New("missing API key for chain")
	}
	dial, err := rpc.DialContext(i.ctx, url+key)
	if err != nil {
		return nil, fmt.Errorf("failed to create dialing client: %w", err)
	}
	client := ethereum.NewClient(ethclient.NewClient(dial))
	i.addClient(chainID, client)
	return client, nil
}
