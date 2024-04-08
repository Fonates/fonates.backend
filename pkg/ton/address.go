package ton

import "github.com/tonkeeper/tongo"

func AddrFriendly(rawAddress string, net string) (string, error) {
	accountID, err := tongo.ParseAccountID(rawAddress)
	if err != nil {
		return "", err
	}
	isTestnet := net == "-3"
	return accountID.ToHuman(true, isTestnet), nil
}

func ValidateTonAddress(address string) bool {
	_, err := tongo.ParseAddress(address)
	return err == nil
}
