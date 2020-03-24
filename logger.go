package nplog

import "context"

type (
	LogLevel       string
	Field          map[string]interface{}
	LoggerInstance int
	LoggerErr      string
)

func (e LoggerErr) Error() string {
	return string(e)
}

const (
	// Uber's zap logger solution
	ZapLogger LoggerInstance = iota
)

const (
	//Debug has verbose message
	Debug LogLevel = "debug"
	//Info is default log level
	Info LogLevel = "info"
	//Warn is for logging messages about possible issues
	Warn LogLevel = "warn"
	//Error is for logging errors
	Error LogLevel = "error"
	//Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	Fatal LogLevel = "fatal"
)

const (
	ErrNotSupportedLoggerInstance = LoggerErr("failed to initialize logger: not supported instance")
)

// NpLogger is a simplified abstraction of the zap.Logger
type NpLogger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	With(fields ...Field) NpLogger
	For(ctx context.Context) NpLogger
}

// NpLoggerOption stores config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type NpLoggerOption struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      LogLevel

	// EnableFile determines if log to file is enable
	EnableFile bool

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string

	// FileJSONFormat determines if log should be printed to file in JSON format
	FileJSONFormat bool

	// FileMaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	FileMaxSize int

	// FileMaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	FileMaxBackups int

	// FileMaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	FileMaxAge int

	// FileCompress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	FileCompress bool

	// FileLevel is log level of file
	FileLevel LogLevel
}

func NewNpLogger(instance LoggerInstance, options NpLoggerOption) (NpLogger, error) {
	switch instance {
	case ZapLogger:
		return newZapLogger(options)
	default:
		return nil, ErrNotSupportedLoggerInstance
	}
}
