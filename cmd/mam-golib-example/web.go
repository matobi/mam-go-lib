package main

import (
	"fmt"
	"net/http"

	"github.com/matobi/mam-go-lib/pkg/version"
)

func initHTTP() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/healthcheck", version.CreateHealthHandler(serviceName))
	port := fmt.Sprintf(":%d", g.cfg.Int(confPort))
	srv := &http.Server{Addr: port, Handler: mux}
	return srv
}
