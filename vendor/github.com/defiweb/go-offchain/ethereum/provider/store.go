package provider

import (
	"errors"
	"math/big"

	"github.com/defiweb/go-offchain/ethereum"
)

var ErrClientNotFound = errors.New("client not found")

type clientele struct {
	clients map[*big.Int]ethereum.Client
}

func newClientele() *clientele {
	return &clientele{
		clients: make(map[*big.Int]ethereum.Client),
	}
}

func (cp *clientele) getClient(chainID *big.Int) (ethereum.Client, error) {
	client, ok := cp.clients[chainID]
	if !ok {
		return nil, ErrClientNotFound
	}
	return client, nil
}

func (cp *clientele) addClient(chainID *big.Int, client ethereum.Client) {
	cp.clients[chainID] = client
}

func (cp *clientele) hasClient(chainID *big.Int) bool {
	_, ok := cp.clients[chainID]
	return ok
}
