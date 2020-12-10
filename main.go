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

func writer(mq <-chan *data.WsMsg) {
	for {
		msg := <-mq
		for client := range clients {
			client.WriteJSON(msg)
		}
	}
}

func main() {
	// Initialize a channel dedicated to websocket messages which will be sent to clients
	msgQ := make(chan *data.WsMsg, 100)

	// Initilize 2 channels to communicate file update signals
	checkerChannel := make(chan struct{})
	historicalSentryChannel := make(chan struct{})

	// Initialize current sentry data upon server start
	currentsentries := &data.Sentries{}
	currentsentries.Update()

	// Initialize future sentry data upon server start
	// immediately send the upcoming sentry to message channel, which in turn will be sent to websocket clients
	futuresentries := &data.SentryPredictions{}
	futuresentries.Update()
	msgQ <- futuresentries.GetClosestFutureSentry().ToWSMessage()

	// Initialize a websocket client used to retrieve current BTC-USDT price
	wsConn := data.NewKlineWebSocket()
	go func(wsc *websocket.Conn, s *data.Sentries) {
		for {
			_, msg, err := wsc.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			price := data.GetCurrentPrice(msg)
			dyde := &data.WsMsg{}
			dyde.M = "dyde"
			dyde.D.T = time.Now().Unix()
			dyde.D.V = s.GetCurrentSentryValue() - price
			msgQ <- dyde
		}
	}(wsConn, currentsentries)

	// Initialize a new router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handler.FileServer(r)
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Set up routes for the router
	r.Get("/ws", wsHandler)
	r.Get("/history", setJsonHeaders(makeHistoryHandler(currentsentries)))
	go writer(msgQ)

	go func() {
		for {
			select {
			case <-checkerChannel:
				log.Printf("Worker received a signal that %v has been updated!\n", config.SentryPredictionFile)
				futuresentries.Update()
				log.Printf("%+v\n", futuresentries.Get())
				msgQ <- futuresentries.GetClosestFutureSentry().ToWSMessage()
			case <-historicalSentryChannel:
				log.Printf("Worker received a signal that %v has been updated!\n", config.HistorySentryFile)
				currentsentries.Update()
				log.Printf("%+v", currentsentries.GetCurrentSentryValue())
			}
		}
	}()

	go util.WatchFile(config.SentryPredictionFile, checkerChannel, 30)
	go util.WatchFile(config.HistorySentryFile, historicalSentryChannel, 6)

	go func() {
		log.Println("Server v0.5 listens on port", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

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

func makeHistoryHandler(s *data.Sentries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.ToJSON(w)
	}
}
