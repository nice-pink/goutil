package main

import (
	"errors"

	"github.com/nice-pink/goutil/pkg/log"
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

	log.Info("bla")
	log.Warn("warn")
	log.Error("This is error")
	err := errors.New("bla")
	log.Err(err, "This is error")

	// extensions := []string{".txt"}
	// val, err := filesystem.GetRegexInAllFiles("bin/", true, ".*(xxx.com.*)", "${1}", extensions)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Values")
	// fmt.Println(val)

}
