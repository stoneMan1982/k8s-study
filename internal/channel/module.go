package channel

import (
	"context"

	"demo/pkg/module"
)

const channelModuleName = "channel"

// ChannelModule implements the module.Module interface.
type ChannelModule struct {
	meta module.MetaData
	svc  *Service
	api  *API
}

func NewChannelModule() *ChannelModule {
	svc := NewService()
	return &ChannelModule{
		meta: module.MetaData{
			Name:        channelModuleName,
			Version:     "v1",
			Description: "Channel management module",
			Author:      "demo",
			Kind:        "internal",
		},
		svc: svc,
		api: NewAPI(svc),
	}
}

func (m *ChannelModule) MetaData() module.MetaData {
	return m.meta
}

func (m *ChannelModule) GetService() interface{} {
	return m.api
}

func (m *ChannelModule) Start(ctx context.Context) error {
	_ = ctx
	return m.svc.Start()
}

func (m *ChannelModule) Stop(ctx context.Context) error {
	_ = ctx
	return m.svc.Stop()
}
