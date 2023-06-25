package logs

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// func InitLogger() {
// 	var writers []io.Writer
// 	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
// 	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
// 	output.FormatLevel = func(i interface{}) string {
// 		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
// 	}
// 	output.FormatMessage = func(i interface{}) string {
// 		return fmt.Sprintf("message : %s ", i)
// 	}
// 	output.FormatFieldName = func(i interface{}) string {
// 		return fmt.Sprintf("| %s:", i)
// 	}
// 	output.FormatFieldValue = func(i interface{}) string {
// 		return fmt.Sprintf("%s", i)
// 	}

// 	writers = append(writers, output)
// 	mw := io.MultiWriter(writers...)
// 	logFile, _ := os.OpenFile("auth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
// 	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

// 	log.Logger = zerolog.New(mw).With().Timestamp().Logger().Output(output)

// 	log.Trace().Msg("zerolog initialize")
// }

type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

// type Logger struct {
// 	*zerolog.Logger
// }

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func InitLogger(config Config) {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("message : %s ", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("| %s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()

	log.Logger = logger

	logger.Output(output)
	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	// return &Logger{
	// 	Logger: &logger,
	// }
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0o744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
