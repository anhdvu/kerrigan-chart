package config

import (
	"os"
	"path/filepath"
)

var (
	appDir, _                = os.Getwd()
	MoonDir                  = filepath.Dir(appDir)
	FrontendDir              = "frontend"
	DataDir                  = filepath.Join(MoonDir, "sentry", "data")
	HistorySentryFile_BTC    = filepath.Join(DataDir, "history_btcusdt.json")
	SentryPredictionFile_BTC = filepath.Join(DataDir, "btcusdt.json")
	HistorySentryFile_ETH    = filepath.Join(DataDir, "history_ethusdt.json")
	SentryPredictionFile_ETH = filepath.Join(DataDir, "ethusdt.json")
	HistorySentryFile_BNB    = filepath.Join(DataDir, "history_bnbnusdt.json")
	SentryPredictionFile_BNB = filepath.Join(DataDir, "bnbusdt.json")
	MultiSaTradeRecords      = filepath.Join(MoonDir, "kerrigan", "endpoints", "multisa", "trade_history.txt")
)
