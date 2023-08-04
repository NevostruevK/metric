package crypt

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	PrivateKeyTitle  = "RSA PRIVATE KEY"
	PublicKeyTitle   = "RSA PUBLIC KEY"
	CertificateTitle = "CERTIFICATE"
)
const HostCertificateAddress = "127.0.0.1"

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

func NewCertificateTemplate() *x509.Certificate {
	return &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		//        IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		IPAddresses: []net.IP{net.ParseIP(HostCertificateAddress), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
}
