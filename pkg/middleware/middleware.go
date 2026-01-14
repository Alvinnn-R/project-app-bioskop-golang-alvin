package middleware

import (
	"go.uber.org/zap"
)

type MiddlewareCostume struct {
	Log *zap.Logger
}

func NewMiddlewareCustome(log *zap.Logger) MiddlewareCostume {
	return MiddlewareCostume{
		Log: log,
	}
}
