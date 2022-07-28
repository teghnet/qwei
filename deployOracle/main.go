package main

// import (
// 	"context"
// 	_ "embed"
// 	"errors"
// 	"fmt"
// 	"log"
//
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/crypto"
//
// 	"github.com/defiweb/go-offchain/bn"
// 	"github.com/defiweb/go-offchain/catalog/maker/com"
// 	"github.com/defiweb/go-offchain/catalog/maker/median"
// 	"github.com/defiweb/go-offchain/catalog/maker/osm"
// 	"github.com/defiweb/go-offchain/ethereum"
// 	"github.com/defiweb/go-offchain/ethereum/provider"
// 	"github.com/defiweb/go-offchain/etherscan"
// 	"github.com/defiweb/go-offchain/examples"
// 	"github.com/defiweb/go-offchain/examples/utils"
// 	"github.com/teghnet/qwei/changelog"
// )
//
// func main() {
// 	var read, write, wait, doDeployOracle, doRelyAdmins, doDenyAdmins, doPermissions, doFeeds bool
// 	read = true
// 	// write = true
// 	// wait = true
// 	// doDeployOracle = true // create actual contracts (for mainnet do it first)
// 	// doFeeds = true       // add feeds (needs already deployed contracts)
// 	// doPermissions = true // kiss and rely on the protocol
// 	// doRelyAdmins = true // rely additional team admins in case deployer is unavailable
// 	// doDenyAdmins = true
//
// 	pip := "PIP_WSTETH"
// 	chainID := ethereum.GoerliChainID
//
// 	ctx := context.Background()
// 	ctx = ethereum.WithChainID(ctx, chainID)
//
// 	account := utils.MustAccount(ethereum.NewPrivateKeyAccount(examples.AccountPrivKey))
// 	ctx = ethereum.WithAccount(ctx, account)
//
// 	address, err := account.Address(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println("from:", etherscan.Address(chainID, address))
// 	client := provider.NewInfura(examples.InfuraAPIKey)
//
// 	nonce := utils.MustUint64(client.PendingNonce(ctx, address))
// 	log.Println("first nonce:", nonce)
// 	txParams := incrementingNonce(client, 1, 80, nonce)
//
// 	list := changelog.AddressList(ctx, client, txParams)
// 	gas := uint64(0)
//
// 	if doDeployOracle {
// 		log.Println("----- Create Median")
//
// 		medianAddress := crypto.CreateAddress(address, txParams.nonce)
// 		log.Println("median", etherscan.Address(chainID, medianAddress))
// 		gas += utils.Run(read, write, wait)(median.Constructor())(ctx, client, txParams)
//
// 		list["MEDIAN_"+pip] = medianAddress
//
// 		log.Println("----- Create OSM")
//
// 		osmAddress := crypto.CreateAddress(address, txParams.nonce)
// 		log.Println("osm", etherscan.Address(chainID, osmAddress))
// 		gas += utils.Run(read, write, wait)(osm.Constructor(medianAddress))(ctx, client, txParams)
//
// 		list[pip] = osmAddress
// 	}
//
// 	osmAddress := utils.GetContract(list, pip, ctx)
// 	ma, err := osm.Src(osmAddress).Read(ctx, client, txParams)
// 	if err != nil {
// 		log.Fatalln("no src in OSM:", etherscan.Address(chainID, osmAddress), err)
// 	}
// 	list["MEDIAN_"+pip] = ma
// 	medianAddress := utils.GetContract(list, "MEDIAN_"+pip, ctx)
//
// 	if doFeeds {
// 		log.Println("----- Setup Feeds")
//
// 		if chainID.Cmp(ethereum.MainnetChainID) == 0 {
// 			gas += utils.Run(read, write, wait)(median.SetBar(medianAddress, 13))(ctx, client, txParams)
// 		} else {
// 			gas += utils.Run(read, write, wait)(median.SetBar(medianAddress, 1))(ctx, client, txParams)
// 		}
//
// 		var requiredFeeds ethereum.AddressList
// 		if chainID.Cmp(ethereum.MainnetChainID) == 0 {
// 			requiredFeeds = changelog.GetFeedsFromOsm(utils.GetContract(list, "PIP_ETH", ctx))(ctx, client, txParams)
// 		} else {
// 			// requiredFeeds.AppendHex(examples.TestFeeds...)
// 		}
//
// 		currentFeeds := changelog.GetFeedsFromOsm(osmAddress)(ctx, client, txParams)
//
// 		if dropFeeds := currentFeeds.Without(requiredFeeds...); len(dropFeeds) > 0 {
// 			fmt.Printf("drop: %#v\n", dropFeeds)
// 			gas += utils.Run(read, write, wait)(median.Drop(medianAddress, dropFeeds...))(ctx, client, txParams)
// 		}
//
// 		if liftFeeds := requiredFeeds.Without(currentFeeds...); len(liftFeeds) > 0 {
// 			fmt.Printf("lift: %#v\n", liftFeeds)
// 			gas += utils.Run(read, write, wait)(median.Lift(medianAddress, liftFeeds...))(ctx, client, txParams)
// 		}
// 	}
//
// 	if doPermissions {
// 		log.Println("----- Set Params & Permissions")
//
// 		mcdPauseProxy := utils.GetContract(list, "MCD_PAUSE_PROXY", ctx)
//
// 		gas += utils.Run(read, write, wait)(com.Kiss(medianAddress, osmAddress))(ctx, client, txParams)
// 		// if wait median.Bud(medianAddress,mcdSpot) else add to checklist
// 		gas += utils.Run(read, write, wait)(com.Rely(medianAddress, mcdPauseProxy))(ctx, client, txParams)
// 		// if wait median.Ward(medianAddress,mcdPauseProxy) else add to checklist
//
// 		mcdSpot := utils.GetContract(list, "MCD_SPOT", ctx)
//
// 		gas += utils.Run(read, write, wait)(com.Kiss(osmAddress, mcdSpot))(ctx, client, txParams)
// 		// if wait osm.Bud(osmAddress,mcdSpot) else add to checklist
// 		gas += utils.Run(read, write, wait)(com.Rely(osmAddress, mcdPauseProxy))(ctx, client, txParams)
// 		// if wait osm.Ward(osmAddress,mcdPauseProxy) else add to checklist
//
// 		if chainID.Cmp(ethereum.MainnetChainID) != 0 {
// 			gas += utils.Run(read, write, wait)(median.SetBar(medianAddress, 1))(ctx, client, txParams)
// 		}
// 	}
//
// 	adminAddress := common.HexToAddress("0x4D6fbF888c374D7964D56144dE0C0cFBd49750D3")
// 	adminAddress0 := common.HexToAddress("0x1f42e41A34B71606FcC60b4e624243b365D99745")
// 	adminAddress1 := common.HexToAddress("0xd00Af2385c4BE41C9E73B0F8A0F189b962fdd18d")
// 	adminAddress2 := common.HexToAddress("0xc3c632b6d38Bd83241e1695c17Cd73B88D9d3Ae3")
//
// 	admins := []common.Address{adminAddress}
// 	if chainID.Cmp(ethereum.MainnetChainID) == 0 {
// 	} else {
// 		admins = append(admins, adminAddress0, adminAddress1, adminAddress2)
// 	}
// 	contracts := []common.Address{
// 		medianAddress,
// 		// osmAddress,
// 	}
//
// 	if doRelyAdmins {
// 		log.Println("----- Rely Admins")
//
// 		for _, a := range admins {
// 			for _, c := range contracts {
// 				if !isWard(c, a, ctx, client, txParams) {
// 					log.Println(etherscan.Address(chainID, c), "rely(", etherscan.Address(chainID, a), ")")
// 					gas += utils.Run(read, write, wait)(com.Rely(c, a))(ctx, client, txParams)
// 				}
// 			}
// 		}
// 	}
//
// 	if doDenyAdmins {
// 		log.Println("----- Deny Admins")
//
// 		for idx, a := range putCurrentAddrAtTheEnd(admins, address) {
// 			if a == address && idx < len(admins)-1 {
// 				continue
// 			}
// 			for _, c := range contracts {
// 				if isWard(c, a, ctx, client, txParams) {
// 					log.Println(etherscan.Address(chainID, c), "deny(", etherscan.Address(chainID, a), ")")
// 					gas += utils.Run(read, write, wait)(com.Deny(c, a))(ctx, client, txParams)
// 				}
// 			}
// 		}
// 	}
//
// 	utils.PrintGas("total:", gas, ctx, txParams)
// 	log.Println("next nonce:", txParams.nonce)
// }
//
// func putCurrentAddrAtTheEnd(admins []common.Address, address common.Address) []common.Address {
// 	var admins2 []common.Address
// 	last := common.Address{}
// 	for _, a := range admins {
// 		if a == address {
// 			last = a
// 			continue
// 		}
// 		admins2 = append(admins2, a)
// 	}
// 	admins = append(admins2, last)
// 	return admins
// }
//
// type txParams2 struct {
// 	client   ethereum.Client
// 	nonce    uint64
// 	tip, max float64
// }
//
// func incrementingNonce(client ethereum.Client, tip, max float64, nonce uint64) *txParams2 {
// 	return &txParams2{
// 		client: client,
// 		nonce:  nonce,
// 		tip:    tip,
// 		max:    max,
// 	}
// }
//
// func (t *txParams2) Amount(_ context.Context) (*ethereum.Value, error) {
// 	return nil, nil
// }
//
// func (t *txParams2) Nonce(_ context.Context) (uint64, error) {
// 	nonce := t.nonce
// 	t.nonce += 1
// 	log.Println("using nonce:", nonce)
// 	return nonce, nil
// }
//
// func (t *txParams2) GasLimit(ctx context.Context, call ethereum.Callable) (uint64, error) {
// 	if call == nil {
// 		return 0, errors.New("call cannot be empty")
// 	}
// 	data, err := call.Data(ctx)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if data == nil {
// 		return ethereum.TransferGas, nil
// 	}
// 	v, err := t.client.EstimateGas(ctx, call)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint64(float64(v) * 1.5), nil
// }
//
// func (t *txParams2) TipValue(_ context.Context) (*ethereum.Value, error) {
// 	return ethereum.GWei(bn.FloatFromFloat64(t.tip)), nil
// }
//
// func (t *txParams2) MaxPrice(_ context.Context) (*ethereum.Value, error) {
// 	return ethereum.GWei(bn.FloatFromFloat64(t.max)), nil
// }
//
// func isWard(c, a common.Address, ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) bool {
// 	w, err := com.Ward(c, a).Read(ctx, client, txParams)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return w
// }
