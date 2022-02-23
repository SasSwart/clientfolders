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
	parentDirectory := strings.TrimSpace(filepath.Dir(destination))
	err := os.MkdirAll(parentDirectory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create destination parent directories: %w", err)
	}

	sourceInfo, err := os.Lstat(source)
	if err != nil {
		return fmt.Errorf("could not stat source file: %w", err)
	}

	// Ignore symlinks
	if sourceInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		return nil
	}

	if sourceInfo.IsDir() {
		err := CopyDir(destination, source)
		if err != nil {
			return fmt.Errorf("could not copy directory: %w", err)
		}
	} else {
		err := CopyFile(destination, source)
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

func CopyFile(destination, source string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destination, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open destination file: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}
	return nil
}
