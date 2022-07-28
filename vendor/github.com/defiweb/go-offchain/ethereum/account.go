package ethereum

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var ErrReadOnlyAccount = errors.New("unable to sign transaction using read only account")
var ErrMissingAccount = errors.New("unable to find account for requested address")
var ErrMissingWallet = errors.New("unable to find wallet for requested address")
var ErrInvalidTXType = errors.New("invalid TX type, *types.Transaction expected")
var ErrInvalidSigner = errors.New("signer address does not match account address")

var ZeroAddress = common.Address{}
var ZeroAccount = NewReadOnlyAccount(ZeroAddress)

type Account interface {
	Address(ctx context.Context) (common.Address, error)
	SignTX(ctx context.Context, tx interface{}) (interface{}, error)
}

type readOnlyAccount struct {
	address common.Address
}

func HexToAccount(hex string) Account {
	return NewReadOnlyAccount(common.HexToAddress(hex))
}

func NewReadOnlyAccount(address common.Address) Account {
	return &readOnlyAccount{address: address}
}

func (a readOnlyAccount) Address(context.Context) (common.Address, error) {
	return a.address, nil
}

func (readOnlyAccount) SignTX(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, ErrReadOnlyAccount
}

type keystoreAccount struct {
	accountManager *accounts.Manager
	passphrase     string
	address        common.Address
	wallet         accounts.Wallet
	account        *accounts.Account
}

func NewKeystoreAccount(keyStorePath, passphrase string, address common.Address) (Account, error) {
	var err error
	if keyStorePath == "" {
		keyStorePath = defaultKeyStorePath()
	}
	log.Println(keyStorePath)
	a := &keystoreAccount{
		accountManager: accounts.NewManager(
			&accounts.Config{InsecureUnlockAllowed: false},
			keystore.NewKeyStore(keyStorePath, keystore.LightScryptN, keystore.LightScryptP),
		),
		passphrase: passphrase,
		address:    address,
	}
	if a.wallet, a.account, err = a.findAccountByAddress(address); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *keystoreAccount) Address(context.Context) (common.Address, error) {
	return a.address, nil
}

func (a *keystoreAccount) SignTX(ctx context.Context, tx interface{}) (interface{}, error) {
	gethTx, ok := tx.(*types.Transaction)
	if !ok {
		return nil, ErrInvalidTXType
	}
	signedTx, err := a.wallet.SignTxWithPassphrase(*a.account, a.passphrase, gethTx, ChainIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (a *keystoreAccount) findAccountByAddress(from common.Address) (accounts.Wallet, *accounts.Account, error) {
	for _, wallet := range a.accountManager.Wallets() {
		for _, account := range wallet.Accounts() {
			if account.Address == from {
				return wallet, &account, nil
			}
		}
	}
	return nil, nil, ErrMissingAccount
}

type privateKeyAccount struct {
	address common.Address
	key     *ecdsa.PrivateKey
}

func NewPrivateKeyAccount(hexKey string) (Account, error) {
	key, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}
	return &privateKeyAccount{address: crypto.PubkeyToAddress(key.PublicKey), key: key}, nil
}

func (a *privateKeyAccount) Address(context.Context) (common.Address, error) {
	return a.address, nil
}

func (a *privateKeyAccount) SignTX(ctx context.Context, tx interface{}) (interface{}, error) {
	rawTx, ok := tx.(*types.Transaction)
	if !ok {
		return nil, types.ErrInvalidTxType
	}
	opts, err := bind.NewKeyedTransactorWithChainID(a.key, ChainIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	if opts.From != a.address {
		return nil, ErrInvalidSigner
	}
	return opts.Signer(opts.From, rawTx)
}

func defaultKeyStorePath() string {
	var defaultKeyStores []string
	switch runtime.GOOS {
	case "darwin":
		defaultKeyStores = []string{
			os.Getenv("HOME") + "/Library/Ethereum/keystore",
			os.Getenv("HOME") + "/Library/Application Support/io.parity.ethereum/keys/ethereum",
		}
	case "windows":
		defaultKeyStores = []string{
			os.Getenv("APPDATA") + "/Ethereum/keystore",
			os.Getenv("APPDATA") + "/Parity/Ethereum/keys",
		}
	default:
		defaultKeyStores = []string{
			os.Getenv("HOME") + "/.ethereum/keystore",
			os.Getenv("HOME") + "/.local/share/io.parity.ethereum/keys/ethereum",
			os.Getenv("HOME") + "/snap/geth/current/.ethereum/keystore",
			os.Getenv("HOME") + "/snap/parity/current/.local/share/io.parity.ethereum/keys/ethereum",
		}
	}
	for _, keyStore := range defaultKeyStores {
		if _, err := os.Stat(keyStore); !os.IsNotExist(err) {
			return keyStore
		}
	}
	return ""
}

type walletAccount struct {
	account accounts.Account
	wallet  accounts.Wallet
}

func (a *walletAccount) Address(ctx context.Context) (common.Address, error) {
	return a.account.Address, nil
}

func (a *walletAccount) SignTX(ctx context.Context, tx interface{}) (interface{}, error) {
	gethTx, ok := tx.(*types.Transaction)
	if !ok {
		return nil, ErrInvalidTXType
	}
	signedTx, err := a.wallet.SignTx(a.account, gethTx, ChainIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func NewWalletAccount(wallet accounts.Wallet, err error) func(path, passphrase string) (Account, func() error, error) {
	return func(path, passphrase string) (Account, func() error, error) {
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create wallet: %w", err)
		}

		if err := wallet.Open(passphrase); err != nil {
			return nil, nil, fmt.Errorf("cannot open wallet: %w", err)
		}

		dp, err := accounts.ParseDerivationPath(path)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot parse derivation path: %w", err)
		}

		account, err := wallet.Derive(dp, true)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot derive account: %w", err)
		}

		return &walletAccount{account: account, wallet: wallet}, wallet.Close, nil
	}
}

func StdInMnemonicAccount(idx int, path, passphrase string) (account Account, closeFn func() error, err error) {
	return NewWalletAccount(StdinHdWallet(idx))(path, passphrase)
}

func LedgerAccount(idx int, path, passphrase string) (account Account, closeFn func() error, err error) {
	return NewWalletAccount(LedgerWallet(idx))(path, passphrase)
}

func PrivateKey(a Account) (*ecdsa.PrivateKey, error) {
	wa, ok := a.(*walletAccount)
	if !ok {
		return nil, errors.New("wrong type")
	}
	wallet, ok := wa.wallet.(*hdwallet.Wallet)
	if !ok {
		return nil, errors.New("wrong wallet type")
	}
	return wallet.PrivateKey(wa.account)
}
