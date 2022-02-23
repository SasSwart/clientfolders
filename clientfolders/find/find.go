package find

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Action func(string, chan error)

func Find(directory string, patterns []string, action Action) ([]string, error) {
	regex, err := regexp.Compile(patterns[0])
	if err != nil {
		return nil, fmt.Errorf("could not search for files due to invalid regex: %w", err)
	}

	lastLevel := len(patterns) == 1

	foundFiles, _ := os.ReadDir(directory)
	filteredFiles := make([]string, 0)

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
					return nil, fmt.Errorf("could not execute callback: %w", err)
				}
			}

			filteredFiles = append(filteredFiles, path)
		} else if file.IsDir() {
			subfiles, err := Find(path, patterns[1:], action)
			if err != nil {
				return nil, fmt.Errorf("could not traverse path: %w", err)
			}
			filteredFiles = append(filteredFiles, subfiles...)
		}
	}

	return filteredFiles, nil
}
