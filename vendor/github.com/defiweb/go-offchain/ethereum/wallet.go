package ethereum

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/defiweb/go-offchain/txt"
)

type KeyStoreWallet struct {
	am         *accounts.Manager
	url        accounts.URL
	passphrase string
}

func NewKeystoreWallet() (*KeyStoreWallet, error) {
	keyStorePath := defaultKeyStorePath()
	log.Println(keyStorePath)

	am := accounts.NewManager(
		&accounts.Config{InsecureUnlockAllowed: false},
		keystore.NewKeyStore(keyStorePath, keystore.LightScryptN, keystore.LightScryptP),
	)

	return &KeyStoreWallet{
		am: am,
		url: accounts.URL{
			Scheme: "file",
			Path:   keyStorePath,
		},
	}, nil
}

func (m *KeyStoreWallet) Wallets() []accounts.Wallet {
	return m.am.Wallets()
}

func (m *KeyStoreWallet) Addresses() []common.Address {
	return m.am.Accounts()
}

func (m *KeyStoreWallet) Wallet(address common.Address) (accounts.Wallet, error) {
	for _, wallet := range m.Wallets() {
		for _, account := range wallet.Accounts() {
			if account.Address == address {
				return wallet, nil
			}
		}
	}
	return nil, ErrMissingWallet
}
func (m *KeyStoreWallet) Account(address common.Address) (accounts.Account, error) {
	for _, wallet := range m.Wallets() {
		for _, account := range wallet.Accounts() {
			if account.Address == address {
				return account, nil
			}
		}
	}
	return accounts.Account{}, ErrMissingAccount
}

func (m *KeyStoreWallet) URL() accounts.URL {
	return m.url
}

func (m *KeyStoreWallet) Status() (string, error) {
	return "", nil
}

func (m *KeyStoreWallet) Open(passphrase string) error {
	return nil
}

func (m *KeyStoreWallet) Close() error {
	return m.am.Close()
}

func (m *KeyStoreWallet) Accounts() []accounts.Account {
	var as []accounts.Account
	for _, wallet := range m.am.Wallets() {
		as = append(as, wallet.Accounts()...)
	}
	return as
}

func (m *KeyStoreWallet) Contains(account accounts.Account) bool {
	for _, wallet := range m.am.Wallets() {
		if wallet.Contains(account) {
			return true
		}
	}
	return false
}

func (m *KeyStoreWallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, errors.New("derive operation not supported")
}

func (m *KeyStoreWallet) SelfDerive(bases []accounts.DerivationPath, chain ethereum.ChainStateReader) {
	log.Println("self derive operation not supported")
}

func (m *KeyStoreWallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	w, err := m.am.Find(account)
	if err != nil {
		return nil, err
	}
	return w.SignData(account, mimeType, data)
}

func (m *KeyStoreWallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	w, err := m.am.Find(account)
	if err != nil {
		return nil, err
	}
	return w.SignDataWithPassphrase(account, passphrase, mimeType, data)
}

func (m *KeyStoreWallet) SignText(account accounts.Account, text []byte) ([]byte, error) {
	w, err := m.am.Find(account)
	if err != nil {
		return nil, err
	}
	return w.SignText(account, text)
}

func (m *KeyStoreWallet) SignTextWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error) {
	w, err := m.am.Find(account)
	if err != nil {
		return nil, err
	}
	return w.SignTextWithPassphrase(account, passphrase, hash)
}

func (m *KeyStoreWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	w, err := m.am.Find(account)
	if err != nil {
		return nil, err
	}
	return w.SignTx(account, tx, chainID)
}

func (m *KeyStoreWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	w, err := m.am.Find(account)
	if err != nil {
		return nil, err
	}
	return w.SignTxWithPassphrase(account, passphrase, tx, chainID)
}

func LedgerWallet(idx int) (accounts.Wallet, error) {
	hub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return nil, fmt.Errorf("cannot create wallet hub: %w", err)
	}
	ws := hub.Wallets()
	if len(ws) == 0 {
		return nil, errors.New("no ledger wallets found")
	}
	if idx >= len(ws) {
		return nil, fmt.Errorf("wallet #%d not found", idx)
	}
	return ws[idx], nil
}

func StdinHdWallet(idx int) (accounts.Wallet, error) {
	f, err := txt.NonEmptyStdIn()
	if err != nil {
		return nil, fmt.Errorf("error opening stdin: %w", err)
	}
	ls, err := txt.ReadNonEmptyLines(f, idx+1, false)
	if err != nil {
		return nil, fmt.Errorf("error reading stdin: %w", err)
	}
	m, err := txt.SelectLine(ls, idx)
	if err != nil {
		return nil, fmt.Errorf("cannot select line: %w", err)
	}
	return hdwallet.NewFromMnemonic(m)
}
