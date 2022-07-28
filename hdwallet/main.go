package main

// func main() {
// 	account, closeFn, err := ethereum.StdInMnemonicAccount(0, "0", "")
// 	if err != nil {
// 		log.Printf("unable to get account from mnemonic: %s", err)
// 		account, closeFn, err = ethereum.LedgerAccount(0, "0", "")
// 		if err != nil {
// 			log.Fatalf("unable to launch ledger: %s", err)
// 		}
// 	}
// 	defer func() {
// 		if err := closeFn(); err != nil {
// 			log.Fatalf("unable to close account: %s", err)
// 		}
// 	}()
// 	address, _ := account.Address(context.TODO())
// 	fmt.Println(etherscan.Address(ethereum.MainnetChainID, address))
// }
