package file

import (
	"errors"
	"fmt"
	"github.com/kr/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Finder struct {
	locations []string
}

func New(locations []string) *Finder {
	finder := Finder{
		locations: locations,
	}

	return &finder
}

func (f *Finder) Find(path string) (string, error) {
	for _, location := range f.locations {
		absolutePath := filepath.Join(location, path)

		if _, err := os.Stat(absolutePath); err == nil {
			return absolutePath, nil
		}
	}

	return "", errors.New(fmt.Sprintf("File %s does not exist", path))
}

func filenameMatchPatterns(filename string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, filename)

		if matched == true || err != nil {
			return true
		}
	}

	return false
}

func (f *Finder) List(ignorePatterns []string) []string {
	fileList := []string{}

	for _, location := range f.locations {
		walker := fs.Walk(location)

		for walker.Step() {
			err := walker.Err()

			if err != nil || filenameMatchPatterns(walker.Path(), ignorePatterns) {
				continue
			}

			fileList = append(fileList, walker.Path())
		}
	}

	return fileList
}
