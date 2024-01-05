package kustomize

import (
	"fmt"

	"github.com/nice-pink/goutil/pkg/filesystem"
)

func CreateKustomization(path string) {
	pathKustomization := path + "/kustomization.yaml"
	if !filesystem.FileExists(pathKustomization) {
		err := filesystem.AppendToFile(pathKustomization, `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:`, true)
		if err != nil {
			panic(err)
		}
	}
}

// metadata

func UpdateNamespace(path string, namespace string) error {
	fmt.Println("Update namespace: " + namespace)
	return addItemToKustomization(path, namespace, "namespace", false)
}

// resource

func AddResourceToKustomization(path string, resource string) error {
	fmt.Println("Add resource: " + resource + " to " + path)
	return addItemToKustomization(path, resource, "resources", true)
}

// compontent

func AddComponentToKustomization(path string, component string) error {
	fmt.Println("Add resource: " + component + " to " + path)
	return addItemToKustomization(path, component, "components", true)
}

// general

func addItemToKustomization(path string, item string, itemType string, addNewline bool) error {
	pathKustomization := path + "/kustomization.yaml"

	// get key from itemType
	key := itemType + ":"

	// add key if does not exist
	if filesystem.ContainsString(pathKustomization, itemType) {
		err := filesystem.AppendToFile(path, key, addNewline)
		if err != nil {
			return err
		}
	}

	// add item
	replacement := key
	if addNewline {
		replacement = replacement + "\n" + item
	} else {
		replacement = replacement + " " + item
	}
	return filesystem.ReplaceInFile(pathKustomization, key, replacement, true)
}
