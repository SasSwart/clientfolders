package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sasswart/clientfolders/clientfolders/find"
)

func Copy(destination, source string) error {
	err := os.MkdirAll(filepath.Dir(destination), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create destination parent directories: %w", err)
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("could not stat source file: %w", err)
	}
	if sourceInfo.IsDir() {
		err := CopyDir(destination, source)
		if err != nil {
			return fmt.Errorf("could not copy directory: %w", err)
		}
	} else {
		err := CopyFile(destination, sourceFile)
		if err != nil {
			return fmt.Errorf("could not copy file: %w", err)
		}
	}

	return nil
}

func CopyDir(target, source string) error {
	_, err := find.Find(source, []string{".*"}, func(path string) error {
		destination := strings.ReplaceAll(path, source, target)
		err := Copy(destination, path)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not copy directory: %w", err)
	}
	return nil
}

func CopyFile(destination string, source *os.File) error {
	destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open destination file: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, source)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}
	return nil
}
