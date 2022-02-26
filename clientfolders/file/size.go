package file

import (
	"fmt"
	"os"
)

func Size(source string) (int64, error) {
	sourceInfo, err := os.Lstat(source)
	if err != nil {
		return 0, fmt.Errorf("could not stat source file: %w", err)
	}

	var size int64
	if sourceInfo.IsDir() {
		subDirSize, err := SizeDir(source)
		if err != nil {
			return 0, fmt.Errorf("could not get directory size: %w", err)
		}
		size = subDirSize
	} else {
		subFileSize, err := SizeFile(source)
		if err != nil {
			return 0, fmt.Errorf("could not get file size: %w", err)
		}

		size = subFileSize
	}

	return size, nil
}

func SizeDir(source string) (int64, error) {
	subfilesChan := make(chan []interface{})
	subErrChan := make(chan error)
	Find(source, []string{".*"}, func(path string) (interface{}, error) {
		size, err := Size(path)
		if err != nil {
			return nil, err
		}

		return size, nil
	}, subfilesChan, subErrChan)
	select {
	case values := <-subfilesChan:
		var size int64 = 0
		for _, value := range values {
			size += value.(int64)
		}
		return size, nil
	case err := <-subErrChan:
		return 0, fmt.Errorf("could not get directory size: %w", err)
	}
}

func SizeFile(source string) (int64, error) {
	sourceInfo, err := os.Lstat(source)

	if err != nil {
		return 0, fmt.Errorf("could not get file size: %w", err)
	}

	return sourceInfo.Size(), nil
}
