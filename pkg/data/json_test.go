package data

import (
	"testing"

	"github.com/nice-pink/goutil/pkg/log"
)

func TestGetJson(t *testing.T) {
	// compare file
	output := map[string]any{}
	err := GetJson("@../../test/test.json", &output)
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetJson: returned error.")
	}
	validateTestMap(output, "TestGetJson::GetJson::file", t)

	//compare string
	input := `{
	    "test": {
	        "object": {
	            "key": "value"
	        }
	    }
	}`

	// get json
	output2 := map[string]any{}
	err = GetJson(input, &output2)
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetJson: returned error.")
	}
	validateTestMap(output2, "TestGetJson::GetJson::string", t)

	// get jsonmap
	output3, err := GetJsonMap(input)
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetJson: returned error.")
	}
	validateTestMap(output3, "TestGetJson::GetJsonMap", t)

	// get json array
	output4, err := GetJsonArray("[" + input + "]")
	if err != nil {
		log.Err(err, "error")
		t.Error("TestGetJson: returned error.")
	}
	validateTestMap(output4[0], "TestGetJson::GetJsonArray", t)
}

func validateTestMap(o map[string]any, testInfo string, t *testing.T) {
	// validate map
	vObj := map[string]any{"key": "value"}
	vTest := map[string]any{}
	vTest["object"] = vObj
	v := map[string]any{}
	v["test"] = vTest

	// validate
	test := o["test"].(map[string]any)
	obj := test["object"].(map[string]any)
	if obj["key"] != "value" {
		t.Error(testInfo, "not valid want != got", v, o)
	}
}
