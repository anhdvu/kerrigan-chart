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

var (
	wsClientsBTC = make(map[*websocket.Conn]bool)
	wsClientsETH = make(map[*websocket.Conn]bool)
	wsClientsBNB = make(map[*websocket.Conn]bool)
	btc          = "btc"
	eth          = "eth"
	bnb          = "bnb"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// register client to a map of websocket clients
	clients[ws] = true
}

func handleWSConn(clients map[*websocket.Conn]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		// register client to a map of websocket clients
		clients[ws] = true
	}
}

func writer(mq <-chan *data.WsMsg) {
	for {
		msg := <-mq
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("%v client is having some issue.\n", client)
				client.Close()
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
	btcPredictionChannel := make(chan struct{})
	btcSentryHistoryChannel := make(chan struct{})
	ethPredictionChannel := make(chan struct{})
	ethSentryHistoryChannel := make(chan struct{})
	bnbPredictionChannel := make(chan struct{})
	bnbSentryHistoryChannel := make(chan struct{})
	multisaChannel := make(chan struct{})

	// Initialize current sentry data upon server start
	btcCurrentSentries := &data.Sentries{}
	btcCurrentSentries.Update(config.HistorySentryFile_BTC)
	ethCurrentSentries := &data.Sentries{}
	ethCurrentSentries.Update(config.HistorySentryFile_ETH)
	bnbCurrentSentries := &data.Sentries{}
	bnbCurrentSentries.Update(config.HistorySentryFile_BNB)

	// Initialize future sentry data upon server start
	btcFutureSentries := &data.SentryPredictions{}
	btcFutureSentries.Update(config.SentryPredictionFile_BTC)
	ethFutureSentries := &data.SentryPredictions{}
	ethFutureSentries.Update(config.SentryPredictionFile_ETH)
	bnbFutureSentries := &data.SentryPredictions{}
	bnbFutureSentries.Update(config.SentryPredictionFile_BNB)

	// Initialize bot trade record data upon server start
	botTradeRecords := &data.BotTradeRecords{}
	botTradeRecords.Update()

	// Initialize a new Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set up routes for the router
	r.Get("/ws", wsHandler)
	r.Get("/ws_btc", handleWSConn(wsClientsBTC))
	r.Get("/ws_eth", handleWSConn(wsClientsETH))
	r.Get("/ws_bnb", handleWSConn(wsClientsBNB))
	r.Get("/history_btc", setJSONHeaders(handler.MakeHistoryHandler(btcCurrentSentries)))
	r.Get("/sentry_btc", setJSONHeaders(handler.MakePredictionHandler(btcFutureSentries)))
	r.Get("/history_eth", setJSONHeaders(handler.MakeHistoryHandler(ethCurrentSentries)))
	r.Get("/sentry_eth", setJSONHeaders(handler.MakePredictionHandler(ethFutureSentries)))
	r.Get("/history_bnb", setJSONHeaders(handler.MakeHistoryHandler(bnbCurrentSentries)))
	r.Get("/sentry_bnb", setJSONHeaders(handler.MakePredictionHandler(bnbFutureSentries)))
	r.Get("/btr", setJSONHeaders(handler.MakeBotTradeRecordHandler(botTradeRecords)))
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
			case <-btcPredictionChannel:
				handleCheckerChannel(btcFutureSentries, config.SentryPredictionFile_BTC, msgQ, btc)
			case <-ethPredictionChannel:
				handleCheckerChannel(ethFutureSentries, config.SentryPredictionFile_ETH, msgQ, eth)
			case <-bnbPredictionChannel:
				handleCheckerChannel(bnbFutureSentries, config.SentryPredictionFile_BNB, msgQ, bnb)
			case <-btcSentryHistoryChannel:
				handleHistoryChannel(btcCurrentSentries, config.HistorySentryFile_BTC)
			case <-ethSentryHistoryChannel:
				handleHistoryChannel(btcCurrentSentries, config.HistorySentryFile_ETH)
			case <-bnbSentryHistoryChannel:
				handleHistoryChannel(btcCurrentSentries, config.HistorySentryFile_BNB)
			case <-multisaChannel:
				botTradeRecords.Update()
				log.Println("New trade was added to trade history file.")
			}
		}
	}()

	go util.WatchFile(config.SentryPredictionFile_BTC, btcPredictionChannel, 6)
	go util.WatchFile(config.HistorySentryFile_BTC, btcSentryHistoryChannel, 6)
	go util.WatchFile(config.SentryPredictionFile_ETH, ethPredictionChannel, 6)
	go util.WatchFile(config.HistorySentryFile_ETH, ethSentryHistoryChannel, 6)
	go util.WatchFile(config.SentryPredictionFile_BNB, bnbPredictionChannel, 6)
	go util.WatchFile(config.HistorySentryFile_BNB, bnbSentryHistoryChannel, 6)
	go util.WatchFile(config.MultiSaTradeRecords, multisaChannel, 6)

	go func() {
		log.Println("Server v0.6 listens on port", s.Addr)
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
	defer cancel()
	s.Shutdown(ctx)
}

func setJSONHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next(w, r)
	}
}

func handleCheckerChannel(sps *data.SentryPredictions, file string, mq chan *data.WsMsg, symbol string) {
	sps.Update(file)
	log.Printf("New data from checker file: %+v\n", sps.Get())
	sp, err := sps.GetClosestFutureSentry()
	if err != nil {
		log.Println(err)
		return
	}
	mq <- sp.ToWSMessage(symbol)
}

func handleHistoryChannel(cs *data.Sentries, file string) {
	cs.Update(file)
	log.Printf("New current sentry value: %+v", cs.GetCurrentSentry())
}
