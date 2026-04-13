package main

import (
	"2/internal/app"
	"2/internal/config"
	"2/internal/repository"
	"2/internal/services/grpcserv"
	"2/internal/services/restserv"
	"context"
	"fmt"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"google.golang.org/grpc"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println("fail to init config:", err)
		return
	}

	logger := setupLogger(cfg.LogLevel)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	dataBase, err := repository.InitDataBase(ctx, connString)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		return
	}

	grpcAuth, err := grpc.Dial("auth-service:50051", grpc.WithInsecure())
	if err != nil {
		logger.Error("failed to connect to auth service", "error", err)
	}

	defer func() {
		if err := grpcAuth.Close(); err != nil {
			logger.Info("failed to close grpc auth service: ", err)
		}
	}()

	grpcAuthClient := auth.NewAuthClient(grpcAuth)

	grpcSevice := grpcserv.NewGrpcBookService(dataBase)
	restService := restserv.NewRestBookService(dataBase)

	application := app.NewMainApp(logger, cfg, grpcSevice, restService, grpcAuthClient)

	go application.GRPCServer.MustRun()
	go application.RESTServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	logger.Info("received signal", sign.String())

	application.GRPCServer.Stop()

	ctxTimeout, cancel := context.WithTimeout(ctx, cfg.ServTimeout)
	defer cancel()

	if err := application.RESTServer.Stop(ctxTimeout); err != nil {
		logger.Error("failed to stop REST server", "error", err)
		return
	}

	logger.Info("rest server stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	default:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return logger
}
