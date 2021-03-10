package http

import "github.com/labstack/echo/v4"

// SetupRoute setup restful api
func SetupRoute(route *echo.Echo, handler *Handler) {
	api := route.Group("/api")
	{
		v1 := api.Group("/v1")
		v1.GET("/protected", handler.ProtectedEndpoint)
	}
}
