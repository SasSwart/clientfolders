package file

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Action func(string, chan error)

func Find(directory string, patterns []string, action Action, filesChan chan<- []string, errChan chan<- error) {
	go func() {
		regex, err := regexp.Compile(patterns[0])
		if err != nil {
			errChan <- fmt.Errorf("could not search for files due to invalid regex: %w", err)
			return
		}

		lastLevel := len(patterns) == 1

		foundFiles, _ := os.ReadDir(directory)
		numDirectories := 0
		filteredFiles := make([]string, 0)

		subfilesChan := make(chan []string)
		subErrChan := make(chan error)

		for _, file := range foundFiles {
			name := file.Name()
			if !regex.Match([]byte(name)) {
				continue
			}
			path := strings.Join([]string{directory, name}, string(os.PathSeparator))

			if lastLevel {
				if action != nil {
					errchan := make(chan error)
					go action(path, errchan)
					err := <-errchan
					if err != nil {
						errChan <- fmt.Errorf("could not execute callback: %w", err)
						return
					}
				}

				filteredFiles = append(filteredFiles, path)
			} else if file.IsDir() {
				numDirectories += 1
				Find(path, patterns[1:], action, subfilesChan, subErrChan)
			}
		}
		if !lastLevel {
			for i := 0; i < numDirectories; i++ {
				select {
				case files := <-subfilesChan:
					filteredFiles = append(filteredFiles, files...)
				case err := <-subErrChan:
					errChan <- fmt.Errorf("could not traverse path: %w", err)
				}
			}
		}

		filesChan <- filteredFiles
	}()
}
