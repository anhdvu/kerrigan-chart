package data

import (
	"kerrigan-chart/util"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var historicalDeltaChannel = make(chan bool, 1)

func NewKlineWebSocket() *websocket.Conn {
	u := url.URL{
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   "/ws/btcusdt@kline_1m",
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func GetDynamicDelta() {
	util.WatchFile("abc", historicalDeltaChannel, 5)
}
