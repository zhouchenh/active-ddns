package logger

import (
	"github.com/rs/zerolog"
	"io"
)

var (
	stdoutLogger = zerolog.New(stdoutConsoleWriter).With().Timestamp().Logger().Level(zerolog.Level(DefaultLogLevel))
	stderrLogger = zerolog.New(stderrConsoleWriter).With().Timestamp().Logger().Level(zerolog.Level(DefaultLogLevel))
)

func SetTimestamp(enabled bool) {
	if enabled {
		stdoutConsoleWriter.PartsExclude = nil
		stderrConsoleWriter.PartsExclude = nil
	} else {
		stdoutConsoleWriter.PartsExclude = []string{zerolog.TimestampFieldName}
		stderrConsoleWriter.PartsExclude = []string{zerolog.TimestampFieldName}
	}
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return stdoutLogger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return stdoutLogger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return stdoutLogger.Info()
}

// Warning starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warning() *zerolog.Event {
	return stdoutLogger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return stderrLogger.Error()
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return stderrLogger.Err(err)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return stderrLogger.Fatal()
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return stderrLogger.Panic()
}

// WithLevel starts a new message with level. Unlike Fatal and Panic
// methods, WithLevel does not terminate the program or stop the ordinary
// flow of a goroutine when used with their respective levels.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level Level) *zerolog.Event {
	if level >= ErrorLevel && level <= PanicLevel {
		return stderrLogger.WithLevel(zerolog.Level(level))
	}
	return stdoutLogger.WithLevel(zerolog.Level(level))
}

// Log starts a new message with no level. Setting GlobalLevel to Disabled
// will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return stdoutLogger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	stdoutLogger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	stdoutLogger.Printf(format, v...)
}

func Writer() io.Writer {
	return &stdoutLogger
}

func ErrorWriter() io.Writer {
	return &stderrLogger
}
