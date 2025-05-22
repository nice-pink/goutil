package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nice-pink/goutil/pkg/data"
)

func main() {
	inputJson := flag.String("inputJson", "", "Input json. Either JSON_STRING or @JSON_FILE.")
	inputJsonArr := flag.String("inputJsonArr", "", "Input json array. Either JSON_STRING or @JSON_FILE.")
	flag.Parse()

	if *inputJson != "" {
		fmt.Println(inputJson)

		jsonMap, err := data.GetJsonMap(*inputJson)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Println(jsonMap)
	}

	if *inputJsonArr != "" {
		fmt.Println(inputJsonArr)

		jsonArr, err := data.GetJsonArray(*inputJsonArr)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Println(jsonArr)
	}
}

// error:
// bin/jsoninput -inputJson '{"key"}'

// string:
// bin/jsoninput -inputJson '{"key":"value","other":"v"}'

// string-arr:
// bin/jsoninput -inputJsonArr '[{"key":"value"},{"key":"value"}]'

// file
// bin/jsoninput -inputJson @cmd/jsoninput/test.json

// file-arr:
// bin/jsoninput -inputJsonArr @cmd/jsoninput/test_arr.json
