package logger

// use zero log for logging
import (
	"chroma-db/internal/constants"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Log is the global logger
var Log zerolog.Logger

// Init initializes the logger
func init() {

	// Default level for this example is info, unless debug flag is present
	if constants.LogLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s:-", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()
	log = log.With().Caller().Logger()

	Log = log
}
