package middleware

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"time"
)

var logger = logrus.New()

func init(){
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func LoggerMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Infof("Method: %s | Route: %s | IP: %s | Time: %s",
			r.Method, r.URL.Path, r.RemoteAddr, duration)
	})
}