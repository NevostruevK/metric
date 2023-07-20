package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypt_Crypt(t *testing.T) {

	t.Run(" normal", func(t *testing.T) {
		raw := []byte("Some message that you need to encrypt")
		cr, err := NewCrypt("public.pem")
		require.NoError(t, err)
		encrypted, err := cr.Crypt(raw)
		require.NoError(t, err)
		require.NotEqual(t, raw, encrypted)
		dcr, err := NewDecrypt("private.pem")
		require.NoError(t, err)
		message, err := dcr.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, raw, message)
	})
}
