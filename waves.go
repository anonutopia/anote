package main

import (
	"log"

	"github.com/anonutopia/gowaves"
)

func initWaves() *gowaves.WavesNodeClient {
	wnc := &gowaves.WavesNodeClient{
		Host:   conf.NodeHost,
		Port:   6869,
		ApiKey: conf.WavesNodeAPIKey,
	}

	return wnc
}

func send() {
	atr := &gowaves.AssetsTransferRequest{
		Amount:    200000,
		AssetID:   "474jTeYx2r2Va35794tCScAXWJG9hU2HcgxzMowaZUnu",
		Fee:       100000,
		Recipient: "3P2EtZMgEN4W49hLXy966D53oHiE52gawhn",
		Sender:    "3PJySTACVDWXFFzVFMPSSzAK3XHfDbekHc4",
	}

	_, err := wnc.AssetsTransfer(atr)
	log.Println(err)
}
