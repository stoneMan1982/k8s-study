package friend

import (
	"context"

	"demo/pkg/module"
)

const friendModuleName = "friend"

// FriendModule implements the module.Module interface.
type FriendModule struct {
	meta module.MetaData
	svc  *Service
	api  *API
}

func NewFriendModule() *FriendModule {
	svc := NewService()
	return &FriendModule{
		meta: module.MetaData{
			Name:        friendModuleName,
			Version:     "v1",
			Description: "Friend management module",
			Author:      "demo",
			Kind:        "internal",
		},
		svc: svc,
		api: NewAPI(svc),
	}
}

func (m *FriendModule) MetaData() module.MetaData {
	return m.meta
}

func (m *FriendModule) GetService() interface{} {
	return m.api
}

func (m *FriendModule) Start(ctx context.Context) error {
	_ = ctx
	return m.svc.Start()
}

func (m *FriendModule) Stop(ctx context.Context) error {
	_ = ctx
	return m.svc.Stop()
}
