package crypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

const (
	PrivateKeyTitle = "RSA PRIVATE KEY"
	PublicKeyTitle  = "RSA PUBLIC KEY"
)

func readFile(fname string) (b []byte, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
	}()
	b, err = io.ReadAll(f)
	return
}

func GetPrivateKey(fname string) (*rsa.PrivateKey, error) {
	op := "GetPrivateKey"
	b, err := readFile(fname)
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	block, _ := pem.Decode(b)
	if block == nil || block.Type != PrivateKeyTitle {
		return nil, fmt.Errorf("%s failed to decode PEM block containing private key", op)
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func GetPublicKey(fname string) (*rsa.PublicKey, error) {
	op := "GetPublicKey"
	b, err := readFile(fname)
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	block, _ := pem.Decode(b)
	if block == nil || block.Type != PublicKeyTitle {
		return nil, fmt.Errorf("%s failed to decode PEM block containing public key", op)
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}
