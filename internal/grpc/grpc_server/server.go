package grpcserver

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/NevostruevK/metric/proto/mproto"
)

type metricsServer struct {
	pb.UnimplementedUsersServer
}
