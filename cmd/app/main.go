package main

import (
	grpcapp "2/internal/app/grpc"
	"2/internal/config"
	"2/internal/repository"
	"2/internal/services/grpcserv"
	"2/internal/services/restserv"
	"2/internal/transport/rest"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	appConfig, err := config.InitConfig()
	if err != nil {
		fmt.Println("fail to init config:", err)
		return
	}

	logger := setupLogger(appConfig.LogLevel)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBName,
	)

	dataBase, err := repository.InitDataBase(ctx, connString)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		return
	}

	serviceDB := restserv.NewBookService(dataBase)
	bookHandler := rest.NewBookHandler(serviceDB)

	engine := gin.Default()
	bookHandler.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:         appConfig.HostAddress,
		Handler:      engine,
		ReadTimeout:  appConfig.ServTimeout,
		WriteTimeout: appConfig.ServTimeout,
		IdleTimeout:  appConfig.ServTimeout,
	}

	// REST server
	go func() {
		logger.Info("REST server started",
			slog.String("addr", appConfig.HostAddress),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("REST server failed", slog.Any("error", err))
			return
		}
	}()

	// gRPC server
	grpcService := grpcserv.NewGrpcBookService(dataBase)
	port, err := strconv.Atoi(appConfig.GrpcPort)
	if err != nil {
		logger.Error("invalid grpc port", slog.Any("error", err))
		return
	}

	grpcApp := grpcapp.New(logger, port, grpcService)
	if grpcApp == nil {
		logger.Error("failed to create gRPC server")
		return
	}

	go func() {
		logger.Info("gRPC server started", slog.Int("port", port))
		grpcApp.MustRun()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Info("shutdown signal received")

	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// REST
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error("REST shutdown error", slog.Any("error", err))
	}

	// gRPC
	grpcApp.Stop()
	logger.Info("gRPC server stopped")

	// DB
	if err := dataBase.CloseDataBase(ctx); err != nil {
		logger.Error("failed to close database", slog.Any("error", err))
		return
	}
	logger.Info("DataBase is closed")
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
