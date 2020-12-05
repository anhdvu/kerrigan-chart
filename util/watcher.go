package util

import (
	"log"
	"os"
	"time"
)

// WatchFile function used to watch for file change of a given file
func WatchFile(path string, c chan bool, t int) error {
	initialStat, err := os.Stat(path)
	if err != nil {
		log.Panic(err)
	}
	tempStat := initialStat.ModTime()
	for {
		stat, err := os.Stat(path)
		if err != nil {
			log.Panic(err)
			return err
		}
		if newStat := stat.ModTime(); newStat != tempStat {
			c <- true
			tempStat = newStat
		} else {
			c <- false
		}
		time.Sleep(time.Duration(t) * time.Second)
	}
}
