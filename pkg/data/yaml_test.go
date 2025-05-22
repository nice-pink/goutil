package data

import (
	"testing"

	"github.com/nice-pink/goutil/pkg/log"
)

func TestGetYaml(t *testing.T) {
	// compare file
	output := map[string]any{}
	err := GetYaml("@../../test/test.yaml", &output)
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetYaml: returned error.")
	}
	validateTestMap(output, "TestGetYaml::GetYaml::file", t)

	//compare string
	input := `test:
  object:
    key: value`

	// get yaml
	output2 := map[string]any{}
	err = GetYaml(input, &output2)
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetYaml: returned error.")
	}
	validateTestMap(output2, "TestGetYaml::GetYaml::string", t)

	// get yamlmap
	output3, err := GetYamlMap(input)
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetYaml: returned error.")
	}
	validateTestMap(output3, "TestGetYaml::GetYamlMap", t)

	// get yaml array
	inputArr := `- test:
    object:
      key: value`

	output4, err := GetYamlArray(inputArr)
	if err != nil {
		log.Err(err, "error")
		t.Error("GetYamlArray: returned error.")
	}
	validateTestMap(output4[0], "TestGetYaml::GetYamlArray", t)
}
