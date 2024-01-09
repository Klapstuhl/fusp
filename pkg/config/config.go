package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Read() (*Config, error) {
	var cfg Config

	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/fusp/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err := filepath.WalkDir("./config.d/", func(path string, entry os.DirEntry, err error) error {
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
	if err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Config struct {
	Proxies    map[string]*Proxy
	Middleware map[string]*Middleware
	Sockets    map[string]*Socket
	Log        Log
	Providers  *Provider
}

type Proxy struct {
	Entrypoint *Entrypoint
	Socket     string
	Middleware []string
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

type Provider struct {
	File *FileProvider
}

type FileProvider struct {
	Directory string
}
