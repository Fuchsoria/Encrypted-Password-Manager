package main

import (
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(debug bool) *zap.SugaredLogger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join("./", "logs", "bot.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	cfg := zap.NewProductionEncoderConfig()

	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.TimeKey = "time"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		w,
		zap.InfoLevel,
	)

	if debug {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			w,
			zap.NewDevelopmentConfig().Level,
		)
	}

	logger := zap.New(core)

	return logger.Sugar()
}
