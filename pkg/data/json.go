package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func DumpJson(i any, filepath string) {
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

func GetJson(input string, output any) error {
	data, err := GetData(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &output)
}

func GetData(input string) ([]byte, error) {
	if !strings.HasPrefix(input, "@") {
		// is json input string
		return []byte(input), nil
	}

	// is json input file
	return os.ReadFile(strings.TrimPrefix(input, "@"))
}

func GetJsonMap(input string) (map[string]any, error) {
	var output map[string]any
	err := GetJson(input, &output)
	return output, err
}

func GetJsonArray(input string) ([]map[string]any, error) {
	var output []map[string]any
	err := GetJson(input, &output)
	return output, err
}
