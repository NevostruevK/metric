// Package sign предназначен для подписи строки по алгоритму sha256.
package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Hash(s, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	if _, err := h.Write([]byte(s)); err != nil {
		return "", fmt.Errorf("hash sign with the error %w", err)
	}
	return (hex.EncodeToString(h.Sum(nil))), nil
}
