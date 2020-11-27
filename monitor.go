package main

import (
	"fmt"
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
				wm.checkTransactionAnote(&t)
			}
		}

		pages, err = wnc.TransactionsAddressLimit(conf.AintAddress, 100)
		if err != nil {
			log.Println(err)
			logTelegram("[wm.start - wnc.TransactionsAddressLimit]" + err.Error())
		}

		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransactionAint(&t)
			}
		}

		time.Sleep(time.Second * 30)
	}
}

func (wm *WavesMonitor) checkTransactionAnote(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxID: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != true {
		if t.Type == 4 &&
			t.Timestamp >= wm.StartedTime &&
			t.Sender != conf.NodeAddress &&
			t.Recipient == conf.NodeAddress &&
			(t.AssetID == "" ||
				t.AssetID == "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS" ||
				t.AssetID == "474jTeYx2r2Va35794tCScAXWJG9hU2HcgxzMowaZUnu") &&
			len(t.Attachment) == 0 {

			wm.purchaseAnote(t)
		}

		tr.Processed = true
		if err := db.Save(tr).Error; err != nil {
			logTelegram("[wm.processTransaction - db.Save] " + err.Error())
		}
	}
}

func (wm *WavesMonitor) checkTransactionAint(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxID: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != true {
		if t.Type == 4 &&
			t.Timestamp >= wm.StartedTime &&
			t.Sender != conf.AintAddress &&
			t.Recipient == conf.AintAddress &&
			(t.AssetID == "" ||
				t.AssetID == "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS" ||
				t.AssetID == "474jTeYx2r2Va35794tCScAXWJG9hU2HcgxzMowaZUnu") {

			wm.purchaseAint(t)
		}

		tr.Processed = true
		if err := db.Save(tr).Error; err != nil {
			logTelegram("[wm.processTransaction - db.Save] " + err.Error())
		}
	}
}

func (wm *WavesMonitor) purchaseAnote(t *gowaves.TransactionsAddressLimitResponse) {
	u := &User{Address: t.Sender}
	db.First(u, u)
	messageTelegram(fmt.Sprintf("ANOTE purchase: %s - %.8f", t.Sender, float64(t.Amount)/float64(satInBtc)), tAnonTeam)
	messageTelegram("We have received your coins and we'll resolve your trade in next 24 hours.", int64(u.TelegramID))
}

func (wm *WavesMonitor) purchaseAint(t *gowaves.TransactionsAddressLimitResponse) {
	// priceAint := float64(tm.PriceRecord) / float64(satInBtc) * 24 * 365
	priceAint := float64(tm.Price) / float64(satInBtc) * 24 * 365
	// priceAint := 1.44
	prices, err := pc.DoRequest()
	if err != nil {
		logTelegram("[monitor.go - 104]" + err.Error())
		return
	}
	invEur := float64(0)

	if t.AssetID == "" {
		invEur = float64(t.Amount) / prices.WAVES / float64(satInBtc)
	} else if t.AssetID == "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS" {
		invEur = float64(t.Amount) / prices.BTC / float64(satInBtc)
	} else if t.AssetID == "474jTeYx2r2Va35794tCScAXWJG9hU2HcgxzMowaZUnu" {
		invEur = float64(t.Amount) / prices.ETH / float64(satInBtc)
	}

	amount := invEur / priceAint
	amountInt := int(amount * float64(satInBtc))

	atr := &gowaves.AssetsTransferRequest{
		Amount:    amountInt,
		AssetID:   conf.AintID,
		Fee:       100000,
		Recipient: t.Sender,
		Sender:    conf.AintAddress,
	}

	if _, err := wnc.AssetsTransfer(atr); err != nil {
		logTelegram("[monitor.go - 129]" + err.Error())
	}

	messageTelegram(fmt.Sprintf("AINT purchase: %.8f €", invEur), tAnonTeam)
}

func initMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
