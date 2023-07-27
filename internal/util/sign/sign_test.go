package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	t.Run("table test", func(t *testing.T) {
		type args struct {
			s   string
			key string
		}
		tests := []struct {
			name    string
			args    args
			want    string
			wantErr bool
		}{
			{
				name:    "normal test",
				args:    args{s: "some metric", key: "some key"},
				want:    "2ac2b6ca2632df66613b05bdc6cfddf6d4e7a2ff389a0078b46a58d5a4b1c8d3",
				wantErr: false,
			},
			{
				name:    "empty strings",
				args:    args{s: "", key: ""},
				want:    "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Hash(tt.args.s, tt.args.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Hash() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("is idempotent", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			got, err := Hash("some metric", "some key")
			require.NoError(t, err)
			assert.Equal(t, "2ac2b6ca2632df66613b05bdc6cfddf6d4e7a2ff389a0078b46a58d5a4b1c8d3", got)
		}
	})
	t.Run("is different", func(t *testing.T) {
		gotBasic, err := Hash("some metric", "some key")
		require.NoError(t, err)
		gotOtherMetric, err := Hash("other metric", "some key")
		require.NoError(t, err)
		gotOtherKey, err := Hash("some metric", "other key")
		require.NoError(t, err)
		assert.NotEqual(t, gotBasic, gotOtherMetric)
		assert.NotEqual(t, gotBasic, gotOtherKey)
		assert.NotEqual(t, gotOtherMetric, gotOtherKey)
	})

}
