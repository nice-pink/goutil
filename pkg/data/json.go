package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func DumpJson(i interface{}, filepath string) {

	j, _ := json.MarshalIndent(i, "", "  ")
	// fmt.Println(string(j))

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if _, err := file.Write(j); err != nil {
		fmt.Println(err)
	}
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func GetJson(input string) (map[string]interface{}, error) {
	var output map[string]interface{}

	if !strings.HasPrefix(input, "@") {
		// is json input string
		err := json.Unmarshal([]byte(input), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	}

	// is json input file
	data, err := os.ReadFile(strings.TrimPrefix(input, "@"))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &output)
	return output, err
}

func GetJsonArray(input string) ([]map[string]interface{}, error) {
	var output []map[string]interface{}

	if !strings.HasPrefix(input, "@") {
		// is json input string
		err := json.Unmarshal([]byte(input), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	}

	// is json input file
	data, err := os.ReadFile(strings.TrimPrefix(input, "@"))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &output)
	return output, err
}
