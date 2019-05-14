package app

import (
	"go.uber.org/dig"
)

type Logger interface {
	Info(fmt string, args ...interface{})
}

type Actor interface {
	Run() error
	Interrupt(error)
}

type ActorResult struct {
	dig.Out
	Actor Actor `group:"actors"`
}

type ActorsResult struct {
	dig.In
	Actors []Actor `group:"actors"`
}
