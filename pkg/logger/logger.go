package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create a custom encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = "" // Disable stacktrace

	// Create a custom core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		config.Level,
	)

	// Create logger
	logger := zap.New(core)

	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, err error, fields ...zapcore.Field) {
	fields = append(fields, zap.Error(err))
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, err error, fields ...zapcore.Field) {
	fields = append(fields, zap.Error(err))
	l.Logger.Fatal(msg, fields...)
}
