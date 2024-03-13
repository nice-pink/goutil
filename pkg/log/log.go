package log

import (
	"fmt"
	"time"
)

func Info(logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " ->"}, logs...)
	fmt.Println(params...)
}

func Warn(logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " ->", "🔶 Warn:"}, logs...)
	fmt.Println(params...)
}

func Error(logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " ->", "❌ Error:"}, logs...)
	fmt.Println(params...)
}

func Err(err error, logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " ->", "❌ Error:"}, logs...)
	fmt.Println(params...)
	fmt.Println(err)
}

func Time() {
	fmt.Println(time.Now())
}
