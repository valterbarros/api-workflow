package utils

import (
	"os"
	"strings"
)

func IterateFiles(path, ignore string, fn func(name string)) {
	files, _ := os.ReadDir(path)

	for _, file := range files {
		name := file.Name()

		if !strings.HasSuffix(name, ignore) {
			continue
		}

		fn(name)
	}
}

func SaveFile(file, data string) {
	os.WriteFile(file, []byte(data), 0644)
}

func ResetFile(file string) {
	os.WriteFile(file, []byte(""), 0644)
}

func GetFileData(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		panic("Error to open file")
	}

	return string(data)
}

func Capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
