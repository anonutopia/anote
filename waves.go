package main

import (
	"encoding/json"
	"log"
	"time"

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
		Amount: 72000000,
		// AssetID:   "8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS",
		AssetID:    "",
		Fee:        100000,
		Recipient:  "3P2EtZMgEN4W49hLXy966D53oHiE52gawhn",
		Sender:     "3PLyWsbNa96tLGZb1dhX6D1tvGhUE9FBf8F",
		Attachment: "fees",
	}

	_, err := wnc.AssetsTransfer(atr)
	log.Println(err)
}

func create_order() string {
	price := 217810
	amount := int(satInBtc) / 100
	aor := &gowaves.AssetsOrderResponse{}
	var err error

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
		Amount:     amount,
		Price:      price,
		MatcherFee: 300000,
		Version:    3,
		Timestamp:  int(time.Now().UnixNano() / int64(time.Millisecond)),
		Expiration: int(time.Now().Add(time.Hour).UnixNano() / int64(time.Millisecond)),
	}

	if aor, err = wnc.AssetsOrder(order); err != nil {
		logTelegram("[waves.go - 63]" + err.Error())
	} else {
		if _, err := wmc.Orderbook(aor); err != nil {
			logTelegram("[token.go - 66]" + err.Error())
		} else {
			b, _ := json.Marshal(aor)
			logTelegram(string(b))
		}
	}

	log.Println(aor.ID)

	return aor.ID
}

func cancel_order(orderId string) {
	// encoding := base58.BitcoinEncoding

	// log.Println(orderId)

	// // log.Println(encoding)

	// decodedKey, _ := encoding.Decode([]byte("2F84usFJtN7devqZJJDu5WVRVfeNMTfMijXsQo2MN5a5"))
	// decodedOrder, _ := encoding.Decode([]byte(orderId))

	// signature := append([]byte("2F84usFJtN7devqZJJDu5WVRVfeNMTfMijXsQo2MN5a5")[:], []byte(orderId)[:]...)

	// signature, _ = encoding.Encode(signature)
	// log.Println(string(signature))

	// encSig, _ := encoding.Encode([]byte("2KzxaA6FAZRr32R1pMunkPMGKCfbNvhPbiv5W4a3TmtXJtfpMnjMVp1SebqgoBMUXJzozP1tFC4RLT3Ewb31f8DW"))

	cor := &gowaves.OrderbookCancelRequest{
		Sender:    "2F84usFJtN7devqZJJDu5WVRVfeNMTfMijXsQo2MN5a5",
		OrderID:   orderId,
		Signature: "286BCQqrk3KDJ78kHupmJmakRpXXDvD7AcAHkxNzFbSTsEGysY6LLL8m4Wy7jcxf5DeDSuLCqSFzazivLrb5HRLW",
	}

	cors, err := wmc.OrderbookCancel(conf.TokenID, "WAVES", cor)

	if err != nil {
		log.Println(err.Error())
		log.Println(cors)
	}
}
