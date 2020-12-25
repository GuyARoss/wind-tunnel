package utilities

import (
	"os"
	"path/filepath"
)

func findFiles(rootDir string, fileType string) []string {
	files := make([]string, 0)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == fileType {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}
