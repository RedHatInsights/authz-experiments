package utils

import (
	"github.com/golang/glog"
	"os"
	"path/filepath"
	"strings"
)

// Walk walks the file tree rooted at root including symlinks, calling fn for
// each file in the tree.
func Walk(root string, fn filepath.WalkFunc) error {
	return walk(root, "", fn)
}

func walk(filename string, linkDirName string, fn filepath.WalkFunc) error {
	if linkDirName == "" {
		linkDirName = filename
	}

	symWalkFunc := func(path string, info os.FileInfo, err error) error {
		relName, relErr := filepath.Rel(filename, path)
		if relErr != nil {
			return relErr
		}

		path = filepath.Join(linkDirName, relName)
		if err != nil {
			return fn(path, info, err)
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, fileErr := filepath.EvalSymlinks(path)
			if fileErr != nil {
				return fileErr
			}

			info, err = os.Lstat(finalPath)
			if err != nil {
				return fn(path, info, err)
			}

			if info.IsDir() {
				return walk(finalPath, path, fn)
			}

			path = finalPath
		}

		if info != nil && !info.IsDir() {
			return fn(path, info, err)
		}

		return nil
	}

	return filepath.Walk(filename, symWalkFunc)
}

func GetProjectRootDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}
	dirs := strings.Split(workingDir, "/")
	var goModPath string
	var rootPath string
	for _, d := range dirs {
		rootPath = rootPath + "/" + d
		goModPath = rootPath + "/go.mod"
		goModFile, err := os.ReadFile(goModPath)
		if err != nil { // if the file doesn't exist, continue searching
			continue
		}
		// The project root directory is obtained based on the assumption that module name,
		// "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager", is contained in the 'go.mod' file.
		// Should the module name change in the code repo then it needs to be changed here too.
		if strings.Contains(string(goModFile), "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager") {
			break
		}
	}
	return rootPath
}
