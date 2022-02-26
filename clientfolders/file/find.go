package file

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Action func(string) (interface{}, error)

func Find(directory string, patterns []string, action Action, valuesChan chan<- []interface{}, errChan chan<- error) {
	go func() {
		regex, err := regexp.Compile(patterns[0])
		if err != nil {
			errChan <- fmt.Errorf("could not search for files due to invalid regex: %w", err)
			return
		}

		lastLevel := len(patterns) == 1

		foundFiles, _ := os.ReadDir(directory)
		numDirectories := 0
		values := make([]interface{}, 0)

		subValuesChan := make(chan []interface{})
		subErrChan := make(chan error)

		for _, file := range foundFiles {
			name := file.Name()
			if !regex.Match([]byte(name)) {
				continue
			}
			path := strings.Join([]string{directory, name}, string(os.PathSeparator))

			if lastLevel {
				if action != nil {
					value, err := action(path)
					if err != nil {
						errChan <- fmt.Errorf("could not execute callback: %w", err)
						return
					}
					values = append(values, value)
				}
			} else if file.IsDir() {
				numDirectories += 1
				Find(path, patterns[1:], action, subValuesChan, subErrChan)
			}
		}

		if !lastLevel {
			for i := 0; i < numDirectories; i++ {
				select {
				case subValues := <-subValuesChan:
					values = append(values, subValues...)
				case err := <-subErrChan:
					errChan <- fmt.Errorf("could not traverse path: %w", err)
				}
			}
		}

		valuesChan <- values
	}()
}
