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
	fmt.Println(logLevel)
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
		PartsOrder: []string{"level", "time", "caller", "message"},
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

	output.FormatCaller = func(i any) string {
		if caller, ok := i.(string); ok {
			// Находим последний слеш в пути
			lastSlash := strings.LastIndex(caller, "/")
			if lastSlash != -1 {
				// Ищем предыдущий слеш (для выделения папки)
				prevSlash := -1
				if lastSlash > 0 {
					prevSlash = strings.LastIndex(caller[:lastSlash], "/")
				}
				// Обрезаем до последней папки (если есть)
				if prevSlash != -1 {
					caller = caller[prevSlash+1:]
				}
				// Удаляем ведущий слеш (если остался)
				if len(caller) > 0 && caller[0] == '/' {
					caller = caller[1:]
				}
			}
			return fmt.Sprintf(" \x1b[90m(%s)\x1b[0m", caller)
		}
		return ""
	}

	Log = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
}
