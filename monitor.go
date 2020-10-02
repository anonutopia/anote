package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

// WavesMonitor represents waves monitoring object
type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		// todo - make sure that everything is ok with 100 here
		pages, err := wnc.TransactionsAddressLimit(conf.NodeAddress, 100)
		if err != nil {
			log.Println(err)
			logTelegram("[wm.start - wnc.TransactionsAddressLimit]" + err.Error())
		}

		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransaction(&t)
			}
		}

		time.Sleep(time.Second)
	}
}

func (wm *WavesMonitor) checkTransaction(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxID: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != true {
		wm.processTransaction(&tr, t)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, t *gowaves.TransactionsAddressLimitResponse) {
	if t.Recipient == conf.NodeAddress {
		log.Println(t)
	}
	if t.Type == 4 &&
		t.Timestamp >= wm.StartedTime &&
		t.Sender != conf.NodeAddress &&
		t.Recipient == conf.NodeAddress &&
		(t.AssetID == "" ||
			t.AssetID == "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS" ||
			t.AssetID == "474jTeYx2r2Va35794tCScAXWJG9hU2HcgxzMowaZUnu") &&
		len(t.Attachment) == 0 {

		log.Println("purchase")
		wm.purchaseAsset(t)
	}

	tr.Processed = true
	if err := db.Save(tr).Error; err != nil {
		logTelegram("[wm.processTransaction - db.Save] " + err.Error())
	}
}

func (wm *WavesMonitor) purchaseAsset(t *gowaves.TransactionsAddressLimitResponse) {
	order := &gowaves.AssetsOrderRequest{}
	log.Println(order)
}

func initMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
