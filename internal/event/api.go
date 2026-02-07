package event

// API exposes the event module capabilities to other modules.
type API struct {
	svc *Service
}

func NewAPI(svc *Service) *API {
	return &API{svc: svc}
}

func (a *API) Status() string {
	if a.svc.Started() {
		return "running"
	}
	return "stopped"
}
