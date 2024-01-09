package proxy

import (
	"context"
	"fmt"
	stdlog "log"
	"net"
	"net/http"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/sirupsen/logrus"
)

type Entrypoint struct {
	server   *http.Server
	listener net.Listener
}

func NewEntrypoint(config *config.Entrypoint, handler http.Handler) (*Entrypoint, error) {
	var network string
	if config.Type == "unix" {
		network = "unix"
	} else if config.Type == "http" {
		network = "tcp"
	} else {
		return nil, fmt.Errorf("Unknown entrypoint type: %q", config.Type)
	}

	listener, err := net.Listen(network, config.Address)
	if err != nil {
		return nil, err
	}

	serverLogger := stdlog.New(logrus.StandardLogger().WriterLevel(logrus.ErrorLevel), "", 0)
	server := &http.Server{Handler: handler, ErrorLog: serverLogger}

	return &Entrypoint{server: server, listener: listener}, nil
}

func (e *Entrypoint) Strat(ctx context.Context) {
	e.server.Serve(e.listener)
}

func (e *Entrypoint) Close() {
	e.server.Close()
}

func (e *Entrypoint) Shutdown(ctx context.Context) error {
	return e.server.Shutdown(ctx)
}
