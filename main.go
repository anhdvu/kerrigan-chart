package main

import (
	"context"
	"kerrigan-chart/config"
	"kerrigan-chart/data"
	"kerrigan-chart/handler"
	"kerrigan-chart/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/gorilla/websocket"
)

// Handling websocket connections
var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// register client to a map of websocket clients
	clients[ws] = true
}

// func getHistory(w http.ResponseWriter, r *http.Request) {
// 	var records data.SentryJsons
// 	records, _ = data.GetSentryRecords()
// 	records.ToJSON(w)
// }

func Ping(mq chan<- *data.WsResponse) {
	for {
		pingMsg := &data.WsResponse{}
		pingMsg.M = "ping"
		pingMsg.D.T = time.Now().Unix()
		pingMsg.D.V = 1.111111
		mq <- pingMsg
		time.Sleep(4 * time.Minute)
	}
}

func writer(mq <-chan *data.WsResponse) {
	for {
		msg := <-mq
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}

func main() {
	msgQ := make(chan *data.WsResponse, 10)
	checkerChannel := make(chan struct{})
	historicalSentryChannel := make(chan struct{})
	sentryHistoryChannel := make(chan data.SentryJsons)
	latestSentryChannel := make(chan float64)
	sentryHistory, latestSentry := data.GetSentryRecords()
	mu := &sync.Mutex{}

	go data.UpdateSentryHistory(&sentryHistory, sentryHistoryChannel, mu)

	wsConn := data.NewKlineWebSocket()
	go func(ls *float64, lsc chan float64, wsc *websocket.Conn) {
		for {
			select {
			case *ls = <-lsc:
				log.Printf("new prediction updated - %v", latestSentry)
			default:
			}
			_, msg, err := wsc.ReadMessage()
			if err != nil {
				log.Println(err)
			}

			price := data.GetCurrentPrice(msg)
			dyde := &data.WsResponse{}
			dyde.M = "dyde"
			dyde.D.T = time.Now().Unix()
			dyde.D.V = latestSentry - price
			msgQ <- dyde
		}
	}(&latestSentry, latestSentryChannel, wsConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handler.FileServer(r)
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Get("/ws", wsHandler)
	r.Get("/history", setJsonHeaders(makeSentryHistoryHandler(&sentryHistory)))
	// r.Get("/history", setJsonHeaders(getHistory))

	go func() {
		for {
			select {
			case <-checkerChannel:
				log.Printf("%v has been updated!\n", config.SentryPredictionFile)
				defer handlepanic()
				msg, err := data.GetSentryPrediction()
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("%+v\n", msg)
					msgQ <- msg
				}
			case <-historicalSentryChannel:
				log.Printf("%v has been updated!\n", config.HistorySentryFile)
				sh, ls := data.GetSentryRecords()
				sentryHistoryChannel <- sh
				log.Println("New sentry history data sent through!")
				latestSentryChannel <- ls
				log.Println("New current sentry data sent through!")
			}
		}
	}()

	go util.WatchFile(config.SentryPredictionFile, checkerChannel, 30)
	go util.WatchFile(config.HistorySentryFile, historicalSentryChannel, 4)

	go Ping(msgQ)
	go writer(msgQ)
	go func() {
		log.Println("Server v0.4.5.2 listens on port", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cancel()
	s.Shutdown(ctx)

}

func setJsonHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fn(w, r)
	}
}

func makeSentryHistoryHandler(records *data.SentryJsons) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records.ToJSON(w)
	}
}

func handlepanic() {
	if a := recover(); a != nil {
		log.Println("RECOVER", a)
	}
}
