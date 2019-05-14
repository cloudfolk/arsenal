package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
)

func NewSignalActor(logger logr.Logger) ActorResult {
	return ActorResult{
		Actor: &signalActor{
			sigs:   make(chan os.Signal),
			logger: logger,
		},
	}
}

type signalActor struct {
	sigs   chan os.Signal
	logger logr.Logger
}

func (r *signalActor) Run() error {
	signal.Notify(r.sigs, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(r.sigs)
	if sig := <-r.sigs; sig != nil {
		r.logger.Info(fmt.Sprintf("Received signal %s ... shutting down", sig))
	}
	return nil
}

func (r *signalActor) Interrupt(err error) {
	close(r.sigs)
}
