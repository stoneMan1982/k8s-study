package module

import "context"

type MetaData struct {
	Name        string
	Version     string
	Description string
	Author      string
	Kind        string
}

type Module interface {
	MetaData() MetaData
	GetService() interface{}

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
