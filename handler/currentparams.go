package handler

import (
	"encoding/json"
	"io/ioutil"
	"kerrigan-chart/config"
	"log"
	"net/http"
)

type Currentparams struct {
	Initial           float64 `json:"initial"`
	Thresh            float64 `json:"thresh"`
	Margin            float64 `json:"margin"`
	LastBuyPoint      float64 `json:"last_buying_point"`
	Stock             float64 `json:"stock"`
	Balance           float64 `json:"balance"`
	LastBuy           float64 `json:"last_buy"`
	LastSell          float64 `json:"last_sell"`
	LastPrediction    float64 `json:"last_pred"`
	Total             float64 `json:"total"`
	Decision          string  `json:"decision_making"`
	DecisionPoint     float64 `json:"decision_making_point"`
	TriggerPrediction float64 `json:"trig_pred"`
}

func GetCurrentParams(w http.ResponseWriter, r *http.Request) {
	if v := r.Header.Get("kerrigan"); v == "abcxyz" {
		raw, err := ioutil.ReadFile(config.CurrentParamsFile)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(raw)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func SetCurrentParams(w http.ResponseWriter, r *http.Request) {
	if v := r.Header.Get("kerrigan"); v == "abcxyz" {
		bodyjson := &Currentparams{}
		if err := json.NewDecoder(r.Body).Decode(bodyjson); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		log.Printf("%+v\n", bodyjson)
		response := &struct {
			Message string `json:"message"`
		}{"Current params have been updated!"}
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// raw, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// err = ioutil.WriteFile(config.CurrentParamsFile, raw, 0644)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotAcceptable)
	// 	return
	// }
}
