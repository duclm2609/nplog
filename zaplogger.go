package nplog

import (
	"context"
	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// zapLogger delegates all calls to the underlying zap.Logger
type zapLogger struct {
	logger *zap.SugaredLogger
}

// Debug logs message with debug level
func (z zapLogger) Debug(msg string, args ...interface{}) {
	z.logger.Debug(msg, args)
}

// Info logs message with info level
func (z zapLogger) Info(msg string, args ...interface{}) {
	z.logger.Info(msg, args)
}

// Error logs message with error level
func (z zapLogger) Error(msg string, args ...interface{}) {
	z.logger.Error(msg, args)
}

// Fatal logs a fatal error message
func (z zapLogger) Fatal(msg string, args ...interface{}) {
	z.logger.Fatal(msg, args)
}

// With creates a child logger, and optionally adds some context fields to that logger
func (z zapLogger) With(fields ...Field) NpLogger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	return zapLogger{z.logger.With(f...)}
}

// For return Elastic APM trace context aware, if available
func (z zapLogger) For(ctx context.Context) NpLogger {
	if traceCtx := apmzap.TraceContext(ctx); traceCtx != nil {
		return zapLogger{logger: z.logger.With(traceCtx)}
	}
	return z
}

// getZapLevel maps with zap log level, default to INFO
func getZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoder return Elastic ECS schema compatible
func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "error.stack_trace"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func newZapLogger(cfg NpLoggerOption) (NpLogger, error) {
	var cores []zapcore.Core

	if cfg.EnableConsole {
		consoleLevel := getZapLevel(cfg.ConsoleLevel)
		writer := zapcore.Lock(os.Stdout)
		consoleCore := zapcore.NewCore(getEncoder(cfg.ConsoleJSONFormat), writer, consoleLevel)
		cores = append(cores, consoleCore)
	}

	if cfg.EnableFile {
		//TODO: default value for configuration
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.FileMaxSize,
			MaxAge:     cfg.FileMaxBackups,
			MaxBackups: cfg.FileMaxAge,
			Compress:   cfg.FileCompress,
		})

		fileCore := zapcore.NewCore(getEncoder(cfg.FileJSONFormat), w, getZapLevel(cfg.FileLevel))
		cores = append(cores, fileCore)
	}

	combinedCores := zapcore.NewTee(cores...)
	logger := zap.New(combinedCores,
		zap.WrapCore((&apmzap.Core{}).WrapCore),
		zap.AddCaller(),
		zap.AddCallerSkip(2)).Sugar()
	return zapLogger{logger: logger}, nil
}
