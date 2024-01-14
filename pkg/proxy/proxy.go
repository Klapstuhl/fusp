package proxy

import (
	"context"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/Klapstuhl/fusp/pkg/middleware"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	Name       string
	socketPath string
	entrypoint *Entrypoint

	chain *middleware.Chain
}

func NewProxy(name string, cfg *config.Proxy) (*Proxy, error) {
	logger := logrus.WithField("proxy", name)

	logger.Debug("Creating Socket")
	socket, err := NewSocket(cfg.Socket)
	if err != nil {
		return nil, err
	}

	logger.Debug("Creating Chain")
	chain, err := middleware.NewChain(cfg.Middleware, socket)
	if err != nil {
		return nil, err
	}

	logger.Debug("Setting entypoint")
	entrypoint, err := NewEntrypoint(cfg.Entrypoint, chain.Start())
	if err != nil {
		return nil, err
	}

	return &Proxy{Name: name, socketPath: cfg.Socket, entrypoint: entrypoint, chain: chain}, nil
}

func (p *Proxy) Start() {
	logrus.WithField("proxy", p.Name).Debug("Proxy started")
	p.entrypoint.Strat()
}

func (p *Proxy) Close() {
	p.entrypoint.Close()
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	logrus.WithField("proxy", p.Name).Debug("Proxy stopped")
	return p.entrypoint.Shutdown(ctx)
}
