package chooseui

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func SelectFileFromPath(dir string, _rec bool) (string, error) {

	dir = filepath.Clean(dir)

	files := make([]string, 0)
	//
	// Find files
	//
	var err error
	if !_rec {
		var _dl []fs.DirEntry
		_dl, err = os.ReadDir(dir)
		for _, d := range _dl {
			// We'll add anything that isn't a directory
			if !d.IsDir() {
				if !strings.Contains(d.Name(), "/.") && !strings.HasPrefix(d.Name(), ".") {
					files = append(files, filepath.Join(dir, d.Name()))
				}
			}
		}
	} else {
		err = filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {

				// Null info?  That probably means that the
				// destination we're trying to walk doesn't exist.
				if info == nil {
					return nil
				}

				// We'll add anything that isn't a directory
				if !info.IsDir() {
					if !strings.Contains(path, "/.") && !strings.HasPrefix(path, ".") {
						files = append(files, path)
					}
				}
				return nil
			})
	}

	if err != nil {
		return "", errors.New(fmt.Sprintf("error walking %s: %s", dir, err.Error()))
	}

	if len(files) < 1 {
		return "", errors.New(fmt.Sprintf("Failed to find any files beneath %s", dir))
	}

	//
	// Launch the UI
	//
	chooser := New(files)
	choice := chooser.Choose()

	return choice, nil
}
