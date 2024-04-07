package configs

var Proof = struct {
	PayloadSignatureKey string `env:"TONPROOF_PAYLOAD_SIGNATURE_KEY" envDefault:"secret`
	PayloadLifeTimeSec  int64  `env:"TONPROOF_PAYLOAD_LIFETIME_SEC" envDefault:"300"`
	ProofLifeTimeSec    int64  `env:"TONPROOF_PROOF_LIFETIME_SEC" envDefault:"300"`
	ExampleDomain       string `env:"TONPROOF_EXAMPLE_DOMAIN" envDefault:"fonates.com"`
}{
	PayloadSignatureKey: "secret",
	PayloadLifeTimeSec:  300,
	ProofLifeTimeSec:    300,
	ExampleDomain:       "fonates.com",
}
