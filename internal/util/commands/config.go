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
	Certificate    string            `json:"certificate"`
	TrustedSubnet  string            `json:"trusted_subnet"`
	RateLimit      int               `json:"rate_limit"`
	Restore        bool              `json:"restore"`
	GRPC           bool              `json:"grpc"`
}

func NewServerConfig() *Config {
	return &Config{
		Address:       defAddress,
		StoreFile:     defStoreFile,
		StoreInterval: duration.NewDuration(defStoreInterval),
		Restore:       defRestore,
		GRPC:          defGRPC,
	}
}

func NewAgentConfig() *Config {
	return &Config{
		ReportInterval: duration.NewDuration(defReportInterval),
		PollInterval:   duration.NewDuration(defPollInterval),
		Address:        defAddress,
		RateLimit:      defRateLimit,
		GRPC:           defGRPC,
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

func (o *Config) SetOption(set func(*Config)) {
	set(o)
}

func WithReportInterval(reportInterval duration.Duration) func(*Config) {
	return func(o *Config) {
		o.ReportInterval = reportInterval
	}
}

func WithPollInterval(pollInterval duration.Duration) func(*Config) {
	return func(o *Config) {
		o.PollInterval = pollInterval
	}
}

func WithStoreInterval(stroreInterval duration.Duration) func(*Config) {
	return func(o *Config) {
		o.StoreInterval = stroreInterval
	}
}

func WithAddress(address string) func(*Config) {
	return func(o *Config) {
		o.Address = address
	}
}

func WithStoreFile(storeFile string) func(*Config) {
	return func(o *Config) {
		o.StoreFile = storeFile
	}
}

func WithHashKey(hashKey string) func(*Config) {
	return func(o *Config) {
		o.HashKey = hashKey
	}
}

func WithDataBaseDSN(dataBaseDSN string) func(*Config) {
	return func(o *Config) {
		o.DataBaseDSN = dataBaseDSN
	}
}

func WithCryptoKey(cryptoKey string) func(*Config) {
	return func(o *Config) {
		o.CryptoKey = cryptoKey
	}
}

func WithCertificate(certificate string) func(*Config) {
	return func(o *Config) {
		o.Certificate = certificate
	}
}

func WithTrustedSubnet(trustedSubnet string) func(*Config) {
	return func(o *Config) {
		o.TrustedSubnet = trustedSubnet
	}
}

func WithRateLimit(rateLimit int) func(*Config) {
	return func(o *Config) {
		o.RateLimit = rateLimit
	}
}

func WithRestore(restore bool) func(*Config) {
	return func(o *Config) {
		o.Restore = restore
	}
}

func WithGRPC(grpc bool) func(*Config) {
	return func(o *Config) {
		o.GRPC = grpc
	}
}
