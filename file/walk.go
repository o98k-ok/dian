package file

import (
	"io/ioutil"
	"path"
)

func Walk(root string) ([]string, error) {
	files, err := ioutil.ReadDir(root) // 读取目录信息
	if err != nil {
		return nil, err
	}

	var results []string
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			results = append(results, path.Join(root, file.Name()))
		}
	}
	return results, nil
}
