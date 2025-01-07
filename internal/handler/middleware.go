package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// LoggingMiddleware расширенный вариант
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			_ = level.Debug(logger).Log(
				"msg", "received request",
				"request", fmt.Sprintf("%+v", request),
			)
			response, err := next(ctx, request)
			if err != nil {
				_ = level.Error(logger).Log("error", err, "request", fmt.Sprintf("%+v", request))
			} else {
				_ = level.Debug(logger).Log("msg", "handled request", "response", fmt.Sprintf("%+v", response))
			}
			return response, err
		}
	}
}

func HTTPLoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = level.Debug(logger).Log(
				"msg", "received request",
				"method", r.Method,
				"url", r.URL.String(),
			)
			next.ServeHTTP(w, r)
			_ = level.Debug(logger).Log(
				"msg", "handled request",
				"method", r.Method,
				"url", r.URL.String(),
			)
		})
	}
}

// HeaderMiddleware для установки заголовков
func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
