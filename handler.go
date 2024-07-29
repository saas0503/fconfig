package fcore

type Middleware func(next Ctx) error
