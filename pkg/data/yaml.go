package data

import (
	"k8s.io/apimachinery/pkg/util/yaml"
)

func GetYaml(input string, output any) error {
	data, err := GetData(input)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &output)
}

func GetYamlMap(input string) (map[string]any, error) {
	var output map[string]any
	err := GetYaml(input, &output)
	return output, err
}

func GetYamlArray(input string) ([]map[string]any, error) {
	var output []map[string]any
	err := GetYaml(input, &output)
	return output, err
}
