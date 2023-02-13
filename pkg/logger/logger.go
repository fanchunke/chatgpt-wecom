package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int

	Level    string
	Filename string
}

func (c *Config) newRollingFileWriter() io.Writer {
	return &lumberjack.Logger{
		Filename:   c.Filename,
		LocalTime:  true,
		MaxBackups: c.MaxBackups, // files
		MaxSize:    c.MaxSize,    // megabytes
		MaxAge:     c.MaxAge,     // days
	}
}

func Configure(config Config) {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writer := zerolog.ConsoleWriter{Out: os.Stderr}
		writers = append(writers, customConsoleWriterFormat(writer))
	}
	if config.FileLoggingEnabled {
		writers = append(writers, config.newRollingFileWriter())
	}

	mw := zerolog.MultiLevelWriter(writers...)

	var level zerolog.Level
	switch config.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(mw).With().Caller().Logger()
}

func customConsoleWriterFormat(writer zerolog.ConsoleWriter) zerolog.ConsoleWriter {
	writer.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%-6s", i))
	}
	writer.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	writer.TimeFormat = time.RFC3339
	return writer
}
