package ethereum

import "github.com/ethereum/go-ethereum/common"

type LogParams interface {
	FromBlock() (uint64, error)
	BlockHash() (*common.Hash, error)
}

type logParams struct {
	fromBlock uint64
	blockHash *common.Hash
}

func LogsFromBlock(block uint64) LogParams {
	return &logParams{
		fromBlock: block,
	}
}

func LogsByBlockHash(block common.Hash) LogParams {
	return &logParams{
		blockHash: &block,
	}
}

func (l logParams) FromBlock() (uint64, error) {
	return l.fromBlock, nil
}

func (l logParams) BlockHash() (*common.Hash, error) {
	return l.blockHash, nil
}
