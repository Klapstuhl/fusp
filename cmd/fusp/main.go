package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Klapstuhl/fusp/pkg/config"
	"github.com/Klapstuhl/fusp/pkg/log"
	"github.com/Klapstuhl/fusp/pkg/proxy"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancle := context.WithCancel(context.Background())

	cfg, err := config.Read()
	if err != nil {
		logrus.Fatalf("Error reading config, %s", err)
	}

	if err := log.Open(cfg.Log.FilePath, cfg.Log.Level); err != nil {
		logrus.Fatalf("Error openign log, %s", err)
	}

	proxies := make([]*proxy.Proxy, 0)

	for proxyName, proxyCfg := range cfg.Proxies {
		proxy, err := proxy.NewProxy(proxyName, proxyCfg)
		if err != nil {
			logrus.WithField("proxy", proxyName).Error(err)
			continue
		}
		proxies = append(proxies, proxy)
	}

	for _, proxy := range proxies {
		go proxy.Start()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			logrus.Info("Shuting down")
			cancle()
		case syscall.SIGUSR1:
			if err := log.Rotate(); err != nil {
				logrus.WithError(err).Errorf("Error rotating log")
			}
		}
	}()

	<-ctx.Done()
	ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	for _, proxy := range proxies {
		if err := proxy.Shutdown(ctx); err != nil {
			logrus.WithError(err).WithField("proxy", proxy.Name).Error("failed to shutdown")
		}
	}
}
