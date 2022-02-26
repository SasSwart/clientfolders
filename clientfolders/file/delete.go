package file

import (
	"fmt"
	"os"

	"github.com/sasswart/clientfolders/clientfolders/find"
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
	subfilesChan := make(chan []string)
	subErrChan := make(chan error)
	find.Find(source, []string{".*"}, func(path string, errchan chan error) {
		err := Delete(path)
		if err != nil {
			errchan <- err
		}

		errchan <- nil
	}, subfilesChan, subErrChan)
	select {
	case <-subfilesChan:
		return nil
	case err := <-subErrChan:
		return fmt.Errorf("could not delete directory: %w", err)
	}
}

func DeleteFile(source string) error {
	err := os.Remove(source)
	if err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}
	return nil
}
