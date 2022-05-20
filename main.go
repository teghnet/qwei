package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/defiweb/go-offchain/bn"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/ethereum/provider"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/defiweb/go-offchain/examples/utils"

	"github.com/teghnet/qwei/vars"
)

func SweepEth(ctx context.Context, c ethereum.Client, to ethereum.Account) error {
	fmt.Println(etherscan.Account(ctx, to))
	bal2, err := c.BalanceOf(ctx, utils.MustAddress(to.Address(ctx)))
	if err != nil {
		return err
	}
	fmt.Printf("%24s %s\n", bal2.Wei(), "bal")

	fmt.Println(etherscan.Context(ctx))
	bal, err := c.BalanceOf(ctx, ethereum.AddressFromContext(ctx))
	if err != nil {
		return err
	}
	fmt.Printf("%24s %s\n", bal.Wei(), "balance")

	feeEstimator := ethereum.NewSuggestedFee(c, 1, 1)
	gasPrice, err := feeEstimator.MaxPrice(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%24s %s\n", gasPrice.GWei(), "suggested")

	gasCost := bn.IntFromInt64(ethereum.TransferGas).Mul(gasPrice.Wei())
	fmt.Printf("%24s %s\n", ethereum.Wei(gasCost).Wei(), "min bal")

	gasPrice = ethereum.GWei(bn.FloatFromInt64(4))
	feeEstimator = ethereum.NewStaticFee(gasPrice, gasPrice)
	fmt.Printf("%24s %s\n", gasPrice.GWei(), "price")

	gasCost = bn.IntFromInt64(ethereum.TransferGas).Mul(gasPrice.Wei())
	fmt.Printf("%24s %s\n", ethereum.Wei(gasCost).GWei(), "cost")

	transferAmount := ethereum.Wei(bal.Wei().Sub(gasCost))
	fmt.Printf("%24s %s\n", transferAmount.GWei(), "transfer")

	if transferAmount.Wei().Cmp(bn.IntFromInt64(0)) <= 0 {
		return errors.New("not enough funds")
	}

	nonce := ethereum.NewPendingNonce(c)
	n, err := nonce.Nonce(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%24d %s\n", n, "nonce")

	if false {
		txParams := ethereum.NewTXParams(
			transferAmount,
			nonce,
			ethereum.GasLimit(ethereum.TransferGas),
			feeEstimator,
		)
		hash, err := c.Transfer(ctx, utils.MustAddress(to.Address(ctx)), txParams)
		if err != nil {
			return err
		}
		fmt.Println(etherscan.Txx(ctx, *hash))
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx = ethereum.WithChainID(ctx, ethereum.MainnetChainID)

	sourceAccount, err := ethereum.NewPrivateKeyAccount(vars.EthKeys[4])
	if err != nil {
		log.Println(err)
		return
	}
	ctx = ethereum.WithAccount(ctx, sourceAccount)

	client := provider.NewAlchemy(vars.AlchemyKeys)
	defer func() {
		err := client.Close(ctx)
		if err != nil {
			log.Printf("error closing client: %s", err)
		}
	}()

	destinationAccount, err := ethereum.NewPrivateKeyAccount(vars.EthKeys[0])
	if err != nil {
		log.Println(err)
		return
	}

	err = SweepEth(ctx, client, destinationAccount)
	if err != nil {
		log.Println(err)
		return
	}
}
