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
		var m map[string]any
		filepath := strings.TrimPrefix(value, "@")
		data, err := os.ReadFile(filepath)
		if err != nil {
			log.Err(err, "payload from file", filepath)
		}
		err = json.Unmarshal(data, &m)
		if err != nil {
			log.Err(err, "unmarshal payload from file", filepath)
		}
		data, _ = json.Marshal(m)
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
