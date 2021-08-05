package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func wrapHTTPLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {

			var (
				reqID  = middleware.GetReqID(r.Context())
				scheme = "http"
			)

			if r.TLS != nil {
				scheme = "https"
			}

			logrus.WithFields(logrus.Fields{
				"request_id":  reqID,
				"protocol":    r.Proto,
				"status_code": ww.Status(),
				"method":      r.Method,
			}).WithTime(time.Now()).Log(
				logrus.InfoLevel,
				fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
			)

		}()

		next.ServeHTTP(ww, r)

	})
}
