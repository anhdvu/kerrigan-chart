package data

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kerrigan-chart/config"
	"kerrigan-chart/util"
	"log"
	"sort"
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
		log.Printf("ERROR: Error during JSON sentry history unmarshaling\nError detail: %v\n", err)
		log.Println(string(raw))
		return
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

func (ss *Sentries) GetCurrentSentry() sentry {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.d[len(ss.d)-1]
}

func (ss *Sentries) GetTrend(h uint) {
	var nCount, pCount, zCount int
	var trend int
	// var TrendMax, TrendMin int
	// TrendMax = int(h)*6 - 1
	// TrendMin = -1 * TrendMax

	ss.mu.Lock()
	s := ss.d[len(ss.d)-int(h)*6:]
	ss.mu.Unlock()

	for i := 0; i < len(s)-1; i++ {
		if s[i+1].Value-s[i].Value > 6 {
			trend++
			pCount++
		} else if s[i+1].Value-s[i].Value < -6 {
			trend--
			nCount++
		} else {
			trend = trend + 0
			zCount++
		}
	}

	countSlice := []int{nCount, pCount, zCount}
	sort.Ints(countSlice)
	max := countSlice[len(s)-1]
	switch max {
	case nCount:
		fmt.Println("Count: bear")
	case pCount:
		fmt.Println("Count: bull")
	case zCount:
		fmt.Println("Count: sideway")
	}

	trendDiff := s[len(s)-1].Value - s[0].Value
	if trendDiff > float64(h*36) {
		fmt.Println("Trend difference: bull")
	} else if trendDiff < float64(-1*int(h*36)) {
		fmt.Println("Trend difference: bear")
	} else {
		fmt.Println("Trend difference: sideway")
	}
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
		log.Printf("ERROR: Error during JSON sentry prediction unmarshaling\nError detail: %v\n", err)
		return
	}
}

func (sp *SentryPredictions) GetClosestFutureSentry() (*SentryPrediction, error) {
	now := time.Now().Unix()
	if len(sp.d) > 0 {
		for _, e := range sp.d {
			if util.ToEpoch(e.Time) < now {
				log.Printf("There was a re-prediction at %v", e.Time)
			} else if util.ToEpoch(e.Time) > now {
				return &e, nil
			}
		}
	} else {
		log.Println("checker.txt is empty for now.")
	}
	return &SentryPrediction{}, fmt.Errorf("There is no valid prediction at the moment.")
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
