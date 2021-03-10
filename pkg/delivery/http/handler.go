package http

import (
	"context"
	"net/http"

	"github.com/karta0898098/dcard-rate-limiter/pkg/ratelimiter/domain"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	rateLimiter domain.RateLimiterService
}

func NewHandler(rateLimiter domain.RateLimiterService) *Handler {
	return &Handler{
		rateLimiter: rateLimiter,
	}
}

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

	return c.JSON(http.StatusOK, claims)
}
