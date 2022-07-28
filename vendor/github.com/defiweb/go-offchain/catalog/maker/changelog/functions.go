package changelog

import (
	"github.com/ethereum/go-ethereum/common"
)

type AddressMap map[string]common.Address

// func ScanChangelog(ctx context.Context, client ethereum.Client, txParams ethereum.TXParams) (AddressMap, string, error) {
// 	var ver string
// 	var cnt *bn.Int
// 	{
// 		responses, err := multicall.Contract.TryAggregate(true, Count(), Version()).Read(ctx, client, txParams)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 		for _, v := range responses {
// 			switch t := v.Call.(type) {
// 			case *CountCall:
// 				cnt, _ = t.Values(v.Result)
// 			case *VersionCall:
// 				ver, _ = t.Values(v.Result)
// 			}
// 		}
// 	}
// 	list := make(map[string]common.Address)
// 	{
// 		var calls []ethereum.Callable
// 		for i := uint64(0); i < cnt.Uint64(); i++ {
// 			calls = append(calls, Get(i))
// 		}
// 		responses, err := multicall.Contract.TryAggregate(true, calls...).Read(ctx, client, txParams)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 		var t *GetCall
// 		for _, v := range responses {
// 			key, address, err := t.Values(v.Result)
// 			if err != nil {
// 				log.Println(err, key, address.String())
// 				continue
// 			}
// 			if address != ethereum.ZeroAddress {
// 				list[key] = address
// 			}
// 		}
// 	}
// 	return list, ver, nil
// }

func (m AddressMap) Keys() []string {
	var kk []string
	for s, _ := range m {
		kk = append(kk, s)
	}
	return kk
}

func (m AddressMap) FilterByKeys(f func(string) bool) AddressMap {
	mm := make(AddressMap, 0)
	for k, v := range m {
		if f(k) {
			mm[k] = v
		}
	}
	return mm
}
