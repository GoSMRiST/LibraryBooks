package grpc

import (
	"2/internal/core"
	grpcproto "2/internal/transport/grpc"
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GrpcBookApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

type BookService interface {
	CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

func NewGrpcBookApp(log *slog.Logger, port string, bookService BookService) *GrpcBookApp {
	gRPCServer := grpc.NewServer()

	grpcproto.NewServer(gRPCServer, bookService)
	return &GrpcBookApp{log: log, gRPCServer: gRPCServer, port: port}
}

func (app *GrpcBookApp) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *GrpcBookApp) Run() error {
	log := app.log

	l, err := net.Listen("tcp", ":"+app.port)
	if err != nil {
		return err
	}

	log.Info("Starting gRPC server on port %d", app.port, slog.String("addr", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (app *GrpcBookApp) Stop() {

	app.log.Info("Stopping gRPC server on port %d", app.port)

	app.gRPCServer.GracefulStop()

	app.log.Info("Stopped gRPC server on port %d", app.port)
}
