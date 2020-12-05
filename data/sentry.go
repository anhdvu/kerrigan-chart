package data

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"kerrigan-chart/config"
	"kerrigan-chart/util"
	"log"
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
type SentryJsonSlice []sentryJson

type WsResponse struct {
	M string `json:"m"`
	D struct {
		T int64   `json:"t"`
		V float64 `json:"v"`
	} `json:"d"`
}

func GetSentryRecords() SentryJsonSlice {
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
	return result
}

func (sjs *SentryJsonSlice) ToJSON(w io.Writer) {
	d := json.NewEncoder(w)
	err := d.Encode(sjs)
	if err != nil {
		log.Panic(err)
	}
}

func GetSentryPrediction() WsResponse {
	data := make([]sentryPred, 0)
	result := WsResponse{}
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
			} else {
				result.M = "sentry"
				result.D.T = util.ToEpoch(e.Time)
				result.D.V = e.Prediction
				break
			}
		}
	} else {
		log.Println("checker.txt is empty for now.")
	}

	return result
}
