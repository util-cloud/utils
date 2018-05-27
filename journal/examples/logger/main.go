package main

import (
	//log "log4go"
	log4go "../../../journal"
	"time"
)

func SetLog() (logger * log4go.Journal) {
	w := log4go.NewFileVoyager()
	w.SetPathPattern("/tmp/logs/error%Y%M%D%H:%m:%s.log")

	logger = log4go.NewJournal()
	logger.Register(w)
	logger.SetLevel(log4go.WARNING)
	return
}

func main() {
	logger := SetLog()
	time.Sleep(5 * time.Second)
	logger1 := SetLog()
	defer logger1.Close()
	defer logger.Close()

	var name = "Linkerist"
	logger.Debug("log4go by %s", name)
	logger.Info("log4go by %s", name)
	logger1.Warn("log4go by %s", name)
	logger.Warn("log4go by %s", name)
	logger.Error("log4go by %s", name)
	logger.Fatal("log4go by %s", name)
}
