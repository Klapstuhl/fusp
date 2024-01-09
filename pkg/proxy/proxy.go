package proxy

import (
	"context"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/Klapstuhl/fusp/pkg/middleware"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	Name   string
	config *config.Proxy

	entrypoint *Entrypoint

	chain *middleware.Chain
}

func NewProxy(name string, config *config.Proxy, middlewares map[string]*config.Middleware) (*Proxy, error) {
	logger := logrus.WithField("proxy", name)

	logger.Debug("Creating Socket")
	socket, err := NewSocket(config.Socket, name)
	if err != nil {
		return nil, err
	}

	logger.Debug("Creating Chain")
	chain, err := middleware.NewChain(config.Middleware, middlewares, socket)
	if err != nil {
		return nil, err
	}

	logger.Debug("Setting entypoint")
	entrypoint, err := NewEntrypoint(config.Entrypoint, chain.Start())
	if err != nil {
		return nil, err
	}

	return &Proxy{Name: name, config: config, entrypoint: entrypoint, chain: chain}, nil
}

func (p *Proxy) Start(ctx context.Context) {
	logrus.WithField("proxy", p.Name).Debug("Proxy started")
	p.entrypoint.Strat(ctx)
}

func (p *Proxy) Close() {
	p.entrypoint.Close()
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	logrus.WithField("proxy", p.Name).Debug("Proxy stopped")
	return p.entrypoint.Shutdown(ctx)
}
