package main

import (
	"github.com/anonutopia/gowaves"
)

func initWaves() *gowaves.WavesNodeClient {
	wnc := &gowaves.WavesNodeClient{
		Host:   "162.0.225.82",
		Port:   6869,
		ApiKey: conf.WavesNodeAPIKey,
	}

	return wnc
}
