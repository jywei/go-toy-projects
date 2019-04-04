package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/honestbee/habitat-catalog-cache/autoexec"
	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/handler"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	conf := config.NewCatalogService()

	service, err := model.New(conf.Database, conf.Cache)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new service failed")
	}
	defer service.Close()

	autoexec, err := autoexec.New(conf.AutoExec, service, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new auto exec failed")
	}
	defer autoexec.Close()

	mux, err := handler.New(conf.BasicAuth, &logger, service)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new router failed")
	}

	srv := &http.Server{
		Addr:         conf.HTTP.ListenAddr,
		Handler:      mux,
		ReadTimeout:  time.Duration(conf.HTTP.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(conf.HTTP.WriteTimeoutSec) * time.Second,
		IdleTimeout:  time.Duration(conf.HTTP.IdleTimeoutSec) * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msgf("http server error")
		}
	}()

	logger.Info().Msgf("catalog cache server started")

	<-done
	logger.Info().Msgf("catalog cache server stoped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msgf("https server shutdown failed")
	}

	logger.Info().Msgf("catalog cache server exit")
}
