package main

import (
	"../../../journal"
	"time"
)

func SetLog() {
	//w := journal.NewFileWriter()
	v := journal.NewFileVoyager()
	/*
	   %Y  year    (eg: 2014)
	   %M  month   (eg: 07)
	   %D  day     (eg: 05)
	   %H  hour    (eg: 18)
	   %m  minute  (eg: 29)

	   notice: No second's variable
	*/
	v.SetPathPattern("/tmp/logs/error%Y-%M-%D-%H:%m:%s.log")

	journal.Register(v)
	journal.SetLevel(journal.ERROR)
}

func main() {
	SetLog()
	defer journal.Close()

	var name = "linkerist"

	for {
		journal.Debug("journal by %s", name)
		journal.Info("journal by %s", name)
		journal.Warn("journal by %s", name)
		journal.Error("journal by %s", name)
		journal.Fatal("journal by %s", name)

		time.Sleep(time.Second * 1)
	}
}
