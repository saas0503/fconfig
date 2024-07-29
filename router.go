package fcore

type MethodKey string

const (
	GET    MethodKey = "GET"
	POST   MethodKey = "POST"
	PUT    MethodKey = "PUT"
	PATCH  MethodKey = "PATCH"
	DELETE MethodKey = "DELETE"
)

type Route struct {
}
