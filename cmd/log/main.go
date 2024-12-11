package main

import (
	"errors"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
)

// func main() {
// 	slog.Debug("debug log")
// 	slog.Info("info log")
// 	slog.Warn("warn log")
// 	slog.Error("error log")
// }

func main() {
	log.Debug("debug log")
	log.Info("info log")
	log.Warn("warn log")
	log.Warning("warning log")
	log.Error("error log")
	log.Err(errors.New("new error"), "error log")
	log.Critical("critical log")
	log.Notify("head log")
	log.Info()

	rl := log.NewRLog("test", 80, "Debug", "January 1th 2006, 01:04:05.100")
	os.Setenv("GU_REMOTE_LOG_DEBUG", "true")
	rl.Debug("test log")
}
