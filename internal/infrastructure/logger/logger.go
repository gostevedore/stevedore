package logger

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogConsoleEncoderName = "console"
	LogJSONEncoderName    = "json"
)

var sugarLogger *zap.SugaredLogger

type Logger struct {
	logger *zap.SugaredLogger
}

// New creates a new logger
func New() *Logger {

	l, _ := zap.NewProduction()
	return &Logger{
		logger: l.Sugar(),
	}
}

// ReloadWithWriter recreates the logger with a new writer
func (l *Logger) ReloadWithWriter(w io.Writer) {
	encoder := generateEncoderEncoder(LogConsoleEncoderName)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		zapcore.DebugLevel,
	)

	logger := zap.New(core)
	l.logger = logger.Sugar()
}

func (l *Logger) Sync() {
	l.logger.Sync()
}

func generateEncoderEncoder(encoderType string) zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	switch encoderType {
	case LogJSONEncoderName:
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

}

// Info
func (l *Logger) Info(msg ...interface{}) {
	l.logger.Info(msg...)
}

// Warn
func (l *Logger) Warn(msg ...interface{}) {
	l.logger.Warn(msg...)
}

// Error
func (l *Logger) Error(msg ...interface{}) {
	l.logger.Error(msg...)
}

// Debug
func (l *Logger) Debug(msg ...interface{}) {
	l.logger.Debug(msg...)
}

// Fatal
func (l *Logger) Fatal(msg ...interface{}) {
	l.logger.Fatal(msg...)
}

// Panic
func (l *Logger) Panic(msg ...interface{}) {
	l.logger.Panic(msg...)
}

// customTimeEncoder is a custom time encoder for logger
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
