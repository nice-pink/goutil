package data

import (
	"testing"

	"github.com/nice-pink/goutil/pkg/log"
)

func TestReadJsonOrYaml(t *testing.T) {
	// json
	j := map[string]any{}
	err := ReadJsonOrYaml("../../test/test.json", &j)
	if err != nil {
		log.Err(err, "read json error")
		t.Error("TestReadJsonOrYaml:: could not read json file")
	}

	// yaml
	y := map[string]any{}
	err = ReadJsonOrYaml("../../test/test.yaml", &y)
	if err != nil {
		log.Err(err, "read yaml error")
		t.Error("TestReadJsonOrYaml:: could not read yaml file")
	}

	// yaml - broken
	yb := map[string]any{}
	err = ReadJsonOrYaml("../../test/test_broken.yaml", &yb)
	if err == nil {
		t.Error("TestReadJsonOrYaml:: should throw error on broken yaml file")
	}
}
