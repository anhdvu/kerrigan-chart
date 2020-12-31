package handler

import (
	"kerrigan-chart/data"
	"net/http"
)

func MakePredictionHandler(sp *data.SentryPredictions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sp.ToJSON(w)
	}
}
