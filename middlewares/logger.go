package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWrite{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		log.Printf("-> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("<- %s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

type loggingResponseWrite struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWrite) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
