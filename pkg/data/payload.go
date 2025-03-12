package data

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
)

func GetPayload(value string) []byte {
	if value == "" {
		// nothing to return
		return nil
	}

	if strings.HasPrefix(value, "@") {
		// get payload from file
		filepath := strings.TrimPrefix(value, "@")
		data, err := os.ReadFile(filepath)
		if err != nil {
			log.Err(err, "payload from file", filepath)
		}
		return data
	}

	// return json string as data
	return []byte(value)
}

func GetPayloadMap(value string) map[string]interface{} {
	data := GetPayload(value)

	m := map[string]interface{}{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Err(err, "unmarshal payload", string(data))
	}
	return m
}
