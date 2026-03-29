package grpcserver

import (
	"context"
	"net"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/galaxy-empire-team/bridge-api/api/gen/go/planet/v1"
	"github.com/galaxy-empire-team/bridge-api/internal/config"
)

const (
	tcpProtocol = "tcp"
)

type PlanetService interface {
	UpdatePlanetResources(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, updatedTime *time.Time) error
}

type Server struct {
	planetpb.UnimplementedPlanetServiceServer
	planetService PlanetService

	logger *zap.Logger
}

func New(planetService PlanetService, logger *zap.Logger) *Server {
	return &Server{
		planetService: planetService,
		logger:        logger,
	}
}

func (s *Server) Start(cfg config.GRPCServer) (func(context.Context) error, error) {
	grpcServer := grpc.NewServer()

	planetpb.RegisterPlanetServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	lis, err := net.Listen(tcpProtocol, cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	go func() {
		s.logger.Info("---  gRPC server started  ---", zap.String("endpoint", cfg.Endpoint))
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	// returns error to pass it to app shutdown pipeline
	return func(_ context.Context) error {
		grpcServer.GracefulStop()
		return nil
	}, nil
}
