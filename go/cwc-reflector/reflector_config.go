package main

import (
	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

type ReflectorConfig struct {
	CWCAddress string
	WebAddress string
}

var defaultConfig = ReflectorConfig{
	CWCAddress: "0.0.0.0:7388",
	WebAddress: "0.0.0.0:7380",
}

func ReadConfig(filename string) *ReflectorConfig {
	cfg := defaultConfig

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		glog.Warningf("Reflector Config file not found %s", filename)
	}

	return &cfg
}
