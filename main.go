package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kerrigan-chart/util"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	root, _           = os.Getwd()
	kDir              = filepath.Dir(root)
	fDir              = filepath.Join(root, "frontend")
	staticDir         = filepath.Join(fDir, "static")
	sentryHistoryF    = filepath.Join(kDir, "historical_delta.txt")
	sentryPredictionF = filepath.Join(kDir, "checker.txt")
)

type Sentry struct {
	Time       string
	Prediction float64 `json:"pred_price"`
	Actual     float64 `json:"actual_price;omitempty"`
}

type SentryPrediction struct {
	Time       string
	Prediction float64 `json:"predict"`
}

type SentryResponse struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type WsResponse struct {
	M string `json:"m"`
	D struct {
		T int64   `json:"t"`
		V float64 `json:"v"`
	} `json:"d"`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hey yo!")
	w.Write([]byte("Hey yo!"))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client to a map of websocket clients
	clients[ws] = true
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadFile(sentryHistoryF)
	if err != nil {
		log.Panic(err)
	}

	data := make([]Sentry, 0)
	err = json.Unmarshal(raw, &data)
	if err != nil {
		log.Panic(err)
	}

	response := make([]SentryResponse, len(data))
	for i, e := range data {
		response[i].Time = parseTime(e.Time)
		response[i].Value = e.Prediction
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panic(err)
	}
}

func parseTime(ts string) int64 {
	// Jan 2 15:04:05 2006 MST
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, ts)
	if err != nil {
		log.Panic(err)
	}
	return t.Unix()
}

func readSentryPrediction(path string, c chan bool, mq chan WsResponse) {
	for {
		if <-c {
			log.Println("There was a file change.")
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				log.Panic(err)
			}
			data := make([]SentryPrediction, 0)
			err = json.Unmarshal(raw, &data)
			if err != nil {
				log.Panic(err)
			}

			msg := WsResponse{}
			msg.M = "sentry"
			msg.D.T = parseTime(data[0].Time)
			msg.D.V = data[0].Prediction
			mq <- msg
			time.Sleep(5 * time.Second)
		} else {
			log.Println("There was NO file change for the past 60 seconds")
			time.Sleep(5 * time.Second)
		}
	}
}

func writer(msgQ chan WsResponse) {
	for {
		msg := <-msgQ
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}

func main() {
	msgQ := make(chan WsResponse, 10)
	fc := make(chan bool, 1)
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Handle("/", http.FileServer(http.Dir(fDir)))
	r.Get("/chart", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "kerrigan-chart.js"))
	})
	r.Get("/history", getHistory)
	r.Get("/ws", wsHandler)

	var interval time.Duration = 5
	go util.WatchFile(sentryPredictionF, fc, interval)
	go readSentryPrediction(sentryPredictionF, fc, msgQ)
	go writer(msgQ)

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(s.ListenAndServe())
}
