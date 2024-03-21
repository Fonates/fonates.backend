package utils

import (
	"github.com/tonkeeper/tongo"
)

func ValidateTonAddress(address string) bool {
	_, err := tongo.ParseAddress(address)
	return err == nil
}
