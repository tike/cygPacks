package main

import (
	//"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const PREV = "[prev]"
const TEST = "[test]"

type RepoFile struct {
	Name string
	size uint64
	md5  string
}

func (r *RepoFile) ThatsMe(path string, info os.FileInfo) bool {
	return r.Name == path
}

func (r *RepoFile) Parse(raw string) (err error) {
	parts := strings.Split(raw, " ")
	if len(parts) < 3 {
		return errors.New("Malformed FileEntry.")
	}
	//logger.Printf(4, "Parts: _%s___%s___%s_", string(parts[0]), string(parts[1]), string(parts[2]))
	r.Name = parts[0]
	r.md5 = parts[2]
	r.size, err = strconv.ParseUint(parts[1], 10, 64)
	return
}

type Version struct {
	Version string
	Install *RepoFile
	Source  *RepoFile
}

func NewVersion() *Version {
	return &Version{
		Install: &RepoFile{},
		Source:  &RepoFile{},
	}
}

func (this *Version) Parse(lines []string) (err error) {
	for _, line := range lines {
		vals := strings.Split(line, ": ")
		if len(vals) < 2 {
			logger.Printf(4, "Version:")
			for i, line := range lines {
				logger.Printf(4, "%d/%d: %s", i, len(lines), line)
			}
			panic("this looks shit:")
			return
		}

		switch vals[0] {
		case "version":
			this.Version = vals[1]

		case "install":
			if err = this.Install.Parse(vals[1]); err != nil {
				logger.Printf(2, "%s install: %s", this, err)
				return
			}

		case "source":
			if err = this.Source.Parse(vals[1]); err != nil {
				logger.Printf(2, "%s source: %s", this, err)
				return
			}
		}
	}

	return
}

type Package struct {
	Name     string
	SDesc    string
	LDesc    string
	Category string
	Requires []string
	Message  string
	Versions []*Version
}

func NewPackage() *Package {
	return &Package{
		Requires: make([]string, 0, 5),
		Versions: make([]*Version, 0, 2),
	}
}

func (p *Package) ThatsMyFile(path string, info os.FileInfo) (meins bool) {
	if strings.Contains(path, p.Name) {
		//logger.Printf(4, "_%s_ vs. _%s_", p.Install.Name, path)
		for _, version := range p.Versions {
			if version.Install.ThatsMe(path, info) || version.Source.ThatsMe(path, info) {
				return true
			}
		}
	}
	return false
}

func (p *Package) String() string {
	return p.Name
}

func (this *Package) ParseHeader(lines []string) (rlines []string) {
	this.Name = lines[0]

	for _, line := range lines[1:] {
		vals := strings.Split(line, ": ")

		if len(vals) == 1 {
			if line != PREV && line != TEST {
				panic(fmt.Sprintf("Malformed line: _%s_", line))
			} else {
				rlines = append(rlines, line)
			}
			continue
		}

		switch vals[0] {
		case "sdesc":
			this.SDesc = vals[1][1 : len(vals[1])-1]

		case "ldesc":
			this.LDesc = vals[1][1:]

		case "category":
			this.Category = vals[1]

		case "requires":
			this.Requires = append(this.Requires, strings.Split(vals[1], " ")...)

		case "message":
			this.Message = vals[1][1 : len(vals[1])-1]
		default:
			rlines = append(rlines, line)
		}
	}

	return
}

func (this *Package) ParseMultiVer(lines []string) (err error) {
	versions := make([][]string, 0, 3)

	nextStart := 0
	for i, line := range lines {
		if line == PREV || line == TEST {
			if i > 0 {
				versions = append(versions, lines[nextStart:i-1])
			}
			nextStart = i + 1
		}
	}
	versions = append(versions, lines[nextStart:])

	for _, lines := range versions {
		v := NewVersion()
		if err = v.Parse(lines); err != nil {
			return
		}
		this.Versions = append(this.Versions, v)
	}
	return
}

func (this *Package) Parse(pack string) (err error) {
	lines := strings.Split(pack, "\n")
	logger.Printf(4, "Package %s", lines[0])

	lines = this.ParseHeader(lines)

	if len(lines) > 3 {
		return this.ParseMultiVer(lines)
	} else {
		v := NewVersion()
		if err = v.Parse(lines); err != nil {
			return
		}
		this.Versions = append(this.Versions, v)
	}

	return
}
