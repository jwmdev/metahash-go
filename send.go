package metahash

import (
	"encoding/hex"
	"errors"
	"log"
	"math/big"
	"strconv"
	"strings"
)

const proxUrl = "http://net-main.MetaKey.com:9999"

type MetaTx struct {
	To    string `json:"to"`
	Value int64  `json:"value"`
	Fee   int64  `json:"fee"`
	Data  string `json:"data"`
	Nonce int64  `json:"nonce"`
}

//TransactionArgs argument
type MetaTxArg struct {
	To     string `json:"to"`
	Value  string `json:"value"`
	Fee    string `json:"fee"`
	Data   string `json:"data"`
	Nonce  string `json:"nonce"`
	Pubkey string `json:"pubkey"`
	Sign   string `json:"sign"`
}

//TransactionResponse response
type MetaTxResp struct {
	Result string
	Params string
	Error  string
}

type MetaKey struct {
	Key Key //private key in ex format
}

//Iniate a new wallet by supplying the wallet address and private key
func InitWallet(walletAddress, privateKey string) (*MetaKey, error) {
	key, err := createKey(PrivateKey(privateKey))
	if err != nil {
		return nil, err
	}
	key.SetAddress(walletAddress)
	mk := &MetaKey{
		Key: key,
	}
	return mk, nil
}

func (mk *MetaKey) Transfer(toAddress string, sendAmount float64) (string, error) {
	//get account balance

	bal, err := FetchBalance(string(mk.Key.Address()))
	if err != nil {
		return "", errors.New("error getting balance")
	}

	amountInt64 := int64(sendAmount * 1e6)

	//check if balance is enough
	availableBalance := bal.Received - bal.Spent

	if availableBalance < amountInt64 {
		log.Print("Insufficient balance")
		return "", errors.New("insufficient balance")
	}

	tx := MetaTx{
		To:    toAddress,
		Value: amountInt64,
		Nonce: bal.CountSpent + 1,
		Fee:   0,
		Data:  "",
	}

	return mk.sendTransaction(&tx)
}

func (mk *MetaKey) sendTransaction(tr *MetaTx) (string, error) {

	data, err := getSignData(tr)
	if err != nil {
		return "", err
	}

	sign, err := mk.SignTransaction(data)
	if err != nil {
		return "", err
	}

	arg := MetaTxArg{
		To:     tr.To,
		Value:  strconv.FormatInt(tr.Value, 10),
		Fee:    strconv.FormatInt(tr.Fee, 10),
		Data:   tr.Data,
		Nonce:  strconv.FormatInt(tr.Nonce, 10),
		Pubkey: string(mk.Key.Public()),
		Sign:   sign,
	}

	sendClient := NewClient(proxUrl)

	resp, err := sendClient.Call("mhc_send", arg)

	var txHash string
	if err == nil {
		err = resp.GetObject(&txHash)
		if err == nil {
			return txHash, nil
		}
		return "", err
	}

	return txHash, nil
}

//SignTransaction generates signature
func (mk *MetaKey) SignTransaction(data []byte) (string, error) {

	sign, err := mk.Key.Sign(data)
	if err != nil {
		return "", err
	}
	return sign, nil
}

func getSignData(tr *MetaTx) ([]byte, error) {
	vrt := NewVarint()
	to, err := hex.DecodeString(strings.TrimPrefix(string(tr.To), "0x"))
	if err != nil {
		log.Panic("Could not convert address into hex")
	}
	vrt.AppendBytes(to)

	if err := vrt.Append(big.NewInt(tr.Value)); err != nil {
		return nil, err
	}

	if err := vrt.Append(big.NewInt(tr.Fee)); err != nil {
		return nil, err
	}

	if err := vrt.Append(big.NewInt(tr.Nonce)); err != nil {
		return nil, err
	}

	//get data length
	l := int64(len(tr.Data) / 2)

	if l > 0 {
		if err := vrt.Append(big.NewInt(l)); err != nil {
			return nil, err
		}
		vrt.AppendString(tr.Data)
	} else {
		vrt.AppendBytes(nil)
	}
	return vrt.GetBytes(), nil
}
