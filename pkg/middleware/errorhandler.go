package middleware

import (
	"context"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrorHandler struct{}

func (e *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, err error) {
	statusCode := http.StatusInternalServerError

	if e, ok := err.(net.Error); ok {
		if e.Timeout() {
			statusCode = http.StatusGatewayTimeout
		} else {
			statusCode = http.StatusBadGateway
		}
	} else if err == context.Canceled {
		statusCode = http.StatusTeapot
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))

	log.Error(statusCode, " ", http.StatusText(statusCode), " caused by: ", err)
}
