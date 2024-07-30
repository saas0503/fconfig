package fcore

import (
	"net/http"
	"sync"
)

type App struct {
	pool sync.Pool
}

func New() *App {
	app := &App{}

	app.pool = sync.Pool{
		New: func() interface{} {
			return NewCtx(app)
		},
	}

	return app
}

func (app *App) AcquireCtx(w http.ResponseWriter, r *http.Request) Ctx {
	ctx, ok := app.pool.Get().(Ctx)
	if !ok {
		panic("acquire ctx fail")
	}
	ctx.Reset(w, r)
	return ctx
}

func (app *App) ReleaseCtx(ctx Ctx) {
	ctx.release()
	app.pool.Put(ctx)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if app.convertMethod(interface{}(method).(MethodKey)) == -1 {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (app *App) convertMethod(s MethodKey) int {
	switch s {
	case GET:
		return 0
	case POST:
		return 1
	case PUT:
		return 2
	case DELETE:
		return 3
	case PATCH:
		return 4
	default:
		return -1
	}
}
