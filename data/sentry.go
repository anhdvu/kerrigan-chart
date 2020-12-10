package data

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"kerrigan-chart/config"
	"kerrigan-chart/util"
	"log"
	"sync"
	"time"
)

type WsMsg struct {
	M string `json:"m"`
	D struct {
		T int64   `json:"t"`
		V float64 `json:"v"`
	} `json:"d"`
}
type sentryRecord struct {
	Time       string
	Prediction float64 `json:"pred_price"`
	Actual     float64 `json:"actual_price;omitempty"`
}

type sentry struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}
type Sentries struct {
	mu sync.Mutex
	d  []sentry
}

// Sentries data type has 3 methods
// Get() to retrieve current data
// Update() to update current data from file historical_delta.txt
// ToJSON() to render json to a writer interface e.g. http.ResponseWriter

func (ss *Sentries) Get() []sentry {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.d
}

func (ss *Sentries) Update() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	raw, err := ioutil.ReadFile(config.HistorySentryFile)
	if err != nil {
		log.Panicf("PANIC: Error reading file %v\nError detail: %v\n", config.HistorySentryFile, err)
	}

	data := make([]sentryRecord, 0)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Panicf("PANIC: Error during JSON sentry history unmarshaling\nError detail: %v\n", err)
	}

	ss.d = make([]sentry, len(data))
	for i, e := range data {
		ss.d[i].Time = util.ToEpoch(e.Time)
		ss.d[i].Value = e.Prediction
	}
}

func (ss *Sentries) ToJSON(w io.Writer) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	d := json.NewEncoder(w)
	err := d.Encode(ss.d)
	if err != nil {
		log.Panic(err)
	}
}

func (ss *Sentries) GetCurrentSentryValue() float64 {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.d[len(ss.d)-1].Value
}

type SentryPrediction struct {
	Time       string
	Prediction float64 `json:"predict"`
}

type SentryPredictions struct {
	mu sync.Mutex
	d  []SentryPrediction
}

// SentryPredictions data type has 3 methods
// Get() to retrieve current data
// Update() to update current data from file checker.txt
// GetClosestFutureSentry() to retrieve the most upcoming sentry prediction

func (sp *SentryPredictions) Get() []SentryPrediction {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.d
}

func (sp *SentryPredictions) Update() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	raw, err := ioutil.ReadFile(config.SentryPredictionFile)
	if err != nil {
		log.Panicf("PANIC: Error reading file %v\nError detail: %v\n", config.SentryPredictionFile, err)
	}
	err = json.Unmarshal(raw, &sp.d)
	if err != nil {
		log.Panicf("PANIC: Error during JSON sentry prediction unmarshaling\nError detail: %v\n", err)
	}
}

func (sp *SentryPredictions) GetClosestFutureSentry() *SentryPrediction {
	now := time.Now().Unix()
	if len(sp.d) > 0 {
		for _, e := range sp.d {
			if util.ToEpoch(e.Time) < now {
				log.Printf("There was a re-prediction at %v", e.Time)
			} else if util.ToEpoch(e.Time) > now {
				return &e
			}
		}
	} else {
		log.Println("checker.txt is empty for now.")
	}
	return &SentryPrediction{}
}

// SentryPrediction data type has 1 method
// ToWSMessage() is a convenient method to convert a sentry prediction to a websocket message
func (sp *SentryPrediction) ToWSMessage() *WsMsg {
	msg := &WsMsg{}
	msg.M = "sentry"
	msg.D.T = util.ToEpoch(sp.Time)
	msg.D.V = sp.Prediction
	return msg
}

// Reserved for future uses
func handlepanic(s string) {
	if err := recover(); err != nil {
		log.Println("RECOVER", err)
	}
}
