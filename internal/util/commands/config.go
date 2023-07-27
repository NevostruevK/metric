package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/NevostruevK/metric/internal/util/commands/duration"
)

type Config struct {
	ReportInterval duration.Duration `json:"report_interval"`
	PollInterval   duration.Duration `json:"poll_interval"`
	StoreInterval  duration.Duration `json:"store_interval"`
	Address        string            `json:"address"`
	StoreFile      string            `json:"store_file"`
	HashKey        string            `json:"hash_key"`
	DataBaseDSN    string            `json:"database_dsn"`
	CryptoKey      string            `json:"crypto_key"`
	TrustedSubnet  string            `json:"trusted_subnet"`
	RateLimit      int               `json:"rate_limit"`
	Restore        bool              `json:"restore"`
}

func NewServerConfig() *Config {
	return &Config{
		Address:       defAddress,
		StoreFile:     defStoreFile,
		StoreInterval: duration.NewDuration(defStoreInterval),
		Restore:       defRestore,
	}
}

func NewAgentConfig() *Config {
	return &Config{
		ReportInterval: duration.NewDuration(defReportInterval),
		PollInterval:   duration.NewDuration(defPollInterval),
		Address:        defAddress,
		RateLimit:      defRateLimit,
	}
}

func (o *Config) ReadConfig(fname string) error {
	if fname == "" {
		return nil
	}
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(o); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (o *Config) setOption(set func(*Config)) {
	set(o)
}

func withReportInterval(reportInterval duration.Duration) func(*Config) {
	return func(o *Config) {
		o.ReportInterval = reportInterval
	}
}

func withPollInterval(pollInterval duration.Duration) func(*Config) {
	return func(o *Config) {
		o.PollInterval = pollInterval
	}
}

func withStoreInterval(stroreInterval duration.Duration) func(*Config) {
	return func(o *Config) {
		o.StoreInterval = stroreInterval
	}
}

func withAddress(address string) func(*Config) {
	return func(o *Config) {
		o.Address = address
	}
}

func withStoreFile(storeFile string) func(*Config) {
	return func(o *Config) {
		o.StoreFile = storeFile
	}
}

func withHashKey(hashKey string) func(*Config) {
	return func(o *Config) {
		o.HashKey = hashKey
	}
}

func withDataBaseDSN(dataBaseDSN string) func(*Config) {
	return func(o *Config) {
		o.DataBaseDSN = dataBaseDSN
	}
}

func withCryptoKey(cryptoKey string) func(*Config) {
	return func(o *Config) {
		o.CryptoKey = cryptoKey
	}
}

func withTrustedSubnet(trustedSubnet string) func(*Config) {
	return func(o *Config) {
		o.TrustedSubnet = trustedSubnet
	}
}

func withRateLimit(rateLimit int) func(*Config) {
	return func(o *Config) {
		o.RateLimit = rateLimit
	}
}

func withRestore(restore bool) func(*Config) {
	return func(o *Config) {
		o.Restore = restore
	}
}
