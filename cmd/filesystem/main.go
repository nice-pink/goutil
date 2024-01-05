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

	val, err := filesystem.GetRegexInAllFiles("bin/", true, ".*(xxx.com.*)", "${1}", ".txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Values")
	fmt.Println(val)

}
