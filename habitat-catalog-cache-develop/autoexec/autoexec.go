package autoexec

import (
	"time"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// AutoExec is a self-trigger executor to push sync-up job at specific time.
type AutoExec interface {
	Close() error
}

type autoexec struct {
	done    chan struct{}
	ticker  *time.Ticker
	service model.Service
	logger  *zerolog.Logger
}

// New returns a AutoExec instance.
func New(conf *config.AutoExec, service model.Service, logger *zerolog.Logger) (AutoExec, error) {
	a := &autoexec{
		ticker:  time.NewTicker(time.Duration(conf.SyncupPeriodSec) * time.Second),
		service: service,
		logger:  logger,
		done:    make(chan struct{}),
	}
	go a.start()

	return a, nil
}

func (a *autoexec) start() {
	for {
		select {
		case <-a.done:
			a.logger.Info().Msgf("autoexec: exit")
			return
		case <-a.ticker.C:
			if err := a.exec(); err != nil {
				a.logger.Error().Err(err).Msgf("autoexec: [start] exec failed")
			}
		}
	}
}

func (a *autoexec) exec() error {
	keys, err := a.service.GetExternalKeys()
	if err != nil {
		return errors.Wrapf(err, "autoexec: [exec] GetExternalKeys failed")
	}

	keyMap := make(map[string]struct{})
	for _, key := range keys {
		if _, exist := keyMap[key]; !exist {
			err = a.service.PushJob(&model.Job{
				Type:  model.ExternalKeyJob,
				Value: key,
			})
			if err != nil {
				return errors.Wrapf(err, "autoexec: [exec] push external key:%s job failed", key)
			}

			keyMap[key] = struct{}{}
		}
	}

	return nil
}

func (a *autoexec) Close() error {
	close(a.done)
	a.ticker.Stop()
	return nil
}
