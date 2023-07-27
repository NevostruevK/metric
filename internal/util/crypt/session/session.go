package session

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
)

const EncryptedSessionKeySize = 512

type Key []byte

func NewKey() (Key, error) {
	return GenerateRandom(2 * aes.BlockSize)
}

func (k Key) Encrypt(public *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, public, k)
}

func (k Key) Decrypt(private *rsa.PrivateKey, encrypted []byte) error {
	return rsa.DecryptPKCS1v15SessionKey(nil, private, encrypted, k)
}

func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
