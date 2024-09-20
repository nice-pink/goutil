package log

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func Newline() {
	fmt.Println()
}

func Verbose(logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " VERBOSE:"}, logs...)
	fmt.Println(params...)
}

func Debug(logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " DEBUG:"}, logs...)
	fmt.Println(params...)
}

func Info(logs ...any) {
	params := append([]any{time.Now().Format(time.DateTime) + " INFO:"}, logs...)
	fmt.Println(params...)
}

func Warn(logs ...any) {
	var params []any
	if strings.ToLower(os.Getenv("LOG_IGNORE_COLOR")) == "true" {
		params = append([]any{time.Now().Format(time.DateTime) + " WARNING:"}, logs...)
	} else {
		params = append([]any{Yellow + time.Now().Format(time.DateTime) + " WARNING:" + Reset}, logs...)
	}
	fmt.Println(params...)
}

func Warning(logs ...any) {
	Warn(logs...)
}

func Error(logs ...any) {
	var params []any
	if strings.ToLower(os.Getenv("LOG_IGNORE_COLOR")) == "true" {
		params = append([]any{time.Now().Format(time.DateTime) + " ERROR:"}, logs...)
	} else {
		params = append([]any{Red + time.Now().Format(time.DateTime) + " ERROR:" + Reset}, logs...)
	}
	fmt.Println(params...)
}

func Err(err error, logs ...any) {
	var params []any
	if strings.ToLower(os.Getenv("LOG_IGNORE_COLOR")) == "true" {
		params = append([]any{time.Now().Format(time.DateTime) + " ERROR:"}, logs...)
	} else {
		params = append([]any{Red + time.Now().Format(time.DateTime) + " ERROR:" + Reset}, logs...)
	}
	fmt.Println(params...)
	fmt.Println(err)
}

func Critical(logs ...any) {
	var params []any
	if strings.ToLower(os.Getenv("LOG_IGNORE_COLOR")) == "true" {
		params = append([]any{time.Now().Format(time.DateTime) + " ERROR:"}, logs...)
	} else {
		params = append([]any{Magenta + time.Now().Format(time.DateTime) + " ERROR:" + Reset}, logs...)
	}
	fmt.Println(params...)
}

// special

func Notify(logs ...any) {
	var params []any
	if strings.ToLower(os.Getenv("LOG_IGNORE_COLOR")) == "true" {
		params = append([]any{time.Now().Format(time.DateTime) + " INFO:"}, logs...)
	} else {
		params = append([]any{Blue + time.Now().Format(time.DateTime) + " INFO:" + Reset}, logs...)
	}
	fmt.Println(params...)
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
		fmt.Printf(Blue+"-%s: %s\n"+Reset, f.Name, f.Value)
	})
	if withNewLine {
		fmt.Println()
	}
}
