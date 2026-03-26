package utils

import (
	"os"
	"testing"
)

func TestSetEnvFromFile(t *testing.T) {
	SetEnvFromFile()

	env := os.Getenv("GH_TOKEN")

	if env == "" {
		t.Errorf("env is empty")
	}
}

type dataReturn struct {
	Url    string
	Labels []map[string]any
}

func TestJsonToGo(t *testing.T) {
	var variable dataReturn

	data, _ := os.Open("../fixtures/json-return.json")
	DecodeJson(data, &variable)

	if len(variable.Labels) < 1 {
		t.Error("An error happen")
	}
}
