package data

import (
	"encoding/json"
	"strings"
)

func GetPayload(value string) []byte {
	if value == "" {
		// nothing to return
		return nil
	}

	if strings.HasPrefix(value, "@") {
		// get payload from file
		m, err := GetJsonMap(value)
		if err != nil {
			return nil
		}
		// marshal data to remove invalid chars
		data, _ := json.Marshal(m)
		return data
	}
	// return json string as data
	return []byte(value)
}

func GetJsonPayload(value string) map[string]any {
	m, _ := GetJsonMap(value)
	return m
}
