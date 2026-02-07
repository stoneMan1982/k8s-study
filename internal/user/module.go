package user

import (
	"context"

	"demo/pkg/module"
)

type UserModule struct {
	meta module.MetaData
	svc  *Service
	api  *API
}

const userModuleName = "user"

func NewUserModule() *UserModule {
	svc := NewService()
	return &UserModule{
		meta: module.MetaData{
			Name:        userModuleName,
			Version:     "v1",
			Description: "User management module",
			Author:      "demo",
			Kind:        "internal",
		},
		svc: svc,
		api: NewAPI(svc),
	}
}

func (u *UserModule) MetaData() module.MetaData {
	return u.meta
}

func (u *UserModule) GetService() interface{} {
	return u.api
}

func (u *UserModule) Start(ctx context.Context) error {
	_ = ctx
	return u.svc.Start()
}

func (u *UserModule) Stop(ctx context.Context) error {
	_ = ctx
	return u.svc.Stop()
}
