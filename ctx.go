package fcore

import (
	"mime/multipart"
	"net/http"
)

type CtxImpl struct {
}

type Ctx interface {
	BaseURL() string
	Body() []byte
	BodyParser(interface{}) error
	FormFile(string) (*multipart.FileHeader, error)
	Get(key string, defaultValue ...string) string
	IP() string
	IPs() []string
	JSON(da interface{}, ctype ...string) error
	Locals(key interface{}, value ...interface{}) interface{}
	MultipartForm() (*multipart.Form, error)
	Next() error
	Params(key string, defaultValue ...string) string
	Queries() map[string]string
	Query(key string, defaultValue ...string) string
	QueryBool(key string, defaultValue ...bool) bool
	QueryFloat(key string, defaultValue ...float64) float64
	QueryInt(key string, defaultValue ...int) int
	QueryParser(interface{}) error
	Redirect(url string, status ...int) error
	Req() *http.Request
	Res() http.ResponseWriter
}
