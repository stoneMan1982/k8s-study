package event

import (
	"context"

	"demo/pkg/module"
)

const eventModuleName = "event"

// EventModule implements the module.Module interface.
type EventModule struct {
	meta module.MetaData
	svc  *Service
	api  *API
}

func NewEventModule() *EventModule {
	svc := NewService()
	return &EventModule{
		meta: module.MetaData{
			Name:        eventModuleName,
			Version:     "v1",
			Description: "Event management module",
			Author:      "demo",
			Kind:        "internal",
		},
		svc: svc,
		api: NewAPI(svc),
	}
}

func (m *EventModule) MetaData() module.MetaData {
	return m.meta
}

func (m *EventModule) GetService() interface{} {
	return m.api
}

func (m *EventModule) Start(ctx context.Context) error {
	_ = ctx
	return m.svc.Start()
}

func (m *EventModule) Stop(ctx context.Context) error {
	_ = ctx
	return m.svc.Stop()
}
