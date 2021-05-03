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
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("%v client is having some issue.\n", client)
				continue
			}
		}
	}
}

func main() {
	// Create a log file to log server activities
	logfile, _ := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer logfile.Close()
	log.SetOutput(logfile)

	// Initialize a channel dedicated to websocket messages which will be sent to clients
	msgQ := make(chan *data.WsMsg, 100)

	// Initilize 2 channels to communicate file update signals
	checkerChannel := make(chan struct{})
	historicalSentryChannel := make(chan struct{})
	multisaChannel := make(chan struct{})

	// Initialize current sentry data upon server start
	currentsentries := &data.Sentries{}
	currentsentries.Update()

	// Initialize future sentry data upon server start
	futuresentries := &data.SentryPredictions{}
	futuresentries.Update()

	// Initialize bot trade record data upon server start
	botTradeRecords := &data.BotTradeRecords{}
	botTradeRecords.Update()

	// Initialize a new Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set up routes for the router
	r.Get("/ws", wsHandler)
	r.Get("/history", setJsonHeaders(handler.MakeHistoryHandler(currentsentries)))
	r.Get("/sentry", setJsonHeaders(handler.MakePredictionHandler(futuresentries)))
	r.Get("/btr", setJsonHeaders(handler.MakeBotTradeRecordHandler(botTradeRecords)))
	handler.FileServer(r, "/", config.FrontendDir)

	// Initialize a custom HTTP server
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go writer(msgQ)

	go func() {
		for {
			select {
			case <-checkerChannel:
				futuresentries.Update()
				log.Printf("New data from checker file: %+v\n", futuresentries.Get())
				if sp, err := futuresentries.GetClosestFutureSentry(); err != nil {
					log.Println(err)
				} else {
					msgQ <- sp.ToWSMessage()
				}
			case <-historicalSentryChannel:
				currentsentries.Update()
				log.Printf("New current sentry value: %+v", currentsentries.GetCurrentSentry())
			case <-multisaChannel:
				botTradeRecords.Update()
				log.Printf("New trade was added to trade history file.")
			}
		}
	}()

	go util.WatchFile(config.SentryPredictionFile, checkerChannel, 6)
	go util.WatchFile(config.HistorySentryFile, historicalSentryChannel, 6)
	go util.WatchFile(config.MultiSaTradeRecords, multisaChannel, 6)

	go func() {
		log.Println("Server v0.5.9 listens on port", s.Addr)
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

// Not used for now
func handlepanic(fn func() error) {
	if r := recover(); r != nil {
		log.Println("RECOVER", r)
	}
}
