package data

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kerrigan-chart/config"
	"kerrigan-chart/util"
	"log"
	"sync"
	"time"
)

type sentryRecord struct {
	Time       string
	Prediction float64 `json:"pred_price"`
	Actual     float64 `json:"actual_price;omitempty"`
}

type sentryPred struct {
	Time       string
	Prediction float64 `json:"predict"`
}

type sentryJson struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}
type SentryJsons []sentryJson

type WsResponse struct {
	M string `json:"m"`
	D struct {
		T int64   `json:"t"`
		V float64 `json:"v"`
	} `json:"d"`
}

func GetSentryRecords() (SentryJsons, float64) {
	raw, err := ioutil.ReadFile(config.HistorySentryFile)
	if err != nil {
		log.Panic(err)
	}

	data := make([]sentryRecord, 0)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Panic(err)
	}

	result := make([]sentryJson, len(data))
	for i, e := range data {
		result[i].Time = util.ToEpoch(e.Time)
		result[i].Value = e.Prediction
	}
	return result, result[len(result)-1].Value
}

func (sjs *SentryJsons) ToJSON(w io.Writer) {
	d := json.NewEncoder(w)
	err := d.Encode(sjs)
	if err != nil {
		log.Panic(err)
	}
}

func GetSentryPrediction() (*WsResponse, error) {
	data := make([]sentryPred, 0)
	result := &WsResponse{}
	raw, err := ioutil.ReadFile(config.SentryPredictionFile)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Panic(err)
	}
	now := time.Now().Unix()

	if len(data) > 0 {
		for _, e := range data {
			if util.ToEpoch(e.Time) < now {
				log.Println("there was a re-prediction.")
			} else if util.ToEpoch(e.Time) > now {
				result.M = "sentry"
				result.D.T = util.ToEpoch(e.Time)
				result.D.V = e.Prediction
				break
			} else {
				return &WsResponse{}, fmt.Errorf("The file has no new data %v", WsResponse{})
			}
		}
	} else {
		log.Println("checker.txt is empty for now.")
	}
	return result, nil
}

func UpdateSentryHistory(sh *SentryJsons, shc chan SentryJsons, mu *sync.Mutex) {
	for {
		mu.Lock()
		*sh = <-shc
		mu.Unlock()
		log.Println((*sh)[len(*sh)-1])
	}
}
