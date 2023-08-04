package commands

import (
	"testing"
	"time"

	"github.com/NevostruevK/metric/internal/util/commands/duration"
	"github.com/stretchr/testify/require"
)

func TestConfig_setOption(t *testing.T) {
	t.Run("ok set all options", func(t *testing.T) {
		exp := Config{
			ReportInterval: duration.NewDuration(time.Second),
			PollInterval:   duration.NewDuration(time.Minute),
			StoreInterval:  duration.NewDuration(time.Hour),
			Address:        "test_address",
			StoreFile:      "test_store_file",
			HashKey:        "test_hash_key",
			DataBaseDSN:    "test_dsn",
			CryptoKey:      "test_crypto_key",
			TrustedSubnet:  "192.168.0.15/24",
			RateLimit:      1234,
			Restore:        false,
		}
		c := Config{}
		c.SetOption(WithReportInterval(duration.NewDuration(time.Second)))
		c.SetOption(WithPollInterval(duration.NewDuration(time.Minute)))
		c.SetOption(WithStoreInterval(duration.NewDuration(time.Hour)))
		c.SetOption(WithAddress("test_address"))
		c.SetOption(WithStoreFile("test_store_file"))
		c.SetOption(WithHashKey("test_hash_key"))
		c.SetOption(WithDataBaseDSN("test_dsn"))
		c.SetOption(WithCryptoKey("test_crypto_key"))
		c.SetOption(WithTrustedSubnet("192.168.0.15/24"))
		c.SetOption(WithRateLimit(1234))
		c.SetOption(WithRestore(false))
		require.Equal(t, exp, c)
	})
}
