package main

import (
	journal "../../../journal"
)

func SetLog() {
	w1 := journal.NewFileWriter()
	w1.SetPathPattern("/tmp/logs/error%Y%M%D%H.log")

	w2 := journal.NewConsoleWriter()

	journal.Register(w1)
	journal.Register(w2)
	journal.SetLevel(journal.ERROR)
}

func main() {
	SetLog()
	defer journal.Close()

	var name = "linkerist"

	journal.Debug("journal by %s", name)
	journal.Info("journal by %s", name)
	journal.Warn("journal by %s", name)
	journal.Error("journal by %s", name)
	journal.Fatal("journal by %s", name)
}
