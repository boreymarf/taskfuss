package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 15:04:05",
		NoColor:    false,
		PartsOrder: []string{"level", "time", "caller", "message"},
	}
	consoleWriter.FormatLevel = func(i any) string {
		var level string
		if l, ok := i.(string); ok {
			switch strings.ToLower(l) {
			case "debug":
				level = "[\x1b[90mDEBUG\x1b[0m]"
			case "trace":
				level = "[\x1b[90mTRACE\x1b[0m]"
			case "info":
				level = "[INFO] "
			case "warn":
				level = "[\x1b[33mWARN\x1b[0m] "
			case "error":
				level = "[\x1b[31mERROR\x1b[0m]"
			case "fatal":
				level = "[\x1b[37;41mFATAL\x1b[0m]"
			case "panic":
				level = "[\x1b[37;41mPANIC\x1b[0m]"
			default:
				level = strings.ToUpper(l)
			}
		}
		return level
	}

	consoleWriter.FormatCaller = func(i any) string {
		if caller, ok := i.(string); ok {
			// Find the last slash in the path
			lastSlash := strings.LastIndex(caller, "/")
			if lastSlash != -1 {
				// Look for the previous slash (to extract the folder)
				prevSlash := -1
				if lastSlash > 0 {
					prevSlash = strings.LastIndex(caller[:lastSlash], "/")
				}
				// Trim to the last folder (if present)
				if prevSlash != -1 {
					caller = caller[prevSlash+1:]
				}
				// Remove leading slash (if still present)
				if len(caller) > 0 && caller[0] == '/' {
					caller = caller[1:]
				}
			}
			return fmt.Sprintf(" \x1b[90m(%s)\x1b[0m", caller)
		}
		return ""
	}

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	multi := io.MultiWriter(consoleWriter, file)

	Log = zerolog.New(multi).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}
