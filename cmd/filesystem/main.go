package main

import (
	"fmt"

	"github.com/nice-pink/goutil/pkg/filesystem"
)

func main() {
	// items := filesystem.FindAllStringsInFile("bin/test/test.txt", ".*xxx.com.*")
	// for _, item := range items {
	// 	fmt.Println(item)
	// }

	// val, err := filesystem.GetRegexInFile("bin/test/test.txt", ".*(xxx.com.*)", "${1}", false)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(val)

	// val, err := filesystem.GetAllRegexInFile("bin/test/test.txt", ".*(xxx.com.*)", "${1}", false)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(val)

	val, err := filesystem.GetRegexInAllFiles("bin/test/", true, ".*(xxx.com.*)", "${1}")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)

}
