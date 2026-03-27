package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetFromRoot(path string) string {
	var Root, _ = filepath.Abs(os.Getenv("HOME") + "/workspace/api-workflow")

	return filepath.Join(Root, path)
}

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
		log.Fatal("Error to open file utils.go", err)
	}

	return string(data)
}

func Capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func SetEnvFromFile() {
	env := GetFromRoot(".env")
	data := GetFileData(env)
	vars := strings.SplitSeq(data, "\n")
	for variable := range vars {
		spt := strings.Split(variable, "=")
		os.Setenv(spt[0], spt[1])
	}
}

func DecodeJson(data io.Reader, variable any) error {
	decoder := json.NewDecoder(data)

	if err := decoder.Decode(&variable); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
