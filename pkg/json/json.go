package json

import (
	"encoding/json"
	"fmt"
	"os"
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
