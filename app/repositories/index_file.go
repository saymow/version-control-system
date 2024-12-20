package repositories

import (
	"os"
	"saymow/version-manager/app/pkg/errors"
	"saymow/version-manager/app/repositories/directories"
	"slices"
)

func (repository *Repository) IndexFile(filepath string) error {
	if repository.isDetachedMode() {
		return &ValidationError{"cannot make changes in detached mode."}
	}

	filepath, err := repository.dir.AbsPath(filepath)
	if err != nil {
		return &ValidationError{err.Error()}
	}

	file, err := os.Open(filepath)
	errors.Check(err)
	defer errors.CheckFn(file.Close)

	object := repository.fs.WriteObject(filepath, file)
	stagedChangeIdx := repository.findStagedChangeIdx(filepath)
	savedObject := repository.findSavedFile(filepath)
	var ChangeType directories.ChangeType

	if savedObject != nil {
		ChangeType = directories.Modification
	} else {
		ChangeType = directories.Creation
	}

	if savedObject != nil && savedObject.ObjectName == object.ObjectName {
		// No changes at all

		if stagedChangeIdx != -1 {
			if repository.index[stagedChangeIdx].GetHash() != object.ObjectName {
				// Remove change file object
				repository.fs.RemoveObject(repository.index[stagedChangeIdx].GetHash())
			}

			// Undo index existing change
			repository.index = slices.Delete(repository.index, stagedChangeIdx, stagedChangeIdx+1)
		}
	} else if stagedChangeIdx != -1 {
		stagedChange := repository.index[stagedChangeIdx]

		if stagedChange.ChangeType == directories.Conflict &&
			stagedChange.GetHash() != object.ObjectName &&
			stagedChange.Conflict.IsObjectTemporary() {
			// Remove conflicted temp file object
			repository.fs.RemoveObject(stagedChange.GetHash())
		}

		if (stagedChange.ChangeType == directories.Creation || stagedChange.ChangeType == directories.Modification) &&
			stagedChange.GetHash() != object.ObjectName {
			// Remove change file object
			repository.fs.RemoveObject(stagedChange.GetHash())
		}

		// Undo index existing change
		repository.index = slices.Delete(repository.index, stagedChangeIdx, stagedChangeIdx+1)
		// Index change
		repository.index = append(repository.index, &directories.Change{ChangeType: ChangeType, File: object})

	} else {
		// Index change
		repository.index = append(repository.index, &directories.Change{ChangeType: ChangeType, File: object})
	}

	return nil
}
