package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadSetupIni(path string) (seBytes []byte, err error) {
	var seFile *os.File

	if seFile, err = os.Open(path); err != nil {
		logger.Printf(1, "Couldn't open %s: %s", path, err)
		return
	}
	defer seFile.Close()

	if seBytes, err = ioutil.ReadAll(seFile); err != nil {
		logger.Printf(1, "failed to read %s: %s", path, err)
		return
	}
	logger.Printf(3, "Read %s (%d kb)", path, len(seBytes)/1024)
	return
}

func main() {
	os.Chdir(dir)
	logger.Printf(3, "Changed to: %s", dir)

	setupInis, err := filepath.Glob("./*/setup.ini")
	if err != nil {
		logger.Printf(0, "Error finding setup.ini: %s", err)
		os.Exit(1)
	}
	if setupInis == nil {
		logger.Printf(0, "No setup.ini found!")
		os.Exit(1)
	}

	repo := &Repo{}

	logger.Printf(3, "Found %d setup.ini files: %s", len(setupInis), setupInis)
	for _, ini := range setupInis {
		var rawBytes []byte
		if rawBytes, err = ReadSetupIni(ini); err != nil {
			continue
		}

		num, err := repo.ParseSetupIni(rawBytes)
		if err != nil {
			logger.Printf(3, "Error parsing package %d: %s", num, err)
		}
		logger.Printf(3, "Found %d packages.", num)

		if err := repo.CleanCache(); err != nil {
			logger.Printf(1, "Cleaning failed: %s", err)
		}

	}

}
