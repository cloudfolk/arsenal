package main

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/oklog/run"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewRunGroup() *run.Group {
	return &run.Group{}
}

func KubernetesConfig() (*rest.Config, error) {
	return config.GetConfig()
}

func NewManager(cfg *rest.Config, logger logr.Logger) (_ manager.Manager, err error) {
	// Create a new Cmd to provide shared dependencies and start components
	return manager.New(cfg, manager.Options{
		MetricsBindAddress: fmt.Sprintf("%s:%d", viper.GetString("metrics.host"), viper.GetInt32("metrics.port")),
	})
}
