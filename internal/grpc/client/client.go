package client

import (
	"github.com/NevostruevK/metric/internal/util/commands"
	"github.com/NevostruevK/metric/internal/util/crypt"
	pb "github.com/NevostruevK/metric/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	C    pb.MetricsClient
}

func NewClient(cfg *commands.Config) (*Client, error) {
	var (
		creds = insecure.NewCredentials()
		err   error
	)
	if cfg.Certificate != "" {
		creds, err = credentials.NewClientTLSFromFile(cfg.Certificate, crypt.HostCertificateAddress)
		if err != nil {
			return nil, err
		}
	}
	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, C: pb.NewMetricsClient(conn)}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

