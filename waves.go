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
		Amount: 10000000,
		// AssetID:   "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS",
		AssetID:   "",
		Fee:       100000,
		Recipient: "3PJySTACVDWXFFzVFMPSSzAK3XHfDbekHc4",
		Sender:    "3PLyWsbNa96tLGZb1dhX6D1tvGhUE9FBf8F",
	}

	_, err := wnc.AssetsTransfer(atr)
	log.Println(err)
}
