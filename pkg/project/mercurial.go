package project

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/yookoala/realpath"
)

// Mercurial contains mercurial data.
type Mercurial struct {
	// Filepath conaints the entity path.
	Filepath string
}

// Detect gets information about the mercurial project for a given file.
func (m Mercurial) Detect() (Result, bool, error) {
	fp, err := realpath.Realpath(m.Filepath)
	if err != nil {
		return Result{}, false,
			Err(fmt.Sprintf("failed to get the real path: %s", err))
	}

	// Take only the directory
	if fileExists(fp) {
		fp = path.Dir(fp)
	}

	// Find for .hg folder
	hgDirectory, ok := findHgConfigDir(fp)

	if ok {
		project := path.Base(path.Join(hgDirectory, ".."))

		branch, err := findHgBranch(hgDirectory)
		if err != nil {
			jww.ERROR.Printf(
				"error finding for branch name from %q: %s",
				hgDirectory,
				err,
			)
		}

		return Result{
			Project: project,
			Branch:  branch,
		}, true, nil
	}

	return Result{}, false, nil
}

func findHgConfigDir(fp string) (string, bool) {
	p := path.Join(fp, ".hg")
	if fileExists(p) {
		return p, true
	}

	dir := filepath.Clean(path.Join(fp, ".."))
	if dir == "/" {
		return "", false
	}

	return findHgConfigDir(dir)
}

func findHgBranch(fp string) (string, error) {
	p := path.Join(fp, "branch")
	if !fileExists(p) {
		return "default", nil
	}

	lines, err := readFile(p)
	if err != nil {
		return "", Err(fmt.Sprintf("failed while opening file %q: %s", fp, err))
	}

	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "default", nil
}

// String returns its name.
func (m Mercurial) String() string {
	return "hg-detector"
}
