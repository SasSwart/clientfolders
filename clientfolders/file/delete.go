package file

import (
	"fmt"
	"os"
)

func Delete(source string) error {

	sourceInfo, err := os.Lstat(source)
	if err != nil {
		return fmt.Errorf("could not stat source file: %w", err)
	}

	if sourceInfo.IsDir() {
		err := DeleteDir(source)
		if err != nil {
			return fmt.Errorf("could not delete directory: %w", err)
		}
	} else {
		err := DeleteFile(source)
		if err != nil {
			return fmt.Errorf("could not delete file: %w", err)
		}
	}

	return nil
}

func DeleteDir(source string) error {
	subErrChan := make(chan error)
	Find(source, []string{".*"}, func(path string) (interface{}, error) {
		err := Delete(path)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, nil, subErrChan)

	if err := <-subErrChan; err != nil {
		return fmt.Errorf("could not delete directory: %w", err)
	}
	return nil
}

func DeleteFile(source string) error {
	err := os.Remove(source)
	if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}
	return nil
}
