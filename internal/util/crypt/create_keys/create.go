package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/NevostruevK/metric/internal/util/crypt"
)

const (
	privateKeyFileName = "private.key"
	publicKeyFileName  = "public.key"
)

func main() {
	privateKey, err := createPrivateKey(privateKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = createPublicKey(publicKeyFileName, privateKey)
	if err != nil {
		log.Fatal(err)
	}
}

func createPrivateKey(fname string) (*rsa.PrivateKey, error) {
	op := "createPrivateKey"
	f, err := os.Create(fname)
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	defer func() {
		err = f.Close()
	}()
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	err = pem.Encode(f, &pem.Block{
		Type:  crypt.PrivateKeyTitle,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	return privateKey, nil
}

func createPublicKey(fname string, privateKey *rsa.PrivateKey) (*rsa.PublicKey, error) {
	op := "createPublicKey"
	f, err := os.Create(fname)
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	defer func() {
		err = f.Close()
	}()
	err = pem.Encode(f, &pem.Block{
		Type:  crypt.PublicKeyTitle,
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		return nil, fmt.Errorf("%s failed with an error %w", op, err)
	}
	return &privateKey.PublicKey, nil
}
