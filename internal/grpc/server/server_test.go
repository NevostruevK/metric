package server_test

import (
	"log"

	grpcserver "github.com/NevostruevK/metric/internal/grpc/server"
	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/crypt"
	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func startClient(cfg *commands.Config) (*grpc.ClientConn, pb.MetricsClient, error) {

	var (
		creds = insecure.NewCredentials()
		err   error
	)
	if cfg.Certificate != "" {
		creds, err = credentials.NewClientTLSFromFile(cfg.Certificate, crypt.HostCertificateAddress)
		if err != nil {
			log.Fatalf("failed to load credentials: %v", err)
		}
	}
	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, err
	}
	return conn, pb.NewMetricsClient(conn), nil
}

func startServer(cfg *commands.Config) (*grpcserver.MetricsServer, error) {
	st := storage.NewMemStorage(false, false, "")
	return grpcserver.NewServer(st, cfg)
}
