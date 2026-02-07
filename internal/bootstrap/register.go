package bootstrap

import (
	"demo/internal/channel"
	"demo/internal/event"
	"demo/internal/friend"
	"demo/internal/user"
	"demo/pkg/module"
)

// RegisterInternalModules wires all internal modules into the provided registry.
func RegisterInternalModules(reg *module.Registry) {
	reg.RegisterModule(channel.NewChannelModule())
	reg.RegisterModule(event.NewEventModule())
	reg.RegisterModule(friend.NewFriendModule())
	reg.RegisterModule(user.NewUserModule())
}
