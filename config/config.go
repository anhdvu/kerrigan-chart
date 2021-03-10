package config

import (
	"os"
	"path/filepath"
)

var (
	cfgDir, _ = os.Getwd()
	// root                 = filepath.Dir(cfgDir)
	KerriganDir          = filepath.Dir(cfgDir)
	FrontendDir          = "frontend"
	HistorySentryFile    = filepath.Join(KerriganDir, "historical_delta.txt")
	SentryPredictionFile = filepath.Join(KerriganDir, "checker.txt")
	MultiSaTradeRecords  = filepath.Join(KerriganDir, "endpoints", "multisa", "trade_history.txt")
)
