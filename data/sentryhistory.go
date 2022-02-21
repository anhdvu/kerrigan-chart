package data

import (
	"encoding/json"
	"io"
	"kerrigan-chart/util"
	"log"
	"os"
	"sync"
)

type sentryRecord struct {
	Time       string
	Prediction float64 `json:"pred_price"`
	Actual     float64 `json:"actual_price"`
}

type sentry struct {
	Time   int64   `json:"time"`
	Value  float64 `json:"value"`
	Actual float64 `json:"actual"`
	Delta  float64 `json:"delta"`
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

func (ss *Sentries) Update(f string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	raw, err := os.ReadFile(f)
	if err != nil {
		log.Panicf("PANIC: Error reading file %v\nError detail: %v\n", f, err)
		return err
	}

	data := make([]sentryRecord, 0)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Printf("ERROR: Error during JSON sentry history unmarshaling\nError detail: %v\n", err)
		log.Println(string(raw))
		return err
	}

	ss.d = make([]sentry, len(data))
	for i, e := range data {
		t, _ := util.ToEpoch(e.Time)
		ss.d[i].Time = t
		ss.d[i].Value = e.Prediction
		ss.d[i].Actual = e.Actual
		ss.d[i].Delta = e.Prediction - e.Actual
	}
	return nil
}

func (ss *Sentries) ToJSON(w io.Writer) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	d := json.NewEncoder(w)
	return d.Encode(ss.d)

}

func (ss *Sentries) GetCurrentSentry() *sentry {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return &ss.d[len(ss.d)-1]
}

func (s *sentry) CalculateDynamicDelta(price float64) (dyde float64, safe float64) {
	dyde = s.Value - price
	safe = (dyde + 900) / 1800
	return dyde, safe
}
