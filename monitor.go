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
	amount, _ := token.issueAmount(t.Amount, t.AssetID, false)
	log.Println(amount)
	user := &User{Address: t.Sender}
	db.First(user, user)

	atr := &gowaves.AssetsTransferRequest{
		Amount:    amount,
		AssetID:   conf.TokenID,
		Fee:       100000,
		Recipient: t.Sender,
		Sender:    conf.NodeAddress,
	}

	_, err := wnc.AssetsTransfer(atr)
	if err != nil {
		log.Printf("[purchaseAsset] error assets transfer: %s", err)
		logTelegram(fmt.Sprintf("[purchaseAsset] error assets transfer: %s", err))
	} else {
		log.Printf("Sent token: %s => %d", t.Sender, amount)
		amount = t.Amount - 200000
		amountR := int(float64(amount) * 0.2)
		amountF := int(float64(amount) * 0.8)

		r := &User{}
		db.First(r, user.ReferralID)

		atr = &gowaves.AssetsTransferRequest{
			Amount:    amountR,
			AssetID:   t.AssetID,
			Fee:       100000,
			Recipient: r.Address,
			Sender:    conf.NodeAddress,
		}

		_, err := wnc.AssetsTransfer(atr)
		if err != nil {
			log.Printf("[purchaseAsset] error Waves referral transfer: %s", err)
			logTelegram(fmt.Sprintf("[purchaseAsset] error Waves referral transfer: %s", err))
		} else {
			log.Printf("Sent waves referral: %s => %d", r.Address, amountR)

			atr = &gowaves.AssetsTransferRequest{
				Amount:    amountF,
				AssetID:   t.AssetID,
				Fee:       100000,
				Recipient: conf.FounderAddress,
				Sender:    conf.NodeAddress,
			}

			_, err := wnc.AssetsTransfer(atr)
			if err != nil {
				log.Printf("[purchaseAsset] error Waves founder transfer: %s", err)
				logTelegram(fmt.Sprintf("[purchaseAsset] error Waves founder transfer: %s", err))
			} else {
				log.Printf("Sent waves founder: %s => %d", conf.FounderAddress, amountF)
			}
		}
	}
}

func initMonitor() {
	wm = &WavesMonitor{}
	wm.start()
}
