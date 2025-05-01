package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006/01/02 15:04:05",
		NoColor:    false,
		PartsOrder: []string{"level", "time", "message"},
	}

	output.FormatLevel = func(i any) string {
		var level string
		if l, ok := i.(string); ok {
			switch l {
			case "debug":
				level = fmt.Sprintf("%-5s", "DEBUG")
			case "info":
				level = fmt.Sprintf("%-5s", "INFO")
			case "warn":
				level = fmt.Sprintf("\x1b[30;43m%-5s\x1b[0m", "WARN") // Black text/Yellow background
			case "error":
				level = fmt.Sprintf("\x1b[30;41m%-5s\x1b[0m", "ERROR") // Black text/Red background
			default:
				level = strings.ToUpper(fmt.Sprintf("%-5s", l))
			}
		}
		return fmt.Sprintf("[%s]", level)
	}

	Log = zerolog.New(output).
		With().
		Timestamp().
		Logger()
}
