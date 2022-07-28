package median

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/defiweb/go-offchain/catalog/multicall"
	"github.com/defiweb/go-offchain/ethereum"
)

func ScanSlots(medianAddr common.Address) func(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (ethereum.AddressList, error) {
	return func(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (ethereum.AddressList, error) {
		var calls []ethereum.Callable
		for i := 0; i < 256; i++ {
			calls = append(calls, Slot(medianAddr, i))
		}
		var list ethereum.AddressList
		responses, err := multicall.Contract().TryAggregate(true, calls...).Read(ctx, client, txParams)
		if err != nil {
			return nil, err
		}
		var t *SlotCall
		for _, v := range responses {
			address, err := t.Values(v.Result)
			if err != nil {
				// log.Println(err, k, address.String())
				continue
			}
			if address != ethereum.ZeroAddress {
				list.Add(address)
			}
		}
		return list, nil
	}
}
