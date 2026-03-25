package grpc

import (
	"2/internal/core"
	grpcserv "2/internal/transport/grpc/book"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

type BookService interface {
	CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

func New(log *slog.Logger, port int, bookService BookService) *App {
	gRPCServer := grpc.NewServer()

	grpcserv.NewServer(gRPCServer, bookService)
	return &App{log: log, gRPCServer: gRPCServer, port: port}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	log := app.log

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return err
	}

	log.Info("Starting gRPC server on port %d", app.port, slog.String("addr", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop() {
	log := app.log

	app.gRPCServer.GracefulStop()

	log.Info("Stopping gRPC server on port %d", app.port)
}
