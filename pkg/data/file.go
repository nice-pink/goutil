package data

import (
	"strings"
)

func ReadJsonOrYaml(filepath string, output any) error {
	if strings.HasSuffix(filepath, ".yaml") || strings.HasSuffix(filepath, ".yml") {
		return GetYaml(filepath, output)
	}
	return GetJson(filepath, output)
}
