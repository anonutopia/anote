package main

import (
	"log"

	"github.com/anonutopia/gowaves"
)

func initWaves() (*gowaves.WavesNodeClient, *gowaves.WavesMatcherClient) {
	wnc := &gowaves.WavesNodeClient{
		Host:   conf.NodeHost,
		Port:   6869,
		ApiKey: conf.WavesNodeAPIKey,
	}

	wmc := &gowaves.WavesMatcherClient{
		Host: "https://matcher.waves.exchange",
		Port: 443,
	}

	return wnc, wmc
}

func send() {
	atr := &gowaves.AssetsTransferRequest{
		Amount:    1100000000,
		AssetID:   "",
		Fee:       100000,
		Recipient: "3P2EtZMgEN4W49hLXy966D53oHiE52gawhn",
		Sender:    "3PJySTACVDWXFFzVFMPSSzAK3XHfDbekHc4",
	}

	_, err := wnc.AssetsTransfer(atr)
	log.Println(err)
}
