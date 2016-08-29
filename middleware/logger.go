package middleware

import (
	"log"
	"net/http"
	"time"
	"youtube/utils/logs"
)

type Logger struct {
	*log.Logger
}

func NewLogger() func(http.Handler) http.Handler {
	l := Logger{logs.New("HTTP")}
	return l.factory
}

func (l Logger) factory(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
			start := time.Now()
			l.Printf("Started %v %v", r.Method, r.URL.Path)

			next.ServeHTTP(w, r)

			l.Printf("Completed %v %v in %v", r.Method, r.URL.Path, time.Since(start))
		})
}
