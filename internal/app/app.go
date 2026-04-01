package app

import (
	grpcapp "2/internal/app/grpc"
	restapp "2/internal/app/rest"
	"2/internal/config"
	"2/internal/services/grpcserv"
	"2/internal/services/restserv"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"log/slog"
)

type MainApp struct {
	GRPCServer *grpcapp.GrpcBookApp
	RESTServer *restapp.RestBookApp
}

func NewMainApp(log *slog.Logger,
	cfg *config.Config,
	grpcServ *grpcserv.GrpcBookService,
	restServ *restserv.RestBookService,
	grpcAuthClient auth.AuthClient,
) *MainApp {
	gRPCSServer := grpcapp.NewGrpcBookApp(log, cfg.GrpcPort, grpcServ)
	restServer := restapp.NewRestBookApp(log, cfg, restServ, grpcAuthClient)

	return &MainApp{
		GRPCServer: gRPCSServer,
		RESTServer: restServer,
	}
}
