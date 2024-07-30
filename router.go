package fcore

import (
	"fmt"
	"github.com/saas0503/futils"
	"sync/atomic"
)

type MethodKey string

const (
	GET    MethodKey = "GET"
	POST   MethodKey = "POST"
	PUT    MethodKey = "PUT"
	PATCH  MethodKey = "PATCH"
	DELETE MethodKey = "DELETE"
)

type Handler func(Ctx) error

type Route struct {
	pos   uint32
	star  bool
	root  bool
	mount bool

	// Path data
	path         string
	routerParser routeParser
	Params       []string

	// Group data
	group *Group

	// Public Data
	Path     string
	Method   MethodKey
	Handlers []Handler
}

func (app *App) register(methods []string, pathRaw string, group *Group, handler Handler, middleware ...Handler) {
	handlers := middleware
	if handler != nil {
		handlers = append(handlers, handler)
	}

	for _, method := range methods {
		method = futils.ToUpper(method)
		if app.convertMethod(interface{}(method).(MethodKey)) == -1 {
			panic(fmt.Sprintf("Invalid method: %s", method))
		}

		// is mounted app
		isMount := group != nil && group.app != app
		// A route requires atleast one ctx handler
		if len(handlers) == 0 && !isMount {
			panic(fmt.Sprintf("missing handler/middleware in route: %s\n", pathRaw))
		}
		// Cannot have an empty path
		if pathRaw == "" {
			pathRaw = "/"
		}
		// Path always start with a '/'
		if pathRaw[0] != '/' {
			pathRaw = "/" + pathRaw
		}
		// Create a stripped path in-case sensitive / trailing slashes
		pathPretty := pathRaw
		// Is path a direct wildcard?
		isStar := pathPretty == "/*"
		// Is path a root slash?
		isRoot := pathPretty == "/"
		// Parse path parameters
		parsedRaw := parseRoute(pathRaw)
		parsedPretty := parseRoute(pathPretty)

		route := Route{
			mount: isMount,
			star:  isStar,
			root:  isRoot,

			path:         RemoveEscapeChar(pathPretty),
			routerParser: parsedPretty,
			Params:       parsedRaw.params,

			group: group,

			Path:     pathRaw,
			Method:   interface{}(method).(MethodKey),
			Handlers: handlers,
		}

		atomic.AddUint32(&app.handlersCount, uint32(len(handlers)))
		app.addRoute(method, &route, isMount)
	}
}

func (app *App) addRoute(method string, route *Route, isMounted ...bool) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	// Check mounted routes
	var mounted bool
	if len(isMounted) > 0 {
		mounted = isMounted[0]
	}

	// Get unique HTTP method identifier
	md := app.convertMethod(interface{}(method).(MethodKey))

	// prevent identically route registration
	l := len(app.stack[md])
	if l > 0 && app.stack[md][l-1].Path == route.Path && !route.mount && !app.stack[md][l-1].mount {
		preRoute := app.stack[md][l-1]
		preRoute.Handlers = append(preRoute.Handlers, route.Handlers...)
	} else {
		// Increment global route position
		route.pos = atomic.AddUint32(&app.routesCount, 1)
		route.Method = interface{}(method).(MethodKey)
		// Add route to the stack
		app.stack[md] = append(app.stack[md], route)
		app.routesRefreshed = true
	}

	// Execute onRoute hooks & change latestRoute if not adding mounted route
	if !mounted {
		app.latestRoute = route
		if err := app.hooks.executeOnRouteHooks(*route); err != nil {
			panic(err)
		}
	}
}
