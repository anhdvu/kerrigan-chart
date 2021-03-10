package handler

import (
	"kerrigan-chart/data"
	"net/http"
)

func MakeBotTradeRecordHandler(btr *data.BotTradeRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		btr.ToJSON(w)
	}
}
