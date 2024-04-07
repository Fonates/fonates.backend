package configs

import (
	"github.com/caarlos0/env/v10"
)

func LoadConfig() error {
	if err := env.Parse(&Proof); err != nil {
		return err
	}
	return nil
}
