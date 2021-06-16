package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func encodeTime(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
	pae.AppendString(t.Format(time.RFC3339))
}

func encodeLevel(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		pae.AppendString("DEBUG ")
	case zapcore.InfoLevel:
		pae.AppendString("INFO  ")
	case zapcore.WarnLevel:
		pae.AppendString("WARN  ")
	case zapcore.ErrorLevel:
		pae.AppendString("ERROR ")
	case zapcore.DPanicLevel:
		pae.AppendString("DPANIC")
	case zapcore.PanicLevel:
		pae.AppendString("PANIC ")
	case zapcore.FatalLevel:
		pae.AppendString("FATAL ")
	default:
		pae.AppendString("??????")
	}
}

func encodeCaller(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {
	// NOP
}

// CreateLogger builds a preconfigured Zap logger instance.
func CreateLogger() *zap.SugaredLogger {
	c := zap.NewProductionConfig()
	c.Encoding = "console"
	c.EncoderConfig.ConsoleSeparator = " "
	c.EncoderConfig.EncodeTime = encodeTime
	c.EncoderConfig.EncodeLevel = encodeLevel
	c.EncoderConfig.EncodeCaller = encodeCaller
	logger, _ := c.Build()
	sugar := logger.Sugar()
	return sugar
}
