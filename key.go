package metahash

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

//Address data type
type Address string

//PrivateKey data type
type PrivateKey string

//PublicKey data type
type PublicKey string

//Key interface
type Key interface {
	Private() PrivateKey
	Public() PublicKey
	Verify(data []byte, sign string) (bool, error)
	Address() Address
	SetAddress(address string)
	Sign(data []byte) (string, error)
}

type key struct {
	priv    *ecdsa.PrivateKey
	pub     *ecdsa.PublicKey
	address string
}

type ecdsaSignature struct {
	R, S *big.Int
}

//NewKey generates metahash key
func NewKey() (Key, error) {
	return newKey()
}

//CreateKey creates metahash keys
func CreateKey(private string) (Key, error) {
	priv := PrivateKey(private)
	return createKey(priv)
}

// https://support.metahash.org/hc/ru/articles/360002712193
func newKey() (Key, error) {
	curve := elliptic.P256() // secp256r1 by default
	rnd := rand.Reader
	priv, err := ecdsa.GenerateKey(curve, rnd)
	if err != nil {
		return nil, err
	}

	pub := priv.Public()
	p, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cant cast [%s] to [*ecdsa.PublicKey]", reflect.TypeOf(p))
	}
	return &key{
		priv: priv,
		pub:  p,
	}, nil

}

func createKey(private PrivateKey) (Key, error) {

	if private == "" {
		return nil, errors.New("private key cannot be empty")
	}

	// r1 key type has prefix 3077
	if strings.HasPrefix(string(private), "3077") {
		decoded, err := hex.DecodeString(string(private))
		if err != nil {
			return nil, err
		}
		priv, err := x509.ParseECPrivateKey(decoded)
		if err != nil {
			return nil, err
		}

		return &key{
			priv: priv,
			pub:  &priv.PublicKey,
			//address: currentAddress, //TODO:generate this address from the public key
		}, nil
	} else if strings.HasPrefix(string(private), "3074") { // k1 key type has prefix 3074
		// k1 type is not yet implemented
		return nil, errors.New("secp256k1 key is currently not supported")
		//return secp256k1.FromHexString(privateKey)
	} else {
		return nil, errors.New("unknown private key")
	}
}

//Private returns private metahash private key
func (t *key) Private() PrivateKey {
	x509EncodedPriv, _ := x509.MarshalECPrivateKey(t.priv)
	return PrivateKey(hex.EncodeToString(x509EncodedPriv))
}

func (t *key) Public() PublicKey {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(t.pub)
	return PublicKey(hex.EncodeToString(x509EncodedPub))
}

//Sign data
func (t *key) Sign(data []byte) (string, error) {
	digest := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, t.priv, digest[:])
	if err != nil {
		return "", err
	}

	b, e := asn1.Marshal(ecdsaSignature{r, s})

	return hex.EncodeToString(b), e
}

//Verify signature
func (t *key) Verify(data []byte, sign string) (bool, error) {
	digest := sha256.Sum256(data)

	decoded, err := hex.DecodeString(string(sign))
	if err != nil {
		return false, err
	}

	var signEcdsa ecdsaSignature

	_, err = asn1.Unmarshal(decoded, &signEcdsa)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(t.pub, digest[:], signEcdsa.R, signEcdsa.S), nil
}

//Address returns address associated with the key
func (t *key) Address() Address {
	return Address(t.address)
}

//SetAddress
//TODO: Replace this function by generating address from the public key
func (t *key) SetAddress(address string) {
	t.address = address
}
