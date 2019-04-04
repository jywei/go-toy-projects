package main

import (
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/honestbee/habitat-catalog-cache/procesor"
	"github.com/honestbee/habitat-catalog-cache/seeker"
	"github.com/rs/zerolog"
)

type healthCheck struct {
	listener net.Listener
	done     chan struct{}
	logger   *zerolog.Logger
}

func newHealthCheck(logger *zerolog.Logger, addr string) (*healthCheck, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	h := &healthCheck{
		listener: l,
		done:     make(chan struct{}),
		logger:   logger,
	}

	go h.start()
	return h, nil
}

func (h *healthCheck) start() {
	for {
		select {
		case <-h.done:
			return
		default:
			conn, err := h.listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					h.logger.Error().Err(err).Msgf("failed accept a connection request")
				}
				continue
			}
			conn.Close()
		}
	}
}

func (h *healthCheck) close() error {
	close(h.done)
	return h.listener.Close()
}

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	conf := config.NewSyncWorker()

	service, err := model.New(conf.Database, conf.Cache)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new service failed")
	}
	defer service.Close()

	seeker, err := seeker.New(conf.Seeker)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new seeker failed")
	}

	proc, err := procesor.New(conf.Procesor, service, seeker)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new procesor failed")
	}

	hc, err := newHealthCheck(&logger, conf.TCPHealthCheckListenAddr)
	if err != nil {
		logger.Fatal().Err(err).Msgf("new health check failed")
	}
	defer hc.close()

	ticker := time.NewTicker(time.Duration(conf.SyncCheckPeriodSec) * time.Second)
	defer ticker.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msgf("sync worker start")
LOOP:
	for {
		select {
		case <-ticker.C:
			job, err := service.PopJob()
			if err != nil {
				switch err {
				case model.ErrNoJob:
				default:
					logger.Error().Msgf("pop job failed:%v", err)
				}
				continue LOOP
			}

			logger.Info().Msgf("receive a job, type:%d, value:%v", job.Type, job.Value)
			if err := proc.Process(job); err != nil {
				logger.Error().Err(err).Msgf("sync process on job:%v failed", job)
			}

		case <-done:
			break LOOP
		}
	}

	logger.Info().Msgf("sync worker exit")
}
