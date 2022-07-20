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
)

func AddressList(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) changelog.AddressMap {
	list, version, err := changelog.ScanChangelog(ctx, client, txParams)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("changelog version\t", version)
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
