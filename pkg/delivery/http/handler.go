package http

import (
	"context"
	"net/http"

	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/domain"
	"github.com/labstack/echo/v4"
)

// Handler aggregate all service
type Handler struct {
	rateLimiter domain.RateLimiterService
}

// NewHandler ...
func NewHandler(rateLimiter domain.RateLimiterService) *Handler {
	return &Handler{
		rateLimiter: rateLimiter,
	}
}

// ProtectedEndpoint entry protected endpoint
// method = GET url = api/v1/protected
func (h *Handler) ProtectedEndpoint(c echo.Context) error {
	var (
		ctx  context.Context
		addr string
		url  string
	)

	addr = c.RealIP()
	url = c.Path()
	ctx = c.Request().Context()

	claims, err := h.rateLimiter.RequireResource(ctx, addr, url)
	if err != nil {
		return err
	}

	// in here do you want service

	return c.JSON(http.StatusOK, claims)
}
