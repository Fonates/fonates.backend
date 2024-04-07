package ton

import "github.com/tonkeeper/tongo/liteapi"

func GetTonNetwork(network string) *liteapi.Client {
	return TonNetworks[network]
}

func GetTonNetworkMainnet() *liteapi.Client {
	net, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		return nil
	}
	return net
}

func GetTonNetworkTestnet() *liteapi.Client {
	net, err := liteapi.NewClientWithDefaultTestnet()
	if err != nil {
		return nil
	}
	return net
}

var TonNetworks = map[string]*liteapi.Client{
	"-239": GetTonNetworkMainnet(),
	"-3":   GetTonNetworkTestnet(),
}
