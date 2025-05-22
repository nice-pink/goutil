package data

import (
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
)

func ReadJsonOrYaml(filepath string, output any) error {
	if strings.HasSuffix(filepath, ".yaml") || strings.HasSuffix(filepath, ".yml") {
		log.Info("Read yaml")
		return GetYaml("@"+filepath, output)
	}
	log.Info("Read json")
	return GetJson("@"+filepath, output)
}
