package ethereum

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	eth "github.com/ethereum/go-ethereum"
	ethABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/defiweb/go-offchain/bn"
)

type Client interface {
	// BalanceOf returns the balance of the given address for the latest block
	// or for the block number from the context.
	BalanceOf(ctx context.Context, owner common.Address) (*Value, error)
	// Transfer transfers the given amount of ether to the given address. It
	// uses an account from the context.
	Transfer(ctx context.Context, to common.Address, txParams TXParams) (*common.Hash, error)
	// Cancel cancels transaction with the given nonce. It uses an account from
	// the context. If the nonce is nil, then the pending nonce is used. If the
	// fee is nil, then the suggested fee is used.
	Cancel(ctx context.Context, nonce NonceProvider, fee FeeEstimator) (*common.Hash, error)
	// Receipt returns the receipt of the given transaction for the latest
	// block or for the block number from the context.
	Receipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// Nonce returns the nonce of the given address for the latest block or for
	// the block number from the context.
	Nonce(ctx context.Context, address common.Address) (uint64, error)
	// PendingNonce returns the pending nonce of the given address for the
	// latest block or for the block number from the context.
	PendingNonce(ctx context.Context, address common.Address) (uint64, error)
	// GasPrice returns the current gas price. It always returns the gas price
	// for the current network conditions.
	GasPrice(ctx context.Context) (*Value, error)
	// TipValue returns the current tip value. It always returns the tip value
	// for the current network conditions.
	TipValue(ctx context.Context) (*Value, error)
	// Block returns the latest block. If a block number is set in the context,
	// then it returns the block with the that number.
	Block(ctx context.Context) (*types.Block, error)
	// BlockByHash returns the block with the given hash.
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	// BlockNumber returns the latest block number. If a block number is set in
	// the context, then it returns the block number from the context.
	BlockNumber(ctx context.Context) (uint64, error)
	// EstimateGas estimates the gas cost of a transaction.
	EstimateGas(ctx context.Context, call Callable) (uint64, error)
	// The Logs function returns Ethereum logs that match the specified filter.
	// If block hash is set in logParams, then it returns logs from that block.
	// If block number is set in logParams, then it returns logs from that
	// block number to the latest or to the block number from the context.
	Logs(ctx context.Context, filter LogFilter, logParams LogParams) ([]types.Log, error)
	// Storage returns the storage value at the given address for the latest
	// block or for the block number from the context.
	Storage(ctx context.Context, address common.Address, key common.Hash) ([]byte, error)
	// NetworkID returns the network ID.
	NetworkID(ctx context.Context) (*big.Int, error)
	// Read executes the given callable on the blockchain. The callable is
	// executed on the latest block or for the block number from the context.
	// If Callable implements the Unpacker interface, then the result of the
	// callable is unpacked before it is returned. If txParams is nil, then
	// transaction parameters are set automatically.
	Read(ctx context.Context, call Callable, txParams TXParams) (interface{}, error)
	// Write sends a transaction to the blockchain. If txParams is nil, then
	// transaction parameters are set automatically.
	Write(ctx context.Context, call Callable, txParams TXParams) (*common.Hash, error)
	// Close closes the client.
	Close(ctx context.Context) error
}

type ClientProvider interface {
	Client(chainID *big.Int) (Client, error)
}

