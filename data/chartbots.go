package data

import (
	"encoding/json"
	"io"
	"kerrigan-chart/config"
	"log"
	"os"
	"sync"
)

type bot struct {
	Name   string  `json:"name"`
	Trades []trade `json:"trades"`
}

type trade struct {
	Time   string  `json:"time"`
	Price  float64 `json:"price"`
	Action string  `json:"action"`
}

type BotTradeRecords struct {
	mu sync.Mutex
	d  []bot
}

func (btr *BotTradeRecords) Get() []bot {
	btr.mu.Lock()
	defer btr.mu.Unlock()
	return btr.d
}

func (btr *BotTradeRecords) Update() {
	btr.mu.Lock()
	defer btr.mu.Unlock()

	raw, err := os.ReadFile(config.MultiSaTradeRecords)
	if err != nil {
		log.Panicf("PANIC: Error reading file %v\nError detail: %v\n", config.MultiSaTradeRecords, err)
	}

	data := make([]bot, 0)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Printf("ERROR: Error during JSON sentry history unmarshaling\nError detail: %v\n", err)
		log.Println(string(raw))
		return
	}

	btr.d = make([]bot, len(data))
	for i, e := range data {
		btr.d[i].Name = e.Name
		btr.d[i].Trades = e.Trades
	}
}

func (btr *BotTradeRecords) ToJSON(w io.Writer) {
	btr.mu.Lock()
	defer btr.mu.Unlock()
	d := json.NewEncoder(w)
	err := d.Encode(btr.d)
	if err != nil {
		log.Panic(err)
	}
}
