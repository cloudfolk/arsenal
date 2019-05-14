package app

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/cloudfolk/arsenal/pkg/apis"
	"github.com/cloudfolk/arsenal/pkg/controller"
	"github.com/cloudfolk/arsenal/pkg/webhook"
)

func NewManagerActor(mgr manager.Manager, logger logr.Logger) ActorResult {
	return ActorResult{
		Actor: &managerActor{
			stopCh: make(chan struct{}),
			mgr:    mgr,
		},
	}
}

type managerActor struct {
	stopCh chan struct{}
	logger logr.Logger
	mgr    manager.Manager
}

func (m *managerActor) Run() error {
	// Setup Scheme for all resources
	if err := apis.AddToScheme(m.mgr.GetScheme()); err != nil {
		m.logger.Error(err, "unable to add APIs to scheme")
		return err
	}

	// Setup all Controllers
	if err := controller.AddToManager(m.mgr); err != nil {
		m.logger.Error(err, "unable to register controllers to the manager")
		return err
	}

	if err := webhook.AddToManager(m.mgr); err != nil {
		m.logger.Error(err, "unable to register webhooks to the manager")
		return err
	}

	return m.mgr.Start(m.stopCh)
}

func (m *managerActor) Interrupt(err error) {
	close(m.stopCh)
}
