package rest

import (
	"2/internal/config"
	"2/internal/middleware"
	"2/internal/services/restserv"
	"2/internal/transport/rest"
	"context"
	"errors"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type RestBookApp struct {
	log        *slog.Logger
	restServer *http.Server
	hostAddr   string
}

func NewRestBookApp(log *slog.Logger, cfg *config.Config, restService *restserv.RestBookService, grpcAuthClient auth.AuthClient) *RestBookApp {
	bookHandler := rest.NewBookHandler(log, restService)

	engine := gin.Default()
	engine.Use(middleware.AuthMiddleware(grpcAuthClient))

	bookHandler.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:         cfg.HostAddress,
		Handler:      engine,
		ReadTimeout:  cfg.ServTimeout,
		WriteTimeout: cfg.ServTimeout,
		IdleTimeout:  cfg.ServTimeout,
	}

	return &RestBookApp{
		log:        log,
		restServer: srv,
		hostAddr:   cfg.HostAddress,
	}
}

func (a *RestBookApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *RestBookApp) Run() error {
	a.log.Info("starting http server on ",
		"addr:", a.hostAddr,
	)

	if err := a.restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("REST server failed", "error", err)
		return err
	}

	return nil
}

func (a *RestBookApp) Stop(ctx context.Context) error {
	if err := a.restServer.Shutdown(ctx); err != nil {
		a.log.Error("REST shutdown error", "error", err)
		return err
	}

	a.log.Info("REST server stopped")

	return nil
}
