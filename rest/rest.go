package rest

import (
	"context"
	"log"
	"time"
	"{{{mytemplate}}}/json"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Current Dependency
type HandlerFunc echo.HandlerFunc
type MiddlewareFunc echo.MiddlewareFunc

func converToEchoType(h HandlerFunc, mws ...MiddlewareFunc) (echo.HandlerFunc, []echo.MiddlewareFunc) {
	handler := echo.HandlerFunc(h)
	var middlewares []echo.MiddlewareFunc
	for _, m := range mws {
		middlewares = append(middlewares, echo.MiddlewareFunc(m))

	}
	return handler, middlewares
}

type Context any

// Rest Server interface with C as the abstrast the  underlying context of web frmaework , L
type RestServer[C Context] interface {
	// Use to hide the Banner of web framework
	HideBanner()
	// Disable HTTP2 , HTTP2 is enabled by default
	DisableHTTP2()
	//  Set Readtimeout
	SetReadTimeout(t time.Duration)
	//  Set Writetimeout
	SetWriteTimeout(t time.Duration)
	// Set Binder to parse body payload
	SetBodyParser(p json.JSONSerializer[C])
	// Set Custom HTTP error handler
	SetErrorHandler(handler func(err error, c C))
	// Get Routing
	Get(route string, handler HandlerFunc, mw ...MiddlewareFunc)
	// Post Routing
	Post(route string, handler HandlerFunc, mw ...MiddlewareFunc)
	// Put Routing
	Put(route string, handler HandlerFunc, mw ...MiddlewareFunc)
	// Register Middleware
	Use(mw ...MiddlewareFunc)
	// Listen for incoming request
	Listen(addr string) error
	// Listen tls , if key type is string , then it's treated as the file path  , else if it is []byte , then its treated as content as is
	ListenTLS(addr string, key interface{}, cert interface{}) error
	// Shut down the  Server
	Shutdown(ctx context.Context) error
}

// Echo WF
type echoSV struct {
	logger *zerolog.Logger
	sv     *echo.Echo
}

func (e *echoSV) DisableHTTP2() {
	e.sv.DisableHTTP2 = true
}

func (e *echoSV) SetReadTimeout(t time.Duration) {
	e.sv.Server.ReadTimeout = t
}

func (e *echoSV) SetWriteTimeout(t time.Duration) {
	e.sv.Server.WriteTimeout = t
}

func (e *echoSV) SetBodyParser(p json.JSONSerializer[echo.Context]) {
	e.sv.JSONSerializer = p
}

func (e *echoSV) SetErrorHandler(handler func(err error, c echo.Context)) {
	e.sv.HTTPErrorHandler = handler
}

func (e *echoSV) Use(mw ...MiddlewareFunc) {
	_, mws := converToEchoType(nil, mw...)
	e.sv.Use(mws...)
}

func (e *echoSV) Get(route string, handler HandlerFunc, mws ...MiddlewareFunc) {
	h, m := converToEchoType(handler, mws...)
	e.sv.GET(route, h, m...)
}

func (e *echoSV) HideBanner() {
	e.sv.HideBanner = true
}

func (e *echoSV) Post(route string, handler HandlerFunc, mws ...MiddlewareFunc) {
	h, m := converToEchoType(handler, mws...)
	e.sv.POST(route, h, m...)
}

func (e *echoSV) Put(route string, handler HandlerFunc, mws ...MiddlewareFunc) {
	h, m := converToEchoType(handler, mws...)
	e.sv.PUT(route, h, m...)
}

func (e *echoSV) Listen(addr string) error {
	return e.sv.Start(addr)
}

func (e *echoSV) ListenTLS(addr string, key interface{}, cert interface{}) error {
	return e.sv.StartTLS(addr, key, cert)
}

func (e *echoSV) Shutdown(ctx context.Context) error {
	return e.sv.Shutdown(ctx)
}

func NewRest(l *zerolog.Logger) RestServer[echo.Context] {

	e := echo.New()
	// For safety use of echo context Logger
	log.SetOutput(l)

	return &echoSV{
		logger: l,
		sv:     e,
	}
}
