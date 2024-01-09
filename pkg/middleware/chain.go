package middleware

import (
	"fmt"
	"net/http"

	"github.com/Klapstuhl/fusp/pkg/config"
)

type Constructor func(http.Handler) (http.Handler, error)

type Chain struct {
	constructors []Constructor
	start        http.Handler
	end          http.Handler
}

func NewChain(middleares []string, configs map[string]*config.Middleware, end http.Handler) (*Chain, error) {
	if end == nil {
		return nil, fmt.Errorf("invalid chain end handler")
	}

	constructors := make([]Constructor, 0)

	for _, name := range middleares {
		cfg, ok := configs[name]
		if !ok {
			return nil, fmt.Errorf("unknown middleware '%s'", name)
		}

		constructor, err := buildConstructor(name, cfg)
		if err != nil {
			return nil, err
		}

		constructors = append(constructors, constructor)
	}

	chain := &Chain{constructors, nil, end}
	if err := chain.assemble(); err != nil {
		return nil, err
	}

	return chain, nil
}

func buildConstructor(name string, cfg *config.Middleware) (Constructor, error) {
	if cfg.AddPrefix != nil {
		return func(next http.Handler) (http.Handler, error) {
			return NewAddPrefix(*cfg.AddPrefix, name, next)
		}, nil
	}

	if cfg.EndpointFilter != nil {
		return func(next http.Handler) (http.Handler, error) {
			return NewEndpointFilter(*cfg.EndpointFilter, name, next)
		}, nil
	}

	if cfg.IPFilter != nil {
		return func(next http.Handler) (http.Handler, error) {
			return NewIPFilter(*cfg.IPFilter, name, next)
		}, nil
	}

	if cfg.MethodFilter != nil {
		return func(next http.Handler) (http.Handler, error) {
			return NewMethodFilter(*cfg.MethodFilter, name, next)
		}, nil
	}

	if cfg.ReplacePath != nil {
		return func(next http.Handler) (http.Handler, error) {
			return NewReplacePath(*cfg.ReplacePath, name, next)
		}, nil
	}

	if cfg.StripPrefix != nil {
		return func(next http.Handler) (http.Handler, error) {
			return NewStripPrefix(*cfg.StripPrefix, name, next)
		}, nil
	}

	return nil, fmt.Errorf("'%s' unkown middleware type", name)
}

func (c *Chain) assemble() error {
	var err error
	next := c.end
	for i := len(c.constructors) - 1; i >= 0; i-- {
		next, err = c.constructors[i](next)
		if err != nil {
			return err
		}
	}
	c.start = next
	return nil
}

func (c *Chain) Start() http.Handler {
	return c.start
}

func (c *Chain) Append(constructors ...Constructor) error {
	c.constructors = append(c.constructors, constructors...)

	err := c.assemble()
	return err
}
