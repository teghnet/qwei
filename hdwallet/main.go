package main

import (
	"context"
	"fmt"
	"log"

	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/defiweb/go-offchain/examples/utils"
)

func main() {
	account, closeFn, err := ethereum.StdInMnemonicAccount(0, "0", "")
	if err != nil {
		log.Printf("unable to get account from mnemonic: %s", err)
		account, closeFn, err = ethereum.LedgerAccount(0, "0", "")
		if err != nil {
			log.Fatalf("unable to launch ledger: %s", err)
		}
	}
	defer func() {
		if err := closeFn(); err != nil {
			log.Fatalf("unable to close account: %s", err)
		}
	}()
	fmt.Println(etherscan.Address(ethereum.MainnetChainID, utils.MustAddress(account.Address(context.TODO()))))
}
