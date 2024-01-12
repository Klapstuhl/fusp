package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Read() (*Config, error) {
	var cfg config

	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/fusp/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err := filepath.WalkDir("/etc/fusp/config.d/", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && entry.Type().IsRegular() {
			if strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml") {
				viper.SetConfigFile(path)
				if err := viper.MergeInConfig(); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg.toRuntime()
}

type config struct {
	Proxies    map[string]*proxy
	Middleware map[string]*Middleware
	Sockets    map[string]*Socket
	Log        Log
}

func (cfg *config) toRuntime() (*Config, error) {
	proxies := make(map[string]*Proxy)

	for name, proxy := range cfg.Proxies {
		socket, ok := cfg.Sockets[proxy.Socket]
		if !ok {
			return nil, fmt.Errorf("unkown socket %s", proxy.Socket)
		}

		middleware := make(map[string]*Middleware)

		for _, mName := range proxy.Middleware {
			mCfg, ok := cfg.Middleware[mName]
			if !ok {
				return nil, fmt.Errorf("unkown middleware %s", mName)
			}

			middleware[mName] = mCfg
		}

		proxies[name] = &Proxy{proxy.Entrypoint, socket.Path, middleware}
	}

	return &Config{proxies, cfg.Log}, nil
}

type Config struct {
	Proxies map[string]*Proxy
	Log     Log
}

type proxy struct {
	Entrypoint *Entrypoint
	Socket     string
	Middleware []string
}

type Proxy struct {
	Entrypoint *Entrypoint
	Socket     string
	Middleware map[string]*Middleware
}

type Entrypoint struct {
	Type    string
	Address string
}

type Socket struct {
	Path string
}

type Middleware struct {
	MethodFilter   *MethodFilter
	EndpointFilter *EndpointFilter
	ReplacePath    *ReplacePath
	AddPrefix      *AddPrefix
	IPFilter       *IPFilter
	StripPrefix    *StripPrefix
}

type AddPrefix struct {
	Prefix string
}

type EndpointFilter struct {
	Block     bool
	Endpoints []string
}

type IPFilter struct {
	Block       bool
	Addresses   []string
	ExcludedIPs []string
	Depth       int
	Header      string
}

type MethodFilter struct {
	Block   bool
	Methods []string
}

type ReplacePath struct {
	Regex       string
	Replacement string
}

type StripPrefix struct {
	Prefix []string
}

type Log struct {
	Level    string
	FilePath string
}

// type Provider struct {
// 	File *FileProvider
// }

// type FileProvider struct {
// 	Directory string
// }
