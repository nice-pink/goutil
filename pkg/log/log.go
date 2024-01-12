package log

import "fmt"

func Info(logs ...any) {
	//params := append([]any{" "}, logs...)
	fmt.Println(logs...)
}

func Warn(logs ...any) {
	params := append([]any{"⚠️ Warn: "}, logs...)
	fmt.Println(params...)
}

func Error(logs ...any) {
	params := append([]any{"❌ Error: "}, logs...)
	fmt.Println(params...)
}

func Err(err error, logs ...any) {
	params := append([]any{"❌ Error: "}, logs...)
	fmt.Println(params...)
	fmt.Println(err)
}
