package data

import (
	"encoding/json"
	"fmt"
	"io"
	"kerrigan-chart/util"
	"log"
	"os"
	"sync"
	"time"
)

type SentryPrediction struct {
	Time       string  `json:"time"`
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

func (sp *SentryPredictions) Update(f string) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	raw, err := os.ReadFile(f)
	if err != nil {
		log.Panicf("PANIC: Error reading file %v\nError detail: %v\n", f, err)
		return err
	}

	err = json.Unmarshal(raw, &sp.d)
	if err != nil {
		log.Printf("ERROR: Error during JSON sentry prediction unmarshaling\nError detail: %v\n", err)
		return err
	}
	return nil
}

func (sp *SentryPredictions) ToJSON(w io.Writer) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	d := json.NewEncoder(w)
	return d.Encode(sp.d)
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
func (sp *SentryPrediction) ToWSMessage(symbol string) *WsMsg {
	msg := &WsMsg{
		M: "sentry",
		D: struct {
			S string  "json:\"s\""
			T int64   "json:\"t\""
			V float64 "json:\"v\""
			E float64 "json:\"e\""
		}{
			S: symbol,
			T: util.ToEpoch(sp.Time),
			V: sp.Prediction,
		},
	}

	// msg.M = "sentry"
	// msg.D.T = util.ToEpoch(sp.Time)
	// msg.D.V = sp.Prediction
	return msg
}
