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
	PriceRecord  uint64
	Price        uint64
	MiningPower  uint64
	TotalSupply  uint64
	TotalHolders uint64
	TotalMiners  uint64
	ActiveMiners uint64
}

func (t *TokenMonitor) saveState() {
	ksip := &KeyValue{Key: "tokenPriceRecord"}
	db.FirstOrCreate(ksip, ksip)
	ksip.ValueInt = t.PriceRecord
	if err := db.Save(ksip).Error; err != nil {
		logTelegram("[token.go - 28] " + err.Error())
	}

	tp := &KeyValue{Key: "tokenPrice"}
	db.FirstOrCreate(tp, tp)
	tp.ValueInt = t.Price
	if err := db.Save(tp).Error; err != nil {
		logTelegram("[token.go - 35] " + err.Error())
	}

	mp := &KeyValue{Key: "miningPower"}
	db.FirstOrCreate(mp, mp)
	mp.ValueInt = t.MiningPower
	if err := db.Save(mp).Error; err != nil {
		logTelegram("[token.go - 42] " + err.Error())
	}

	ts := &KeyValue{Key: "totalSupply"}
	db.FirstOrCreate(ts, ts)
	ts.ValueInt = t.TotalSupply
	if err := db.Save(ts).Error; err != nil {
		logTelegram("[token.go - 49] " + err.Error())
	}

	hld := &KeyValue{Key: "totalHolders"}
	db.FirstOrCreate(hld, hld)
	hld.ValueInt = t.TotalHolders
	if err := db.Save(hld).Error; err != nil {
		logTelegram("[token.go - 56] " + err.Error())
	}

	tm := &KeyValue{Key: "totalMiners"}
	db.FirstOrCreate(tm, tm)
	tm.ValueInt = t.TotalMiners
	if err := db.Save(tm).Error; err != nil {
		logTelegram("[token.go - 63] " + err.Error())
	}

	am := &KeyValue{Key: "activeMiners"}
	db.FirstOrCreate(am, am)
	am.ValueInt = t.ActiveMiners
	if err := db.Save(am).Error; err != nil {
		logTelegram("[token.go - 72] " + err.Error())
	}
}

func (t *TokenMonitor) loadState() {
	ksip := &KeyValue{Key: "tokenPriceRecord"}
	db.FirstOrCreate(ksip, ksip)
	if ksip.ValueInt > 0 {
		t.PriceRecord = ksip.ValueInt
	}

	tp := &KeyValue{Key: "tokenPrice"}
	db.FirstOrCreate(tp, tp)
	if tp.ValueInt > 0 {
		t.Price = tp.ValueInt
	}

	mp := &KeyValue{Key: "miningPower"}
	db.FirstOrCreate(mp, mp)
	if mp.ValueInt > 0 {
		t.MiningPower = mp.ValueInt
	}

	ts := &KeyValue{Key: "totalSupply"}
	db.FirstOrCreate(ts, ts)
	if ts.ValueInt > 0 {
		t.TotalSupply = ts.ValueInt
	}

	hld := &KeyValue{Key: "totalHolders"}
	db.FirstOrCreate(hld, hld)
	if hld.ValueInt > 0 {
		t.TotalHolders = hld.ValueInt
	}

	tm := &KeyValue{Key: "totalMiners"}
	db.FirstOrCreate(tm, tm)
	if tm.ValueInt > 0 {
		t.TotalMiners = tm.ValueInt
	}

	am := &KeyValue{Key: "activeMiners"}
	db.FirstOrCreate(am, am)
	if am.ValueInt > 0 {
		t.ActiveMiners = am.ValueInt
	}
}

func (t *TokenMonitor) start() {
	go func() {
		for {
			t.miningPower()

			t.saveState()

			time.Sleep(time.Minute * 30)
		}
	}()

	for {
		price := t.getPrice()
		priceInt := uint64(price * float64(satInBtc))
		t.Price = uint64(price * float64(satInBtc))

		if price > (float64(t.PriceRecord)/float64(satInBtc) + 0.0005) {
			t.PriceRecord = priceInt

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

		t.calculateSupply()

		t.countMiners()

		t.saveState()

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
		amount := int(satInBtc)

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
			Timestamp:  int(time.Now().UnixNano() / int64(time.Millisecond)),
			Expiration: int(time.Now().Add(time.Hour).UnixNano() / int64(time.Millisecond)),
		}

		if aor, err := wnc.AssetsOrder(order); err != nil {
			logTelegram("[token.go - 214]" + err.Error())
		} else {
			if _, err := wmc.Orderbook(aor); err != nil {
				logTelegram("[token.go - 217]" + err.Error())
			}
		}
	}
}

func (t *TokenMonitor) miningPower() {
	var mp float64
	var users []*User
	db.Where("mining = true").Find(&users)
	for _, u := range users {
		mp += u.miningPower()
	}
	t.MiningPower = uint64(mp * 100)
}

func (t *TokenMonitor) countMiners() {
	var users []*User
	db.Where("mining_activated is not null").Find(&users).Count(&t.TotalMiners)
	db.Where("mining = true").Find(&users).Count(&t.ActiveMiners)
}

func (t *TokenMonitor) calculateSupply() {
	supply := uint64(0)
	holders := uint64(0)

	if nsr, err := wnc.NodeStatus(); err != nil {
		logTelegram("[token.go - 186]" + err.Error())
	} else {
		if abdr, err := wnc.AssetsBalanceDistribution(conf.TokenID, nsr.BlockchainHeight-3, 100, ""); err != nil {
			logTelegram("[token.go - 189]" + err.Error())
		} else {
			for a, i := range abdr.Items {
				if !stringInSlice(a, conf.Exclude) {
					supply += uint64(i)
					holders++
				}
			}

			for abdr.HasNext {
				if abdr, err = wnc.AssetsBalanceDistribution(conf.TokenID, nsr.BlockchainHeight-3, 100, abdr.LastItem); err != nil {
					logTelegram("[token.go - 104]" + err.Error())
				} else {
					for a, i := range abdr.Items {
						if !stringInSlice(a, conf.Exclude) {
							supply += uint64(i)
							holders++
						}
					}
				}
			}

			t.TotalSupply = supply
			t.TotalHolders = holders
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
