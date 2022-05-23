package changelog

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/catalog/maker/changelog"
	"github.com/defiweb/go-offchain/catalog/maker/median"
	"github.com/defiweb/go-offchain/catalog/maker/osm"
	"github.com/defiweb/go-offchain/ethereum"
	"github.com/defiweb/go-offchain/etherscan"
	"github.com/defiweb/go-offchain/examples/utils"
)

func AddressList(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) changelog.AddressMap {
	list, version, err := changelog.ScanChangelog(ctx, client, txParams)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("changelog version\t", version)
	chainID := ethereum.ChainIDFromContext(ctx)
	if chainID.Cmp(ethereum.MainnetChainID) == 0 {
		list[utils.PipStEth] = common.HexToAddress("0x79ED6619640C1c1d9F3E64555172406FE72788B7")
		list[utils.PipWStEth] = common.HexToAddress("0xFe7a2aC0B945f12089aEEB6eCebf4F384D9f043F")
	} else if chainID.Cmp(ethereum.GoerliChainID) == 0 {
		list[utils.PipStEth] = common.HexToAddress("0xc3A7D6C81675C7AEB10B0D25757FB0f64aE0Daab")
		list[utils.PipWStEth] = common.HexToAddress("0x323eac5246d5BcB33d66e260E882fC9bF4B6bf41")
	} else if chainID.Cmp(ethereum.KovanChainID) == 0 {
		list[utils.PipStEth] = common.HexToAddress("0x82CB3F17D1E319ddf34155969B9F350923746830")
		list[utils.PipWStEth] = common.HexToAddress("0xb87E347Ca9AB5f7698521545A9157A48175A6CC2")
	}
	return list
}
func GetFeedsFromOsm(contract common.Address) func(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) ethereum.AddressList {
	return func(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) ethereum.AddressList {
		medianAddr, err := osm.Src(contract).Read(ctx, client, txParams)
		if err != nil {
			log.Println("no src in OSM:", etherscan.Address(ethereum.ChainIDFromContext(ctx), contract), err)
			return nil
		}

		log.Println("scanning median:", etherscan.Address(ethereum.ChainIDFromContext(ctx), medianAddr))
		feeds, err := median.ScanSlots(medianAddr)(ctx, client, txParams)
		if err != nil {
			log.Println("scan failed in Median:", etherscan.Address(ethereum.ChainIDFromContext(ctx), medianAddr), err)
			return nil
		}
		return feeds
	}
}
