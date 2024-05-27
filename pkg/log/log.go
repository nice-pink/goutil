package log

import (
	"flag"
	"fmt"
	"time"
)

func Newline() {
	fmt.Println()
}

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

func Flags(withNewLine bool) {
	if withNewLine {
		fmt.Println()
	}
	fmt.Println("Flags:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("-%s: %s\n", f.Name, f.Value)
	})
	if withNewLine {
		fmt.Println()
	}
}
