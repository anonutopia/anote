package main

import (
	"time"

	"github.com/anonutopia/gowaves"
)

type SustainingService struct {
	LastSell time.Time
	LastSend time.Time
}

func (sus *SustainingService) start() {
	sus.LastSell = time.Now()
	sus.LastSend = time.Now()
	for {
		sus.checkState()

		time.Sleep(time.Second * 10)
	}
}

func (sus *SustainingService) checkState() {
	if abr, err := wnc.AddressesBalance(conf.NodeAddress); err != nil {
		return
	} else {
		if abr.Balance < 10000000 && time.Since(sus.LastSell) > time.Duration(time.Minute*5) {
			sus.sell()
		} else if abr.Balance < 30000000 {
			sus.createLimitOrder()
		}
	}

	if time.Since(sus.LastSend) > time.Duration(time.Minute*5) {
		sus.sendToNode()
	}
}

func (sus *SustainingService) createLimitOrder() {
	mlp := &KeyValue{Key: "myLastPrice"}
	db.FirstOrCreate(mlp, mlp)

	osr, err := wmc.OrderbookStatus(conf.TokenID, "WAVES")
	if err != nil {
		logTelegram("[sustaining.go - 39]" + err.Error())
		return
	}

	if mlp.ValueInt == 0 || mlp.ValueInt > uint64(osr.Ask) {
		price := osr.Ask - 1
		amount := 30000000 * satInBtc / uint64(price)

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
			OrderType:  "sell",
			Amount:     int(amount),
			Price:      price,
			MatcherFee: 300000,
			Version:    3,
			Timestamp:  int(time.Now().UnixNano() / int64(time.Millisecond)),
			Expiration: int(time.Now().Add(time.Hour*48).UnixNano() / int64(time.Millisecond)),
		}

		if aor, err := wnc.AssetsOrder(order); err != nil {
			logTelegram("[sustaining.go - 66]" + err.Error())
		} else {
			if _, err := wmc.Orderbook(aor); err != nil {
				logTelegram("[sustaining.go - 69]" + err.Error())
			}
		}

		mlp.ValueInt = uint64(price)
		db.Save(mlp)
	}
}

func (sus *SustainingService) sell() {
	opr, err := wmc.OrderbookPair(conf.TokenID, "WAVES", 10)
	if err != nil {
		logTelegram("[sustaining.go - 39]" + err.Error())
		return
	}

	waves := uint64(0)
	price := uint64(0)
	amount := uint64(0)
	leftToBuy := uint64(0)

	for i := 0; waves < 10000000; i++ {
		w := opr.Bids[i].Amount * opr.Bids[i].Price / satInBtc
		waves += w
		price = opr.Bids[i].Price
		if waves < 10000000 {
			amount = opr.Bids[i].Amount
		} else {
			leftToBuy = 10000000 - waves + w
			newAmount := leftToBuy * satInBtc / opr.Bids[i].Price
			amount += newAmount
		}
	}

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
		OrderType:  "sell",
		Amount:     int(amount),
		Price:      int(price),
		MatcherFee: 300000,
		Version:    3,
		Timestamp:  int(time.Now().UnixNano() / int64(time.Millisecond)),
		Expiration: int(time.Now().Add(time.Hour).UnixNano() / int64(time.Millisecond)),
	}

	if aor, err := wnc.AssetsOrder(order); err != nil {
		logTelegram("[sustaining.go - 214]" + err.Error())
	} else {
		if _, err := wmc.Orderbook(aor); err != nil {
			logTelegram("[sustaining.go - 217]" + err.Error())
		} else {
			sus.LastSell = time.Now()
		}
	}
}

func (sus *SustainingService) sendToNode() {
	if abr, err := wnc.AddressesBalance("3PPc3AP75DzoL8neS4e53tZ7ybUAVxk2jAb"); err != nil {
		return
	} else if abr.Balance >= 20000000 {
		atr := &gowaves.AssetsTransferRequest{
			Amount:     abr.Balance - 10000000,
			Fee:        100000,
			Recipient:  conf.NodeAddress,
			Sender:     "3PPc3AP75DzoL8neS4e53tZ7ybUAVxk2jAb",
			Attachment: "fees",
		}

		if _, err := wnc.AssetsTransfer(atr); err != nil {
			logTelegram("[sustaining.go - 144]" + err.Error())
		} else {
			sus.LastSend = time.Now()
		}
	}
}

func initSustainingService() {
	sus := &SustainingService{}
	if !conf.Dev {
		go sus.start()
	}
}
