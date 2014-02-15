package main

import (
	l "log"
	"log5"
	"time"
)

// The logger
var logger log5.Log5

func setupLogger() (logger log5.Log5, err error) {
	logger = log5.NewLog5(2)
	if logger, err = logger.Add("", verbosity-1, l.Ltime|l.Lshortfile); err != nil {
		return
	}
	if logFile != "" {
		if logFile == "auto" {
			logFile = "cygPacks_" + time.Now().Format("2006-01-02__15-04-05")
		}
		if logger, err = logger.Add(logFile, logVerbosity-1, l.Ltime|l.Lmicroseconds|l.Lshortfile); err != nil {
			return
		}
		//logger[1].Logger.Printf("Will log %s-level to %s\n", log5.Levels[*verbose-1], *logFile)
	}
	//logger[0].Logger.Println("Will log", log5.Levels[*verbose-1], "level to screen")
	return
}
