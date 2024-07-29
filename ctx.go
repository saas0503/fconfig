package fcore

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Ctx interface {
	BaseURL() string
	BodyParser(interface{}) error
	Get(key interface{}) interface{}
	IP() string
	JSON(data interface{}) error
	Locals(key interface{}, value interface{})
	Next() error
	Params(key string, defaultValue ...string) string
	Queries() map[string]string
	QueryBool(key string, defaultValue ...bool) bool
	QueryFloat(key string, defaultValue ...float64) float64
	QueryInt(key string, defaultValue ...int64) int64
	Req() *http.Request
	Res() http.ResponseWriter
}

type CtxImpl struct {
	req *http.Request
	res http.ResponseWriter
}

// BaseURL implements Ctx.
func (ctx *CtxImpl) BaseURL() string {
	return ctx.req.URL.Scheme + "://" + ctx.req.Host
}

// BodyParser implements Ctx.
func (ctx *CtxImpl) BodyParser(payload interface{}) error {
	err := json.NewDecoder(ctx.req.Body).Decode(&payload)
	return err
}

// Get implements Ctx.
func (ctx *CtxImpl) Get(key interface{}) interface{} {
	return ctx.req.Context().Value(key)
}

// IP implements Ctx.
func (ctx *CtxImpl) IP() string {
	ipAddress := ctx.req.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = ctx.req.RemoteAddr
	}

	return ipAddress
}

// JSON implements Ctx.
func (ctx *CtxImpl) JSON(data interface{}) error {
	ctx.res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(ctx.res).Encode(data)
	return err
}

// Locals implements Ctx.
func (ctx *CtxImpl) Locals(key interface{}, value interface{}) {
	contx := context.WithValue(ctx.req.Context(), key, value)
	ctx.req = ctx.req.WithContext(contx)
}

// Next implements Ctx.
func (ctx *CtxImpl) Next() error {
	return nil
}

// Params implements Ctx.
func (ctx *CtxImpl) Params(key string, defaultValue ...string) string {
	para := ctx.req.PathValue(key)
	if para == "" {
		para = defaultValue[0]
	}

	return para
}

// Queries implements Ctx.
func (ctx *CtxImpl) Queries() map[string]string {
	query := make(map[string]string)
	q := ctx.req.URL.Query()
	for k, v := range q {
		query[k] = strings.Join(v, "")
	}

	return query
}

func (ctx *CtxImpl) Query(key string, defaultValue ...string) string {
	q := ctx.req.URL.Query().Get(key)
	if q == "" {
		q = defaultValue[0]
	}

	return q
}

// QueryBool implements Ctx.
func (ctx *CtxImpl) QueryBool(key string, defaultValue ...bool) bool {
	q := ctx.req.URL.Query().Has(key)
	if q {
		query := ctx.req.URL.Query().Get(key)
		boolValue, err := strconv.ParseBool(query)
		if err != nil {
			panic(err)
		}
		return boolValue
	}
	return defaultValue[0]
}

// QueryFloat implements Ctx.
func (ctx *CtxImpl) QueryFloat(key string, defaultValue ...float64) float64 {
	q := ctx.req.URL.Query().Has(key)
	if q {
		query := ctx.req.URL.Query().Get(key)
		floatValue, err := strconv.ParseFloat(query, 64)
		if err != nil {
			panic(err)
		}
		return floatValue
	}
	return defaultValue[0]
}

// QueryInt implements Ctx.
func (ctx *CtxImpl) QueryInt(key string, defaultValue ...int64) int64 {
	q := ctx.req.URL.Query().Has(key)
	if q {
		query := ctx.req.URL.Query().Get(key)
		intValue, err := strconv.ParseInt(query, 0, 32)
		if err != nil {
			panic(err)
		}
		return intValue
	}
	return defaultValue[0]
}

// Req implements Ctx.
func (ctx *CtxImpl) Req() *http.Request {
	return ctx.req
}

// Res implements Ctx.
func (ctx *CtxImpl) Res() http.ResponseWriter {
	return ctx.res
}

func NewCtx(w http.ResponseWriter, r *http.Request) Ctx {
	return &CtxImpl{
		req: r,
		res: w,
	}
}
