package project_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/wakatime/wakatime-cli/pkg/project"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubversion_Detect(t *testing.T) {
	skipIfBinaryNotFound(t)

	fp, tearDown := setupTestSvn(t)
	defer tearDown()

	s := project.Subversion{
		Filepath: path.Join(fp, "wakatime-cli/src/pkg/file.go"),
	}

	result, detected, err := s.Detect()
	require.NoError(t, err)

	assert.True(t, detected)
	assert.Equal(t, project.Result{
		Project: "wakatime-cli",
		Branch:  "trunk",
	}, result)
}

func TestSubversion_Detect_Branch(t *testing.T) {
	skipIfBinaryNotFound(t)

	fp, tearDown := setupTestSvnBranch(t)
	defer tearDown()

	s := project.Subversion{
		Filepath: path.Join(fp, "wakatime-cli/src/pkg/file.go"),
	}

	result, detected, err := s.Detect()
	require.NoError(t, err)

	assert.True(t, detected)
	assert.Equal(t, project.Result{
		Project: "wakatime-cli",
		Branch:  "billing",
	}, result)
}

func setupTestSvn(t *testing.T) (fp string, tearDown func()) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "wakatime-svn")
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(tmpDir, "wakatime-cli/src/pkg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	tmpFile, err := os.Create(path.Join(tmpDir, "wakatime-cli/src/pkg/file.go"))
	require.NoError(t, err)

	tmpFile.Close()

	copyDir(t, "testdata/svn", path.Join(tmpDir, "wakatime-cli/.svn"))

	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func setupTestSvnBranch(t *testing.T) (fp string, tearDown func()) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "wakatime-svn")
	require.NoError(t, err)

	err = os.MkdirAll(path.Join(tmpDir, "wakatime-cli/src/pkg"), os.FileMode(int(0700)))
	require.NoError(t, err)

	tmpFile, err := os.Create(path.Join(tmpDir, "wakatime-cli/src/pkg/file.go"))
	require.NoError(t, err)

	tmpFile.Close()

	copyDir(t, "testdata/svn_branch", path.Join(tmpDir, "wakatime-cli/.svn"))

	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func copyDir(t *testing.T, src string, dst string) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	require.NoError(t, err)

	if !si.IsDir() {
		return
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if err == nil {
		return
	}

	err = os.MkdirAll(dst, si.Mode())
	require.NoError(t, err)

	entries, err := ioutil.ReadDir(src)
	require.NoError(t, err)

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			copyDir(t, srcPath, dstPath)
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			copyFile(t, srcPath, dstPath)
		}
	}
}

func findSvnBinary() (string, bool) {
	locations := []string{
		"svn",
		"/usr/bin/svn",
		"/usr/local/bin/svn",
	}

	for _, loc := range locations {
		cmd := exec.Command(loc, "--version")

		err := cmd.Run()
		if err != nil {
			jww.ERROR.Printf("failed while calling %s --version: %s", loc, err)
			continue
		}

		return loc, true
	}

	return "", false
}

func skipIfBinaryNotFound(t *testing.T) {
	_, found := findSvnBinary()
	if !found {
		t.Skip("Skipping because the lack of svn binary in this machine.")
	}
}
