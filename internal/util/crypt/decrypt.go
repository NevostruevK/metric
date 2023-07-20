package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"errors"

	"github.com/NevostruevK/metric/internal/util/crypt/session"
)

var ErrInputIsTooSmall = errors.New("input data is smaller than session key size")

type Decrypt struct {
	PrivateKey *rsa.PrivateKey
	Nonce      [12]byte
}

func NewDecrypt(fname string) (*Decrypt, error) {
	PrivateKey, err := GetPrivateKey(fname)
	if err != nil {
		return nil, err
	}
	return &Decrypt{PrivateKey: PrivateKey, Nonce: *new([12]byte)}, nil
}

func (d Decrypt) Decrypt(encrypted []byte) ([]byte, error) {
	k, err := session.NewKey()
	if err != nil {
		return encrypted, err
	}
	if len(encrypted) < session.EncryptedSessionKeySize {
		return encrypted, ErrInputIsTooSmall
	}
	err = k.Decrypt(d.PrivateKey, encrypted[:session.EncryptedSessionKeySize])
	if err != nil {
		return encrypted, err
	}
	aesblock, err := aes.NewCipher(k)
	if err != nil {
		return encrypted, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return encrypted, err
	}
	return aesgcm.Open(nil, d.Nonce[:], encrypted[session.EncryptedSessionKeySize:], nil)
}
