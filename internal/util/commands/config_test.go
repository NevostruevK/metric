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
		c.setOption(withReportInterval(duration.NewDuration(time.Second)))
		c.setOption(withPollInterval(duration.NewDuration(time.Minute)))
		c.setOption(withStoreInterval(duration.NewDuration(time.Hour)))
		c.setOption(withAddress("test_address"))
		c.setOption(withStoreFile("test_store_file"))
		c.setOption(withHashKey("test_hash_key"))
		c.setOption(withDataBaseDSN("test_dsn"))
		c.setOption(withCryptoKey("test_crypto_key"))
		c.setOption(withTrustedSubnet("192.168.0.15/24"))
		c.setOption(withRateLimit(1234))
		c.setOption(withRestore(false))
		require.Equal(t, exp, c)
	})
}
