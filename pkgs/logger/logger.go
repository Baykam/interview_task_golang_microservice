package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger ekosistemindeki tüm katmanların kullanacağı ortak arayüz.
// Bu sayede usecase veya repository katmanları arkada Zap olduğunu bilmez.
type Logger interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
	Fatalf(msg string, keysAndValues ...any)
}

// zapLogger struct'ı dışarıya kapatıldı (unexported). Dışarı sadece interface sızacak.
type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

// NewLogger artık somut bir struct değil, yukarıdaki "Logger" interface'ini döner!
func NewLogger() Logger {
	writeSyncer := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)

	rawLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // CallerSkip(1) logun atıldığı asıl satırı gösterir

	return &zapLogger{
		sugaredLogger: rawLogger.Sugar(),
	}
}

// Interface metotlarının gerçekleştirilmesi (Implementation)
func (l *zapLogger) Info(msg string, keysAndValues ...any) {
	l.sugaredLogger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...any) {
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Debug(msg string, keysAndValues ...any) {
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...any) {
	l.sugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalf(msg string, keysAndValues ...any) {
	l.sugaredLogger.Fatalf(msg, keysAndValues...)
}
