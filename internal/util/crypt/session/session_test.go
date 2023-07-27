package session

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandom(t *testing.T) {
	t.Run("zero size", func(t *testing.T) {
		b, err := GenerateRandom(0)
		require.NoError(t, err)
		assert.True(t, len(b) == 0)
	})
	t.Run("is different", func(t *testing.T) {
		const size = 32
		b1, err := GenerateRandom(size)
		require.NoError(t, err)
		b2, err := GenerateRandom(size)
		require.NoError(t, err)
		assert.NotEqual(t, b1, b2)
		assert.True(t, !bytes.Equal(b1, b2))
	})
}

func TestNewKey(t *testing.T) {
	t.Run("is different", func(t *testing.T) {
		k1, err := NewKey()
		require.NoError(t, err)
		k2, err := NewKey()
		require.NoError(t, err)
		assert.NotEqual(t, k1, k2)
		assert.True(t, !bytes.Equal(k1, k2))
	})
}

func TestEncrypt(t *testing.T) {
	t.Run("is different", func(t *testing.T) {
		sk, err := NewKey()
		require.NoError(t, err)
		pk, err := rsa.GenerateKey(rand.Reader, 4096)
		require.NoError(t, err)
		encr1, err := sk.Encrypt(&pk.PublicKey)
		require.NoError(t, err)
		encr2, err := sk.Encrypt(&pk.PublicKey)
		require.NoError(t, err)
		assert.NotEqual(t, encr1, encr2)
	})
	t.Run("is constant size", func(t *testing.T) {

		for i := 0; i < 2; i++ {
			sk, err := NewKey()
			require.NoError(t, err)
			pk, err := rsa.GenerateKey(rand.Reader, 4096)
			require.NoError(t, err)
			encr, err := sk.Encrypt(&pk.PublicKey)
			require.NoError(t, err)
			assert.True(t, len(encr) == EncryptedSessionKeySize)
		}
	})
}

func TestDecrypt(t *testing.T) {
	skAgent, err := NewKey()
	require.NoError(t, err)
	pk1, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)
	encr, err := skAgent.Encrypt(&pk1.PublicKey)
	require.NoError(t, err)
	t.Run("test ok", func(t *testing.T) {
		skServer, err := NewKey()
		require.NoError(t, err)
		err = skServer.Decrypt(pk1, encr)
		require.NoError(t, err)
		assert.Equal(t, skAgent, skServer)
	})
}
