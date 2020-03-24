package nplog

import (
	"context"
	"go.elastic.co/apm/module/apmzap"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NPLogger is the default logging wrapper that can create
// zapLogger instances either for a given Context or context-less.
type NPLogger struct {
	logger *zap.Logger
}

// NewNPLogger creates a new NPLogger.
func NewNPLogger(logger *zap.Logger) NPLogger {
	return NPLogger{logger: logger}
}

// Bg creates a context-unaware zapLogger.
func (b NPLogger) Bg() Logger {
	return zapLogger(b)
}

// For returns a Elastic APM context-aware Logger, if available
func (b NPLogger) For(ctx context.Context) Logger {
	if traceCtx := apmzap.TraceContext(ctx); traceCtx != nil {
		return zapLogger{logger: b.logger.With(traceCtx...)}
	}
	return b.Bg()
}

// With creates a child zapLogger, and optionally adds some context fields to that zapLogger.
func (b NPLogger) With(fields ...zapcore.Field) NPLogger {
	return NPLogger{logger: b.logger.With(fields...)}
}
