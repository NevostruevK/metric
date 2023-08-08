package crypt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	t.Run("is idempotent", func(t *testing.T) {
		b1, err := readFile("private.pem")
		require.NoError(t, err)
		b2, err := readFile("private.pem")
		require.NoError(t, err)
		assert.True(t, bytes.Equal(b1, b2))
	})
	t.Run("is different", func(t *testing.T) {
		b1, err := readFile("private.pem")
		require.NoError(t, err)
		b2, err := readFile("public.pem")
		require.NoError(t, err)
		assert.True(t, !bytes.Equal(b1, b2))
	})
	t.Run("error no file", func(t *testing.T) {
		b, err := readFile("nonExistentFile")
		require.Error(t, err)
		assert.True(t, len(b) == 0)
	})
}

func TestGetPrivateKey(t *testing.T) {
	t.Run("is idempotent", func(t *testing.T) {
		k1, err := GetPrivateKey("private.pem")
		require.NoError(t, err)
		k2, err := GetPrivateKey("private.pem")
		require.NoError(t, err)
		assert.Equal(t, k1, k2)
	})
	t.Run("is not nil", func(t *testing.T) {
		k, err := GetPrivateKey("private.pem")
		require.NoError(t, err)
		require.NotNil(t, k)
		require.NotEmpty(t, k)
	})
	t.Run("error not a private.key", func(t *testing.T) {
		_, err := GetPrivateKey("public.pem")
		require.Error(t, err)
	})
}

func TestGetPublicKey(t *testing.T) {
	t.Run("is idempotent", func(t *testing.T) {
		k1, err := GetPublicKey("public.pem")
		require.NoError(t, err)
		k2, err := GetPublicKey("public.pem")
		require.NoError(t, err)
		assert.Equal(t, k1, k2)
	})
	t.Run("is not nil", func(t *testing.T) {
		k, err := GetPublicKey("public.pem")
		require.NoError(t, err)
		require.NotNil(t, k)
		require.NotEmpty(t, k)
	})
	t.Run("error not a public.key", func(t *testing.T) {
		_, err := GetPublicKey("private.pem")
		require.Error(t, err)
	})
}
