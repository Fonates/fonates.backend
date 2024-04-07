package ton

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"fonates.backend/pkg/configs"
	"github.com/rs/zerolog/log"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
)

const (
	tonProofPrefix   = "ton-proof-item-v2/"
	tonConnectPrefix = "ton-connect"
)

type Domain struct {
	LengthBytes uint32 `json:"lengthBytes"`
	Value       string `json:"value"`
}

type MessageInfo struct {
	Timestamp int64  `json:"timestamp"`
	Domain    Domain `json:"domain"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
	StateInit string `json:"state_init"`
}

type TonProof struct {
	Address string      `json:"address"`
	Network string      `json:"network"`
	Proof   MessageInfo `json:"proof"`
}

type ParsedMessage struct {
	Workchain int32
	Address   []byte
	Timstamp  int64
	Domain    Domain
	Signature []byte
	Payload   string
	StateInit string
}

type Payload struct {
	ExpirtionTime int64
	Signature     string
}

type AccountInfo struct {
	Address struct {
		Bounceable    string `json:"bounceable"`
		NonBounceable string `json:"non_bounceable"`
		Raw           string `json:"raw"`
	} `json:"address"`
	Balance int64  `json:"balance"`
	Status  string `json:"status"`
}

func SignatureVerify(pubkey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pubkey, message, signature)
}

func ConvertTonProofMessage(tp *TonProof) (*ParsedMessage, error) {
	addr, err := tongo.ParseAccountID(tp.Address)
	if err != nil {
		return nil, err
	}

	var parsedMessage ParsedMessage

	sig, err := base64.StdEncoding.DecodeString(tp.Proof.Signature)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	parsedMessage.Workchain = addr.Workchain
	parsedMessage.Address = addr.Address[:]
	parsedMessage.Domain = tp.Proof.Domain
	parsedMessage.Timstamp = tp.Proof.Timestamp
	parsedMessage.Signature = sig
	parsedMessage.Payload = tp.Proof.Payload
	parsedMessage.StateInit = tp.Proof.StateInit
	return &parsedMessage, nil
}

func CreateMessage(message *ParsedMessage) ([]byte, error) {
	wc := make([]byte, 4)
	binary.BigEndian.PutUint32(wc, uint32(message.Workchain))

	ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(ts, uint64(message.Timstamp))

	dl := make([]byte, 4)
	binary.LittleEndian.PutUint32(dl, message.Domain.LengthBytes)
	m := []byte(tonProofPrefix)
	m = append(m, wc...)
	m = append(m, message.Address...)
	m = append(m, dl...)
	m = append(m, []byte(message.Domain.Value)...)
	m = append(m, ts...)
	m = append(m, []byte(message.Payload)...)
	// log.Info(string(m))
	messageHash := sha256.Sum256(m)
	fullMes := []byte{0xff, 0xff}
	fullMes = append(fullMes, []byte(tonConnectPrefix)...)
	fullMes = append(fullMes, messageHash[:]...)
	res := sha256.Sum256(fullMes)
	// log.Info(hex.EncodeToString(res[:]))
	return res[:], nil
}

func CheckProof(ctx context.Context, address tongo.AccountID, net *liteapi.Client, tonProofReq *ParsedMessage) (bool, error) {
	pubKey, err := GetWalletPubKey(ctx, address, net)
	if err != nil {
		if tonProofReq.StateInit == "" {
			log.Error().Msg("get wallet address error")
			return false, err
		}
		if ok, err := CompareStateInitWithAddress(address, tonProofReq.StateInit); err != nil || !ok {
			return ok, err
		}
		pubKey, err = ParseStateInit(tonProofReq.StateInit)
		if err != nil {
			log.Error().Msg("parse wallet state init error")
			return false, err
		}
	}

	if time.Now().After(time.Unix(tonProofReq.Timstamp, 0).Add(time.Duration(configs.Proof.ProofLifeTimeSec) * time.Second)) {
		msgErr := "proof has been expired"
		log.Error().Msg(msgErr)
		return false, fmt.Errorf(msgErr)
	}

	if tonProofReq.Domain.Value != configs.Proof.ExampleDomain {
		msgErr := fmt.Sprintf("wrong domain: %v", tonProofReq.Domain)
		log.Error().Msg(msgErr)
		return false, fmt.Errorf(msgErr)
	}

	mes, err := CreateMessage(tonProofReq)
	if err != nil {
		log.Error().Msg("create message error")
		return false, err
	}

	return SignatureVerify(pubKey, mes, tonProofReq.Signature), nil
}
