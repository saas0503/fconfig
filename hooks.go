package fcore

// ListenData is a struct to use it with OnListenHandler
type ListenData struct {
	Host string
	Port string
	TLS  bool
}

// OnRouteHandler Handlers define a function to create hooks for Fiber.
type (
	OnRouteHandler     = func(Route) error
	OnNameHandler      = OnRouteHandler
	OnGroupHandler     = func(Group) error
	OnGroupNameHandler = OnGroupHandler
	OnListenHandler    = func(ListenData) error
	OnShutdownHandler  = func() error
	OnForkHandler      = func(int) error
	OnMountHandler     = func(*App) error
)

// Hooks is a struct to use it with App.
type Hooks struct {
	// Embed app
	app *App

	// Hooks
	onRoute     []OnRouteHandler
	onName      []OnNameHandler
	onGroup     []OnGroupHandler
	onGroupName []OnGroupNameHandler
	onListen    []OnListenHandler
	onShutdown  []OnShutdownHandler
	onFork      []OnForkHandler
	onMount     []OnMountHandler
}

func (h *Hooks) executeOnRouteHooks(route Route) error {
	// Check mounting
	if h.app.mountFields.mountPath != "" {
		route.path = h.app.mountFields.mountPath + route.path
		route.Path = route.path
	}

	for _, v := range h.onRoute {
		if err := v(route); err != nil {
			return err
		}
	}

	return nil
}
