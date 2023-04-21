package file

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Walk(root string) []string {
	var result []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	return result
}

type QueryResult struct {
	Filename string
	Line     int
	Content  string
}

func queryFile(wg *sync.WaitGroup, chans chan QueryResult, path string, query string) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if !strings.Contains(scanner.Text(), query) {
			continue
		}

		chans <- QueryResult{
			Line:     i,
			Content:  scanner.Text(),
			Filename: path,
		}
	}
}

func Grep(root string, query string) []QueryResult {
	result := make(chan QueryResult, 10)
	var wg sync.WaitGroup
	var res []QueryResult

	for _, f := range Walk(root) {
		if strings.Contains(f, query) {
			result <- QueryResult{Line: 1, Filename: f, Content: f}
		}

		wg.Add(1)
		go queryFile(&wg, result, f, query)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for c := range result {
		res = append(res, c)
	}
	return res
}
