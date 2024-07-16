package middlewares

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func LoggingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Infof("Started request: method=%s uri=%s", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
			log.Infof("Completed request: method=%s uri=%s time=%v", r.Method, r.RequestURI, time.Since(start))
		})
	}
}

func Middleware1() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//before
			next.ServeHTTP(w, r)
			//after
		})
	}
}

func Middleware2() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//before
			next.ServeHTTP(w, r)
			//after
		})
	}
}

func NormalRoute() Middleware {
	return ChainMiddleware(
		LoggingMiddleware(),
		Middleware1(),
		Middleware2(),
	)
}
