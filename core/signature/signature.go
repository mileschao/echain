package signature

import (
	"errors"

	"github.com/ontio/ontology-crypto/keypair"
	ontsig "github.com/ontio/ontology-crypto/signature"
)

var (
	//ErrVerifySign verify signature failed
	ErrVerifySign = errors.New("signature verification failed")
	//ErrInvalidSignData invalide signature data
	ErrInvalidSignData = errors.New("invalid signature data")
	//ErrNotEnoughtSignature not enought signature
	ErrNotEnoughtSignature = errors.New("not enought signature")
)

//Signatory the one who sign
type Signatory interface {
	//PrivateKey private key
	PrivateKey() keypair.PrivateKey
	//PublicKey public key
	PublicKey() keypair.PublicKey
	//Scheme signature scheme
	Scheme() ontsig.SignatureScheme
}

// Sign get the signature of data by private key
func Sign(signatory Signatory, data []byte) ([]byte, error) {
	signature, err := ontsig.Sign(signatory.Scheme(), signatory.PrivateKey(), data, nil)
	if err != nil {
		return nil, err
	}
	return ontsig.Serialize(signature)
}

// Verify check the signature of data by public
func Verify(pubKey keypair.PublicKey, data []byte, signature []byte) error {
	sigObj, err := ontsig.Deserialize(signature)
	if err != nil {
		return ErrInvalidSignData
	}
	if !ontsig.Verify(pubKey, data, sigObj) {
		return ErrVerifySign
	}
	return nil
}

// VerifyMultiSignature check whether more than m signature have been signed by the keys in key array
func VerifyMultiSignature(data []byte, keys []keypair.PublicKey, m int, sigs [][]byte) error {
	//TODO: ugly code, make it more readable
	n := len(keys)

	if len(sigs) < m {
		return ErrNotEnoughtSignature
	}

	mask := make([]bool, n)
	for i := 0; i < m; i++ {
		valid := false

		sig, err := ontsig.Deserialize(sigs[i])
		if err != nil {
			return ErrInvalidSignData
		}
		for j := 0; j < n; j++ {
			if mask[j] {
				continue
			}
			if ontsig.Verify(keys[j], data, sig) {
				mask[j] = true
				valid = true
				break
			}
		}

		if valid == false {
			return ErrNotEnoughtSignature
		}
	}
	return nil
}
