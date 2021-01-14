package data

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gorilla/websocket"
)

type Delta struct {
	StartTime int64   `json:"st"`
	EndTime   int64   `json:"et"`
	Price     float64 `json:"p"`
	TotalVol  float64 `json:"tv"`
	ActualVol float64 `json:"av"`
	Ratio     float64 `json:"r"`
	Trades    int     `json:"n"`
}

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

func MakeDelta(msg []byte) *Delta {
	re := regexp.MustCompile(`"[tTcvnV]":"?(\d+\.?\d+)`)
	matches := re.FindAllSubmatch(msg, -1)
	startTime, _ := strconv.ParseInt(string(matches[0][1]), 10, 64)
	endTime, _ := strconv.ParseInt(string(matches[1][1]), 10, 64)
	price, _ := strconv.ParseFloat(string(matches[2][1]), 64)
	totalVol, _ := strconv.ParseFloat(string(matches[3][1]), 64)
	actualVol, _ := strconv.ParseFloat(string(matches[5][1]), 64)
	ratio := actualVol / totalVol
	trades, _ := strconv.Atoi(string(matches[4][1]))
	return &Delta{
		StartTime: startTime,
		EndTime:   endTime,
		Price:     price,
		TotalVol:  totalVol,
		ActualVol: actualVol,
		Ratio:     ratio,
		Trades:    trades,
	}
}

func GetCurrentPrice(msg []byte) float64 {
	re := regexp.MustCompile(`"c":"?(\d+\.?\d+)`)
	matches := re.FindAllSubmatch(msg, -1)
	if len(matches) > 0 {
		price, _ := strconv.ParseFloat(string(matches[0][1]), 64)
		return price
	}
	return 0
}

func (delta *Delta) ToJSON(w io.Writer) {
	d := json.NewEncoder(w)
	err := d.Encode(delta)
	if err != nil {
		log.Panic(err)
	}
}
