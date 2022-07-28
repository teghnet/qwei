package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type LogFilter interface {
	// Addresses restricts matches to events created by specific contracts.
	Addresses(ctx context.Context) (AddressList, error)
	// Topics list restricts matches to particular event topics.
	Topics(ctx context.Context) ([][]common.Hash, error)
}

type logFilter struct {
	addresses AddressList
	topics    [][]common.Hash
}

func NewLogFilter(addresses AddressList, topics [][]common.Hash) LogFilter {
	return &logFilter{
		addresses: addresses,
		topics:    topics,
	}
}

func (e logFilter) Addresses(_ context.Context) (AddressList, error) {
	return e.addresses, nil
}

func (e logFilter) Topics(_ context.Context) ([][]common.Hash, error) {
	return e.topics, nil
}
