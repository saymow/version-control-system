package handlers

import (
	"os"
)

func TrackFiles(paths []string) {
	dir, err := os.Getwd()
	check(err)

	repository := getRepository(dir)

	for _, path := range paths {
		repository.IndexFile(path)
	}

	repository.SaveIndex()
}
