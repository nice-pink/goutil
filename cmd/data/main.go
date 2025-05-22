package main

import (
	"fmt"

	"github.com/nice-pink/goutil/pkg/data"
)

func main() {
	o := map[string]any{}
	err := data.GetYaml("@cmd/data/test.yaml", &o)
	if err != nil {
		fmt.Println(err.Error())
	}
	print(o)

	m, err := data.GetYamlMap("@cmd/data/test.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	print(m)

}

func print(o map[string]any) {
	fmt.Println(o)

	test := o["test"].(map[string]any)
	fmt.Println(test)

	obj := test["object"].(map[string]any)
	fmt.Println(obj)

	fmt.Println(obj["key"])
}
