package xgrpc

import (
	"github.com/AltScore/gothic/pkg/xlogger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type Config struct {
	Address string
}

type Server struct {
	*grpc.Server
	logger xlogger.Logger
	config Config
}

func NewServer(logger xlogger.Logger, config Config, unaryInterceptors ...grpc.UnaryServerInterceptor) *Server {
	interceptors := []grpc.UnaryServerInterceptor{
		NewLoggerInterceptor(logger),
	}
	interceptors = append(interceptors, unaryInterceptors...)
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors...),
	}

	return &Server{
		Server: grpc.NewServer(opts...),
		logger: logger,
		config: config,
	}
}

// Start starts the grpc server
// This is a non-blocking call, it will start the server in a goroutine
func (s *Server) Start() error {
	address := s.config.Address

	lis, err := net.Listen("tcp", address)

	if err != nil {
		s.logger.Error("failed to listen for grpc server", zap.String("address", address), zap.Error(err))
		return err
	}

	go func() {
		s.logger.Info("Starting grpc server", zap.String("address", address))
		if err := s.Server.Serve(lis); err != nil {
			s.logger.Error("failed to start grpc server", zap.String("address", address), zap.Error(err))
			panic(err)
		}
	}()
	return nil
}
