package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
)

type AddressList []common.Address

func NewAddressList(addrs ...string) AddressList {
	var l AddressList
	l.AppendHex(addrs...)
	return l
}

func (c *AddressList) Add(addr common.Address) {
	for _, a := range *c {
		if a == addr {
			return
		}
	}
	*c = append(*c, addr)
}

func (c *AddressList) Append(addrs ...common.Address) {
	for _, a := range addrs {
		c.Add(a)
	}
}

func (c *AddressList) AppendHex(addrs ...string) {
	for _, a := range addrs {
		c.Add(common.HexToAddress(a))
	}
}

func (c *AddressList) WithoutHex(addrs ...string) AddressList {
	var list AddressList
	list.AppendHex(addrs...)
	return c.Without(list...)
}

func (c *AddressList) Without(addrs ...common.Address) AddressList {
	return c.Filter(Except(addrs...))
}

func (c *AddressList) Filter(f func(address common.Address) bool) AddressList {
	var list AddressList
	for _, a := range *c {
		if f(a) {
			list.Add(a)
		}
	}
	return list
}

func Except(addrs ...common.Address) func(address common.Address) bool {
	m := make(map[common.Address]byte)
	for _, a := range addrs {
		m[a] = 1
	}
	return func(a common.Address) bool {
		_, is := m[a]
		return !is
	}
}
