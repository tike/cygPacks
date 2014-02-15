package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Repo struct {
	release      string
	arch         string
	setupVersion string

	packs []*Package
}

func trimm(in []byte) string {
	var openQuote bool
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case 34:
			if openQuote {
				openQuote = false
			} else {
				openQuote = true
			}
		case 10:
			if openQuote {
				in[i] = 32
			}
		}
	}

	if in[len(in)-1] == 10 {
		in = in[:len(in)-1]
	}

	return string(in)
}

func (r *Repo) ParseSetupIni(raw []byte) (num int, err error) {
	sRaw := trimm(raw)
	packages := strings.Split(sRaw, "\n\n@ ")

	for i, pack := range packages[1:] {
		logger.Printf(4, "============================== %5d / %5d ================================", i, len(packages))
		this := NewPackage()
		if err = this.Parse(pack); err != nil {
			logger.Println(1, "Parsing: %s: %s", this.Name, err)
			return i, err
		}
		logger.Printf(4, "%#v", this)
		r.packs = append(r.packs, this)

	}

	return len(r.packs), err
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
			logger.Printf(4, "will keep:   %s", path)
			break
		}
	}

	if !belongs {
		logger.Printf(2, "deleting: %s", path)
		if !dryrun {
			if err = os.Remove(path); err != nil {
				logger.Printf(1, "Couldn't delete: %s because of %s", path, err)
			}
		}
	}
	return
}

func (r *Repo) CleanCache() error {
	return filepath.Walk("./", r.Belongs)
}
