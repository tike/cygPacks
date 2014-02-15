package main

import (
	"bytes"
	"os"
	"path/filepath"
)

type Repo struct {
	release      string
	arch         string
	setupVersion string

	packs []*Package
}

func (r *Repo) ParseSetupIni(raw []byte) int {
	packages := bytes.Split(raw, []byte("\n@ "))

	for _, pack := range packages[1:] {
		this := NewPackage()
		this.Parse(pack)
		r.packs = append(r.packs, this)
	}

	return len(r.packs)
}

func (r *Repo) Belongs(path string, info os.FileInfo, inErr error) (err error) {
	if inErr != nil {
		return inErr
	}

	if info.IsDir() {
		return
	}

	if filepath.Base(path) == "setup.ini" {
		return
	}

	var belongs bool
	for _, pack := range r.packs {
		if pack.ThatsMyFile(filepath.ToSlash(path), info) {
			belongs = true
			logger.Printf(3, "will keep:   %s", path)
			break
		}
	}

	if !belongs {
		logger.Printf(3, "will delete: %s", path)
	}
	return nil
}

func (r *Repo) CleanCache() error {
	return filepath.Walk("./", r.Belongs)
}
