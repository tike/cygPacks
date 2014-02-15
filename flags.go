package main

import (
	"flag"
	"fmt"
	"os"
)

// saving options
var dir string
var dryrun bool

// general options
var verbosity int
var logFile string
var logVerbosity int

func init() {
	flag.StringVar(&dir, "d", "ftp%3a%2f%2fftp.gwdg.de%2fpub%2flinux%2fsources.redhat.com%2fcygwin%2f", "parent directory in which directory for this crawl is created")
	flag.IntVar(&verbosity, "v", 4, "verbosity level 1-5")
	flag.StringVar(&logFile, "l", "", "log file name")
	flag.IntVar(&logVerbosity, "lv", 4, "verbosity of logfile output")
	flag.BoolVar(&dryrun, "n", false, "do everything exept touching the filesystem.")
	flag.Parse()

	var err error
	if logger, err = setupLogger(); err != nil {
		fmt.Println("Error setting up logging: %s", err)
		os.Exit(127)
	}
}

// the Usage function...
func Usage() {
	fmt.Printf("%s [Options]", os.Args[0])
	flag.PrintDefaults()
}
