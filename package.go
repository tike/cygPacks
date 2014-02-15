package main

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	// "fmt"
)

type RepoFile struct {
	Name string
	size uint64
	md5  string
}

func (r *RepoFile) ThatsMe(path string, info os.FileInfo) bool {
	return r.Name == path
}

func (r *RepoFile) Parse(raw []byte) (err error) {
	parts := bytes.Split(raw, []byte(" "))
	if len(parts) < 3 {
		return errors.New("Malformed FileEntry.")
	}
	//logger.Printf(4, "Parts: _%s___%s___%s_", string(parts[0]), string(parts[1]), string(parts[2]))
	r.Name = string(parts[0])
	r.md5 = string(parts[2])
	r.size, err = strconv.ParseUint(string(parts[1]), 10, 64)
	return
}

type Package struct {
	Name     string
	SDesc    string
	LDesc    string
	Category string
	Requires []string
	Version  string
	Install  *RepoFile
	Source   *RepoFile
}

func NewPackage() *Package {
	return &Package{
		Requires: make([]string, 0, 5),
		Install:  new(RepoFile),
		Source:   new(RepoFile),
	}
}

func (p *Package) ThatsMyFile(path string, info os.FileInfo) (meins bool) {
	if strings.Contains(path, p.Name) {
		//logger.Printf(4, "_%s_ vs. _%s_", p.Install.Name, path)
		return p.Install.ThatsMe(path, info) || p.Source.ThatsMe(path, info)
	}
	return false
}

func (p *Package) String() string {
	return p.Name //+"_"+ p.Install.String()
}

func (this *Package) Parse(pack []byte) {
	lines := bytes.Split(pack, []byte("\n"))

	this.Name = string(lines[0])
	lines = lines[1:]

	var ldesc []byte
	for _, line := range lines {
		vals := bytes.Split(line, []byte(": "))

		if len(vals) <= 1 {
			if bytes.Compare(vals[0], []byte("[prev]")) == 0 {
				logger.Printf(4, "Skipping prev Version: %s", this)
				break
			}
			ldesc = append(ldesc, line...)
			continue
		}

		key := string(vals[0])
		switch key {
		case "sdesc":
			quoted := string(vals[1])
			this.SDesc = quoted[1 : len(quoted)-1]

		case "category":
			this.Category = string(vals[1])

		case "requires":
			reqs := bytes.Split(vals[1], []byte(" "))
			for _, req := range reqs {
				this.Requires = append(this.Requires, string(req))
			}

		case "version":
			this.Version = string(vals[1])

		case "install":
			if err := this.Install.Parse(vals[1]); err != nil {
				logger.Printf(2, "%s install: %s", this, err)
			}

		case "source":
			if err := this.Source.Parse(vals[1]); err != nil {
				logger.Printf(2, "%s source: %s", this, err)
			}

		}
	}

	this.LDesc = string(ldesc)
	return
}
