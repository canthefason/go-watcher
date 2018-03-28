package watcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func(w *Watcher) getGoFiles(rootDir string) ([]string, error) {
	var files []string
	ext := ".go"

	err := filepath.Walk(rootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func(w *Watcher) getImportsRoot(files []string,wd string) ([]string, error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return nil, ErrPathNotSet
	}

	importsFolders := make([]string, 0)
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(content), "\n")
		startImport := false

		for _, line := range lines {
			//start of import section
			if strings.Contains(line, "import") {
				startImport = true
			}

			if startImport {
				if strings.Index(line, w.watchRecursiveRoot) == -1 {
					continue
				}

				startIndex := strings.Index(line, "\"")
				endIndex := strings.LastIndex(line, "\"")

				if startIndex != -1 && endIndex != -1 {
					//import line
					importStr := filepath.Join(goPath, "src", line[startIndex+1:len(line)-1])
					if strings.HasPrefix(importStr, wd) {
						continue
					}

					if !sliceContains(importsFolders, importStr) {
						importsFolders = append(importsFolders, importStr)
					}
				}

				//end of import section
				if strings.Contains(line, ")") {
					break
				}
			}
		}
	}

	return importsFolders, nil
}
