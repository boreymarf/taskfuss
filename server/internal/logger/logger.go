package logger

import (
	"fmt"
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

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 15:04:05",
		NoColor:    false,
		PartsOrder: []string{"level", "time", "message"},
	}

	output.FormatLevel = func(i any) string {
		var level string
		if l, ok := i.(string); ok {
			switch strings.ToLower(l) {
			case "debug":
				level = "[\x1b[90mDEBUG\x1b[0m]"
			case "trace":
				level = "[\x1b[90mTRACE\x1b[0m]"
			case "info":
				level = "[INFO]"
			case "warn":
				level = "[\x1b[33mWARN\x1b[0m]"
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
		return fmt.Sprintf("%-7s  ", level) // Два пробела после скобки
	}

	Log = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}
