package isuutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func exampleEcho() {
	_, err := InitializeTracerProvider()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(otelecho.Middleware(serviceName))
}

func exampleMux() {
	_, err := InitializeTracerProvider()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Use(otelmux.Middleware(serviceName))
}

func exampleGin() {
	_, err := InitializeTracerProvider()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(otelgin.Middleware(serviceName))
}
