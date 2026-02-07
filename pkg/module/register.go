package module

import "context"

// Registry keeps track of modules and controls their lifecycle.
type Registry struct {
	modules map[string]Module
}

func NewRegistry() *Registry {
	return &Registry{modules: make(map[string]Module)}
}

func (r *Registry) RegisterModule(m Module) {
	r.modules[m.MetaData().Name] = m
}

func (r *Registry) GetModule(name string) (Module, bool) {
	m, ok := r.modules[name]
	return m, ok
}

func (r *Registry) GetModules() []Module {
	var ms []Module
	for _, m := range r.modules {
		ms = append(ms, m)
	}
	return ms
}

func (r *Registry) StartModules(ctx context.Context) error {
	for _, m := range r.modules {
		if err := m.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) StopModules(ctx context.Context) error {
	for _, m := range r.modules {
		if err := m.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) GetModuleService(name string) (interface{}, bool) {
	m, ok := r.modules[name]
	if !ok {
		return nil, false
	}
	return m.GetService(), true
}

// DefaultRegistry keeps backward-compatible global helpers.
var DefaultRegistry = NewRegistry()

func RegisterModule(m Module) {
	DefaultRegistry.RegisterModule(m)
}

func GetModule(name string) (Module, bool) {
	return DefaultRegistry.GetModule(name)
}

func GetModules() []Module {
	return DefaultRegistry.GetModules()
}

func StartModules(ctx context.Context) error {
	return DefaultRegistry.StartModules(ctx)
}

func StopModules(ctx context.Context) error {
	return DefaultRegistry.StopModules(ctx)
}

func GetModuleService(name string) (interface{}, bool) {
	return DefaultRegistry.GetModuleService(name)
}

func GetServiceFromRegistry[T any](reg *Registry, name string) (T, bool) {
	var zero T
	if reg == nil {
		return zero, false
	}
	m, ok := reg.modules[name]
	if !ok {
		return zero, false
	}
	svc, ok := m.GetService().(T)
	if !ok {
		return zero, false
	}
	return svc, true
}

func GetService[T any](name string) (T, bool) {
	return GetServiceFromRegistry[T](DefaultRegistry, name)
}
