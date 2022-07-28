package ethereum

import "context"

type NonceProvider interface {
	Nonce(ctx context.Context) (uint64, error)
}

type Nonce uint64

func (n Nonce) Nonce(_ context.Context) (uint64, error) {
	return uint64(n), nil
}

type pendingNonce struct {
	client Client
}

func NewPendingNonce(client Client) NonceProvider {
	return &pendingNonce{client: client}
}

func (n *pendingNonce) Nonce(ctx context.Context) (uint64, error) {
	address, err := AccountFromContext(ctx).Address(ctx)
	if err != nil {
		return 0, err
	}
	return n.client.PendingNonce(ctx, address)
}
