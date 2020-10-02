package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anonutopia/gowaves"
	ui18n "github.com/unknwon/i18n"
)

const priceFactorLimit = uint64(0.0001 * float64(satInBtc))

// TokenMonitor represents issuing token
type TokenMonitor struct {
	PriceRecord uint64
}

// func (t *TokenMonitor) issueAmount(investment int, assetID string, dryRun bool) (int, float64) {
// 	amount := int(0)

// 	oldPrice := t.Price
// 	oldPriceFactor := t.PriceFactor
// 	oldTierPrice := t.TierPrice
// 	oldTierPriceFactor := t.TierPriceFactor

// 	p, err := pc.DoRequest()
// 	if err != nil {
// 		log.Printf("[token.issueAmount] error pc.DoRequest: %s", err)
// 		logTelegram(fmt.Sprintf("[token.issueAmount] error pc.DoRequest: %s", err))
// 		return 0, 0
// 	}

// 	var cryptoPrice float64
// 	var investmentEur float64

// 	log.Println(assetID)

// 	if assetID == "" {
// 		cryptoPrice = p.WAVES
// 	} else if assetID == "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS" {
// 		cryptoPrice = p.BTC
// 	} else if assetID == "474jTeYx2r2Va35794tCScAXWJG9hU2HcgxzMowaZUnu" {
// 		cryptoPrice = p.ETH
// 	} else {
// 		return amount, float64(t.Price) / float64(satInBtc)
// 	}

// 	priceChanged := false
// 	investmentEur = float64(investment) / float64(satInBtc) / cryptoPrice

// 	log.Println(investmentEur)

// 	for investment > 10 {
// 		log.Printf("token: %d %d %d %d", t.Price, t.PriceFactor, t.TierPrice, t.TierPriceFactor)
// 		log.Printf("cryptoPrice: %f", cryptoPrice)
// 		log.Printf("investment: %d", investment)

// 		tierAmount := uint64(float64(investment) / cryptoPrice / float64(t.Price) * float64(satInBtc))

// 		log.Printf("tierAmount: %d", tierAmount)

// 		if tierAmount > t.TierPrice {
// 			tierAmount = t.TierPrice
// 		}

// 		log.Printf("tierAmount: %d", tierAmount)

// 		tierInvestment := int(float64(tierAmount) * float64(t.Price) * cryptoPrice / float64(satInBtc))

// 		log.Printf("tierInvestment: %d", tierInvestment)

// 		amount = amount + int(tierAmount)

// 		log.Printf("amount: %d", amount)

// 		investment = investment - tierInvestment

// 		log.Printf("investment: %d", investment)

// 		t.TierPrice = t.TierPrice - tierAmount
// 		t.TierPriceFactor = t.TierPriceFactor - tierAmount

// 		log.Printf("token: %d %d %d %d", t.Price, t.PriceFactor, t.TierPrice, t.TierPriceFactor)

// 		if t.TierPrice == 0 {
// 			t.TierPrice = 1000 * satInBtc
// 			t.Price = t.Price + t.PriceFactor
// 			priceChanged = true
// 		}

// 		if t.TierPriceFactor == 0 {
// 			t.TierPriceFactor = 1000000 * satInBtc
// 			if t.PriceFactor > priceFactorLimit {
// 				t.PriceFactor = t.PriceFactor - priceFactorLimit
// 			}
// 		}

// 		if !dryRun {
// 			t.saveState()
// 			log.Printf("token: %d %d %d %d", t.Price, t.PriceFactor, t.TierPrice, t.TierPriceFactor)
// 		}
// 	}

// 	newPrice := float64(0)

// 	if priceChanged {
// 		newPrice = float64(t.Price) / float64(satInBtc)
// 		log.Println(newPrice)
// 	}

// 	if !dryRun {
// 		sendInvestmentMessages(investmentEur, newPrice)
// 	} else {
// 		t.Price = oldPrice
// 		t.PriceFactor = oldPriceFactor
// 		t.TierPrice = oldTierPrice
// 		t.TierPriceFactor = oldTierPriceFactor
// 	}

// 	log.Println(dryRun)

// 	return amount, newPrice
// }

func (t *TokenMonitor) saveState() {
	ksip := &KeyValue{Key: "tokenPriceRecord"}
	db.FirstOrCreate(ksip, ksip)
	ksip.ValueInt = t.PriceRecord
	if err := db.Save(ksip).Error; err != nil {
		logTelegram("[token.go - 121] " + err.Error())
	}
}

func (t *TokenMonitor) loadState() {
	ksip := &KeyValue{Key: "tokenPriceRecord"}
	db.FirstOrCreate(ksip, ksip)

	if ksip.ValueInt > 0 {
		t.PriceRecord = ksip.ValueInt
	}
}

func (t *TokenMonitor) start() {
	for {
		price := t.getPrice()
		priceInt := uint64(price * float64(satInBtc))

		if price > (float64(t.PriceRecord)/float64(satInBtc) + 0.0005) {
			t.PriceRecord = priceInt
			t.saveState()
			msg := fmt.Sprintf(ui18n.Tr(lang, "priceRise"), price)
			msgHr := fmt.Sprintf(ui18n.Tr(langHr, "priceRise"), price)
			if conf.Dev {
				messageTelegram(msg, tAnonOps)
			} else {
				messageTelegram(msg, tAnon)
				messageTelegram(msgHr, tAnonBalkan)
			}
		}

		t.checkLastOrder()

		time.Sleep(time.Second * 30)
	}
}

func (t *TokenMonitor) getPrice() float64 {
	osr, err := wmc.OrderbookStatus(conf.TokenID, "WAVES")
	if err != nil {
		logTelegram("[token.go - 150]" + err.Error())
	}

	p, err := pc.DoRequest()
	if err != nil {
		log.Println("[token.go - 155]" + err.Error())
		logTelegram("[token.go - 156]" + err.Error())
	}

	price := float64(osr.LastPrice) / float64(satInBtc) / p.WAVES

	return price
}

func (t *TokenMonitor) checkLastOrder() {
	osr, err := wmc.OrderbookStatus(conf.TokenID, "WAVES")
	if err != nil {
		logTelegram("[token.go - 150]" + err.Error())
		return
	}

	if osr.LastSide != "buy" {
		price := osr.Ask
		amount := 1

		order := &gowaves.AssetsOrderRequest{
			SenderPublicKey:  "2F84usFJtN7devqZJJDu5WVRVfeNMTfMijXsQo2MN5a5",
			MatcherPublicKey: "9cpfKN9suPNvfeUNphzxXMjcnn974eme8ZhWUjaktzU5",
			AssetPair: struct {
				AmountAsset string `json:"amountAsset"`
				PriceAsset  string `json:"priceAsset"`
			}{
				AmountAsset: conf.TokenID,
				PriceAsset:  "WAVES",
			},
			OrderType:  "buy",
			Amount:     amount,
			Price:      price,
			MatcherFee: 300000,
			Version:    3,
		}

		if aor, err := wnc.AssetsOrder(order); err != nil {
			logTelegram("[token.go - 214]" + err.Error())
		} else {
			log.Println(aor.Signature)
		}
	}
}

func initTokenMonitor() *TokenMonitor {
	tm := &TokenMonitor{
		PriceRecord: uint64(0),
	}

	tm.loadState()

	go tm.start()

	return tm
}
