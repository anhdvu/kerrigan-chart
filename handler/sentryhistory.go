package handler

import (
	"kerrigan-chart/data"
	"net/http"
)

func MakeHistoryHandler(ss *data.Sentries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ss.ToJSON(w)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
