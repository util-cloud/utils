package main

import (
	"../../../journal"
)

func SetLog() {
	w := journal.NewConsoleWriter()
	w.SetColor(true)

	journal.Register(w)
	journal.SetLevel(journal.DEBUG)
	journal.SetLayout("2006-01-02 15:04:05")
}

func main() {
	SetLog()
	defer journal.Close()

	var name = "skoo"
	journal.Debug("journal by %s", name)
	journal.Info("journal by %s", name)
	journal.Warn("journal by %s", name)
	journal.Error("journal by %s", name)
	journal.Fatal("journal by %s", name)
}
