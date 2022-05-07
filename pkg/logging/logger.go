package logging

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const timeFormat = "2006-01-02T15:04:05.000MST"

func init() {
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		arr := strings.Split(file, "/")
		return arr[len(arr)-1] + ":" + strconv.Itoa(line)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}).Level(zerolog.InfoLevel)
}

func NewLogger() zerolog.Logger {
	return log.With().Timestamp().Caller().Logger()
}
