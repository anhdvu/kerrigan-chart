package handler

import (
	"kerrigan-chart/data"
	"net/http"
)

func MakeHistoryHandler(ss *data.Sentries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ss.ToJSON(w)
	}
}
