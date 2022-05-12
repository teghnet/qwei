package main

import (
	"github.com/teghnet/qwei/chain"
	"github.com/teghnet/qwei/client"
	"github.com/teghnet/qwei/vars"
)

func run() error {
	c, err := client.NewInfura(vars.InfuraAPIKey).Client(chain.MainnetChainID)
	if err != nil {
		return err
	}
	c.Close()

	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
