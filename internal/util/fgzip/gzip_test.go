// Package fgzip модуль компрессии/декомпрессии данных по алгоритму gzip
package fgzip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompress(t *testing.T) {
	t.Run("is idempotent", func(t *testing.T) {
		input := []byte("48656c6c6f20476f7068657221")
		output, err := Compress(input)
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			outLoop, err := Compress(input)
			require.NoError(t, err)
			assert.Equal(t, output, outLoop)
		}
	})

	t.Run("is different", func(t *testing.T) {
		input1 := []byte("48656c6c6f20476f7068657221")
		input2 := []byte("21726568706f47206fc6c65686")
		output1, err := Compress(input1)
		require.NoError(t, err)
		output2, err := Compress(input2)
		require.NoError(t, err)
		assert.NotEqual(t, output1, output2)
		assert.NotEqual(t, input1, output1)
		assert.NotEqual(t, input2, output2)
	})
}

func TestDecompress(t *testing.T) {
	input1 := []byte("48656c6c6f20476f7068657221")
	output1, err := Compress(input1)
	require.NoError(t, err)
	input2 := []byte("")
	output2, err := Compress(input2)
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "normal input",
			input:   output1,
			want:    input1,
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   output2,
			want:    input2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decompress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}

}
