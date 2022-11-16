package otelgin

import (
	"github.com/gin-gonic/gin"
	"github.com/wang007/otel-go/metrics"
)

type RewriteStatus func(c *gin.Context, recommendStatus string) string

type RewriteActiveService func(c *gin.Context, recommendActiveService string) string

type config struct {
	httpCallCollector    *metrics.HttpCallCollector
	service              string
	RewriteStatus        RewriteStatus
	RewriteActiveService RewriteActiveService
}

type Option interface {
	apply(c *config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func WithRewriteStatus(f RewriteStatus) Option {
	return optionFunc(func(config *config) {
		config.RewriteStatus = f
	})
}

func WithRewriteActiveService(f RewriteActiveService) Option {
	return optionFunc(func(config *config) {
		config.RewriteActiveService = f
	})
}

func WithHttpCallCollector(h *metrics.HttpCallCollector) Option {
	return optionFunc(func(c *config) {
		c.httpCallCollector = h
	})
}

func WithService(service string) Option {
	return optionFunc(func(c *config) {
		c.service = service
	})
}
