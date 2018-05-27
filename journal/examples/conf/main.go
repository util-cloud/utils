package main

import (
	journal "../../../journal"
)

func main() {
	if err := journal.SetupLogWithConf("log.json"); err != nil {
		panic(err)
	}
	defer journal.Close()

	var name = "linkerist"
	journal.Debug("journal by %s", name)
	journal.Info("journal by %s", name)
	journal.Warn("journal by %s", name)
	journal.Error("journal by %s", name)
	journal.Fatal("journal by %s", name)
}
