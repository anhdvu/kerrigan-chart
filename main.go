package main

import (
	"kerrigan-chart/config"
	"kerrigan-chart/data"
	"kerrigan-chart/handler"
	"kerrigan-chart/util"
	"log"
	"net/http"
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

type SentryPrediction struct {
	Time       string
	Prediction float64 `json:"predict"`
}

type WsResponse struct {
	M string `json:"m"`
	D struct {
		T int64   `json:"t"`
		V float64 `json:"v"`
	} `json:"d"`
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
	var records data.SentryJsonSlice
	records = data.GetSentryRecords()
	records.ToJSON(w)
}

func getMarkers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Placeholder for markers data."))
}

func readSentryPrediction(c chan bool, mq chan data.WsResponse) {
	for {
		select {
		case <-c:
			log.Printf("%v has been modified.", config.SentryPredictionFile)
			mq <- data.GetSentryPrediction()
		default:
			log.Println("No file modification detected for the last 30 seconds.")
			time.Sleep(30 * time.Second)
		}
	}
}

func makeMarkerFile(c chan bool) {
	if <-c {
		log.Println("There seems to be a new trade.")

	} else {
		log.Println("records.txt has not changed for the last hour.")
	}
}

func Ping(mq chan<- data.WsResponse) {
	for {
		pingMsg := data.WsResponse{}
		pingMsg.M = "ping"
		pingMsg.D.T = time.Now().Unix()
		pingMsg.D.V = 1.111111
		mq <- pingMsg
		time.Sleep(4 * time.Minute)
	}
}

func writer(mq <-chan data.WsResponse) {
	for {
		msg := <-mq
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}

func main() {
	msgQ := make(chan data.WsResponse, 10)
	checkerChannel := make(chan bool, 1)
	// markerChannel := make(chan bool, 1)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handler.FileServer(r)

	r.Get("/history", setJsonHeaders(getHistory))
	r.Get("/markers", setJsonHeaders(getMarkers))

	r.Get("/ws", wsHandler)

	go util.WatchFile(config.SentryPredictionFile, checkerChannel, 30)

	// go util.WatchFile(recordF, markerChannel, 3600)
	go readSentryPrediction(checkerChannel, msgQ)
	go writer(msgQ)
	go Ping(msgQ)

	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Server v0.4.0 listens on port", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func setJsonHeaders(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		f(w, r)
	}
}
