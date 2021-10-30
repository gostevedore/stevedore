package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	errors "github.com/apenella/go-common-utils/error"
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

func NewLogger(w io.Writer, encoderType string) *Logger {

	encoder := generateEncoderEncoder(encoderType)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		zapcore.DebugLevel)
	logger := zap.New(core)

	return &Logger{
		logger: logger.Sugar(),
	}
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
	l.logger.Info(msg)
}

// Warn
func (l *Logger) Warn(msg ...interface{}) {
	l.logger.Warn(msg)
}

// Error
func (l *Logger) Error(msg ...interface{}) {
	l.logger.Error(msg)
}

// Debug
func (l *Logger) Debug(msg ...interface{}) {
	l.logger.Debug(msg)
}

// Fatal
func (l *Logger) Fatal(msg ...interface{}) {
	l.logger.Fatal(msg)
}

// Panic
func (l *Logger) Panic(msg ...interface{}) {
	l.logger.Panic(msg)
}

//Init initializes suggarlogger
func Init(logfile, encoderType string) error {
	var encoder zapcore.Encoder

	if sugarLogger == nil {
		writer, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New("(logger::Init)", fmt.Sprintf("Error opening log file '%s'", logfile), err)
		}

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = customTimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		switch encoderType {
		case LogConsoleEncoderName:
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		case LogJSONEncoderName:
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		default:
			return errors.New("(logger::Init)", fmt.Sprintf("Unknown encoder '%s'", encoder), err)
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(writer),
			zapcore.DebugLevel)
		logger := zap.New(core)
		sugarLogger = logger.Sugar()
	}

	return nil
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// Info
func Info(msg ...interface{}) {
	sugarLogger.Info(msg)
}

// Warn
func Warn(msg ...interface{}) {
	sugarLogger.Warn(msg)
}

// Error
func Error(msg ...interface{}) {
	sugarLogger.Error(msg)
}

// Debug
func Debug(msg ...interface{}) {
	sugarLogger.Debug(msg)
}

// Fatal
func Fatal(msg ...interface{}) {
	sugarLogger.Fatal(msg)
}

// Panic
func Panic(msg ...interface{}) {
	sugarLogger.Panic(msg)
}
