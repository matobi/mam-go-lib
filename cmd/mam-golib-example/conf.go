package main

import (
	"github.com/matobi/mam-go-lib/pkg/conf"
)

const (
	confPort = "port"
)

func initConf(profile string) (*conf.Config, error) {
	c := conf.NewConfig(profile)
	c.Add(conf.VtInt, confPort, "8080")
	return c.LogAndValidate()
}
