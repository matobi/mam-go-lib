package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matobi/mam-go-lib/pkg/conf"
	"github.com/matobi/mam-go-lib/pkg/logger"
	"github.com/rs/zerolog/log"
)

type global struct {
	httpClient     *http.Client
	cfg            *conf.Config
	isShuttingDown bool
}

var (
	g *global
)

const (
	serviceName = "mam-golib-example"
)

func main() {
	profile := os.Getenv("profile")
	if profile == "" {
		log.Fatal().Msg("missing env 'profile'")
	}
	logger.InitLogger(serviceName, profile)

	cfg, err := initConf(profile)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}

	g = &global{
		httpClient: &http.Client{Timeout: time.Second * 10},
		cfg:        cfg,
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	httpSrv := initHTTP()
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !g.isShuttingDown {
			log.Fatal().Err(err).Msg("failed to serve http")
		}
	}()

	<-stopChan // wait for SIGINT
	g.isShuttingDown = true
	log.Info().Msg("got kill signal")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Info().Err(err).Msg("failed shutdown http server")
	}
	//wg.Wait()
	log.Info().Msg("bye")
}
