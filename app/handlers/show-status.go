package handlers

import (
	"fmt"
	"os"
	"saymow/version-manager/app/pkg/errors"
	"saymow/version-manager/app/repositories"
)

func printStatus(status *repositories.Status) {

	fmt.Println("Tracked changes:")
	for _, conflictedFile := range status.Staged.ConflictedFilesPaths {
		fmt.Printf("\t- %s (conflicted)\t%s\r\n", conflictedFile.Filepath, conflictedFile.Message)
	}
	for _, path := range status.Staged.CreatedFilesPaths {
		fmt.Printf("\t- %s (created)\r\n", path)
	}
	for _, path := range status.Staged.ModifiedFilePaths {
		fmt.Printf("\t- %s (modified)\r\n", path)
	}
	for _, path := range status.Staged.RemovedFilePaths {
		fmt.Printf("\t- %s (removed)\r\n", path)
	}

	fmt.Println("Untracked changes:")
	for _, path := range status.WorkingDir.UntrackedFilePaths {
		fmt.Printf("\t- %s (created)\r\n", path)
	}
	for _, path := range status.WorkingDir.ModifiedFilePaths {
		fmt.Printf("\t- %s (modified)\r\n", path)
	}
	for _, path := range status.WorkingDir.RemovedFilePaths {
		fmt.Printf("\t- %s (removed)\r\n", path)
	}
}

func ShowStatus() {
	dir, err := os.Getwd()
	errors.Check(err)

	repository := repositories.GetRepository(dir)
	status := repository.GetStatus()
	printStatus(status)
}