type ethClient interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	FilterLogs(ctx context.Context, q eth.FilterQuery) ([]types.Log, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	StorageAt(ctx context.Context, account common.Address, key common.Hash, block *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call eth.CallMsg, block *big.Int) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	EstimateGas(ctx context.Context, msg eth.CallMsg) (uint64, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	Close()
}

// ErrRevert may be returned by client.Call method in case of EVM revert.
type ErrRevert struct {
	Message string
	Err     error
}

func (e ErrRevert) Error() string {
	return fmt.Sprintf("reverted: %s", e.Message)
}

func (e ErrRevert) Unwrap() error {
	return e.Err
}

// client implements the Client interface.
type client struct {
	ethClient ethClient
}

// NewClient returns a new Client instance.
func NewClient(c ethClient) Client {
	return &client{ethClient: c}
}

// NewRPCClient returns a new RPC client instance.
func NewRPCClient(url string) (Client, error) {
	dial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return NewClient(ethclient.NewClient(dial)), nil
}

// BalanceOf implements the Client interface.
func (e *client) BalanceOf(ctx context.Context, address common.Address) (*Value, error) {
	balance, err := e.ethClient.BalanceAt(ctx, address, blockNumberFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return Wei(bn.IntFromBigInt(balance)), nil
}

// Transfer implements the Client interface.
func (e *client) Transfer(ctx context.Context, to common.Address, txParams TXParams) (*common.Hash, error) {
	return e.Write(ctx, NewTransferCall(StaticAddress(to)), txParams)
}

// Cancel implements the Client interface.
func (e *client) Cancel(ctx context.Context, nonce NonceProvider, fee FeeEstimator) (*common.Hash, error) {
	if nonce == nil {
		nonce = NewPendingNonce(e)
	}
	if fee == nil {
		fee = NewSuggestedFee(e, defaultTipValueMultiplier, defaultMaxPriceMultiplier)
	}
	address, err := AccountFromContext(ctx).Address(ctx)
	if err != nil {
		return nil, err
	}
	return e.Transfer(
		ctx,
		address,
		NewTXParams(Wei(bn.IntFromInt64(0)), nonce, GasLimit(TransferGas), fee),
	)
}

// Receipt implements the Client interface.
func (e *client) Receipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return e.ethClient.TransactionReceipt(ctx, txHash)
}

// Nonce implements the Client interface.
func (e *client) Nonce(ctx context.Context, address common.Address) (uint64, error) {
	return e.ethClient.NonceAt(ctx, address, blockNumberFromContext(ctx))
}

// PendingNonce implements the Client interface.
func (e *client) PendingNonce(ctx context.Context, address common.Address) (uint64, error) {
	return e.ethClient.PendingNonceAt(ctx, address)
}

// GasPrice implements the Client interface.
func (e *client) GasPrice(ctx context.Context) (*Value, error) {
	price, err := e.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return Wei(bn.IntFromBigInt(price)), nil
}

// TipValue implements the Client interface.
func (e *client) TipValue(ctx context.Context) (*Value, error) {
	price, err := e.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	return Wei(bn.IntFromBigInt(price)), nil
}

// Block implements the Client interface.
func (e *client) Block(ctx context.Context) (*types.Block, error) {
	return e.ethClient.BlockByNumber(ctx, blockNumberFromContext(ctx))
}

// BlockByHash implements the Client interface.
func (e *client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return e.ethClient.BlockByHash(ctx, hash)
}

func (e *client) BlockNumber(ctx context.Context) (uint64, error) {
	block := BlockNumberFromContext(ctx)
	if block > 0 {
		return block, nil
	}
	return e.ethClient.BlockNumber(ctx)
}

// EstimateGas implements the Client interface.
func (e *client) EstimateGas(ctx context.Context, call Callable) (uint64, error) {
	if call == nil {
		return 0, errors.New("call cannot be empty")
	}
	from := AccountFromContext(ctx)
	to, err := NilWhenZero(call.Address(ctx))
	if err != nil {
		return 0, err
	}
	data, err := call.Data(ctx)
	if err != nil {
		return 0, err
	}
	address, err := from.Address(ctx)
	if err != nil {
		return 0, err
	}
	gas, err := e.ethClient.EstimateGas(ctx, eth.CallMsg{
		From: address,
		To:   to,
		Data: data,
	})
	if err != nil {
		return 0, err
	}
	return gas, nil
}

// Logs implements the Client interface.
func (e *client) Logs(ctx context.Context, filter LogFilter, logParams LogParams) ([]types.Log, error) {
	var err error
	addresses, err := filter.Addresses(ctx)
	if err != nil {
		return nil, err
	}
	topics, err := filter.Topics(ctx)
	if err != nil {
		return nil, err
	}
	blockHash, err := logParams.BlockHash()
	if err != nil {
		return nil, err
	}
	fromBlock, err := logParams.FromBlock()
	if err != nil {
		return nil, err
	}
	f := eth.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
	}
	if blockHash == nil {
		f.FromBlock = new(big.Int).SetUint64(fromBlock)
		f.ToBlock = blockNumberFromContext(ctx)
	} else {
		f.BlockHash = blockHash
	}
	logs, err := e.ethClient.FilterLogs(ctx, f)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Storage implements the Client interface.
func (e *client) Storage(ctx context.Context, address common.Address, key common.Hash) ([]byte, error) {
	return e.ethClient.StorageAt(ctx, address, key, blockNumberFromContext(ctx))
}

// NetworkID implements the Client interface.
func (e *client) NetworkID(ctx context.Context) (*big.Int, error) {
	return e.ethClient.NetworkID(ctx)
}

// TransactionReceipt implements the Client interface.
func (e *client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return e.ethClient.TransactionReceipt(ctx, txHash)
}

// Read implements the Client interface.
func (e *client) Read(ctx context.Context, call Callable, txParams TXParams) (interface{}, error) {
	if call == nil {
		return nil, errors.New("call cannot be empty")
	}
	if txParams == nil {
		txParams = NewTXParams(nil, nil, nil, nil)
	}
	tx, err := e.newTxData(ctx, call, txParams)
	if err != nil {
		return nil, err
	}
	from, err := tx.account.Address(ctx)
	if err != nil {
		return nil, err
	}
	cm := eth.CallMsg{
		From:      from,
		To:        tx.address,
		Gas:       tx.gasLimit,
		GasFeeCap: tx.gasFeeCap,
		GasTipCap: tx.gasTipCap,
		Value:     tx.amount.Wei().BigInt(),
		Data:      tx.data,
	}
	res, err := e.ethClient.CallContract(ctx, cm, blockNumberFromContext(ctx))
	if err := isRevertErr(err); err != nil {
		return nil, err
	}
	if err := isRevertResp(res); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if u, ok := call.(Unpacker); ok {
		return u.Unpack(res)
	}
	return res, err
}

// Write implements the Client interface.
func (e *client) Write(ctx context.Context, call Callable, txParams TXParams) (*common.Hash, error) {
	if call == nil {
		return nil, errors.New("call cannot be empty")
	}
	if txParams == nil {
		txParams, _ = NewAutoTXParams(e, nil)
	}
	tx, err := e.newTxData(ctx, call, txParams)
	if err != nil {
		return nil, err
	}
	if tx.nonce, err = txParams.Nonce(ctx); err != nil {
		return nil, err
	}
	stx, err := tx.account.SignTX(ctx, types.NewTx(&types.DynamicFeeTx{
		ChainID:   tx.chainID,
		Nonce:     tx.nonce,
		Gas:       tx.gasLimit,
		GasFeeCap: tx.gasFeeCap,
		GasTipCap: tx.gasTipCap,
		To:        tx.address,
		Value:     tx.amount.Wei().BigInt(),
		Data:      tx.data,
	}))
	if err != nil {
		return nil, err
	}
	hash := stx.(*types.Transaction).Hash()
	return &hash, e.ethClient.SendTransaction(ctx, stx.(*types.Transaction))
}

func (e *client) Close(_ context.Context) error {
	e.ethClient.Close()
	return nil
}

type txData struct {
	account   Account
	chainID   *big.Int
	address   *common.Address
	amount    *Value
	nonce     uint64
	gasTipCap *big.Int
	gasFeeCap *big.Int
	gasLimit  uint64
	data      CallData
}

func (e *client) newTxData(ctx context.Context, call Callable, txParams TXParams) (*txData, error) {
	var err error
	tx := &txData{}
	tx.account = AccountFromContext(ctx)
	if tx.chainID = ChainIDFromContext(ctx); tx.chainID == nil {
		if tx.chainID, err = e.NetworkID(ctx); err != nil {
			return nil, err
		}
		ctx = WithChainID(ctx, tx.chainID)
	}
	if tx.address, err = NilWhenZero(call.Address(ctx)); err != nil {
		return nil, err
	}
	if tx.amount, err = txParams.Amount(ctx); err != nil {
		return nil, err
	}
	if tx.data, err = call.Data(ctx); err != nil {
		return nil, err
	}
	tipValue, err := txParams.TipValue(ctx)
	if err != nil {
		return nil, err
	}
	if tipValue != nil {
		tx.gasTipCap = tipValue.Wei().BigInt()
	}
	maxPrice, err := txParams.MaxPrice(ctx)
	if err != nil {
		return nil, err
	}
	if maxPrice != nil {
		tx.gasFeeCap = maxPrice.Wei().BigInt()
	}
	if tx.gasLimit, err = txParams.GasLimit(ctx, call); err != nil {
		return nil, err
	}
	return tx, nil
}

func isRevertResp(res []byte) error {
	revert, err := ethABI.UnpackRevert(res)
	if err != nil {
		return nil
	}
	return ErrRevert{Message: revert, Err: nil}
}

var revertRE = regexp.MustCompile("(0x[a-zA-Z0-9]+)")

func isRevertErr(vmErr error) error {
	if terr, is := vmErr.(rpc.DataError); is {
		// Some RPC servers returns "revert" data as a hex encoded string, here
		// we are trying to parse it:
		if str, ok := terr.ErrorData().(string); ok {
			match := revertRE.FindStringSubmatch(str)
			if len(match) == 2 && len(match[1]) > 2 {
				bytes, err := hex.DecodeString(match[1][2:])
				if err != nil {
					return nil
				}
				revert, err := ethABI.UnpackRevert(bytes)
				if err != nil {
					return nil
				}
				return ErrRevert{Message: revert, Err: vmErr}
			}
		}
	}
	return nil
}

func blockNumberFromContext(ctx context.Context) *big.Int {
	b := BlockNumberFromContext(ctx)
	if b == 0 {
		return nil
	}
	return new(big.Int).SetUint64(b)
}

var _ Client = (*client)(nil)
