package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"errors"

	"github.com/NevostruevK/metric/internal/util/crypt/session"
)

var ErrCryptNotInit = errors.New("entity Crypt isn't initialized")

type Crypt struct {
	PublicKey *rsa.PublicKey
	Nonce     [12]byte
}

func NewCrypt(fname string) (*Crypt, error) {
	if fname == "" {
		return nil, nil
	}
	PublicKey, err := GetPublicKey(fname)
	if err != nil {
		return nil, err
	}
	return &Crypt{PublicKey: PublicKey, Nonce: *new([12]byte)}, nil
}

func (c *Crypt) Crypt(raw []byte) ([]byte, error) {
	if c == nil {
		return raw, ErrCryptNotInit
	}
	k, err := session.NewKey()
	if err != nil {
		return raw, err
	}
	aesblock, err := aes.NewCipher(k)
	if err != nil {
		return raw, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return raw, err
	}
	eKey, err := k.Encrypt(c.PublicKey)
	if err != nil {
		return raw, err
	}
	encrypted := make([]byte, 0, len(raw)+session.EncryptedSessionKeySize+len(k))
	encrypted = append(encrypted, eKey...)
	encrypted = append(encrypted, aesgcm.Seal(nil, c.Nonce[:], raw, nil)...)
	return encrypted, nil
}
