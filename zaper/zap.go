// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/8/13

package zaper

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel  = "debug"
	InfoLevel   = "info"
	WarnLevel   = "warn"
	ErrorLevel  = "error"
	DpanicLevel = "dpanic"
	PanicLevel  = "panic"
	FatalLevel  = "fatal"
)

func NewZap(level string, w zapcore.WriteSyncer, fields ...zap.Field) *zap.Logger {
	var l zapcore.Level
	switch level {
	case DebugLevel:
		l = zap.DebugLevel
	case InfoLevel:
		l = zap.InfoLevel
	case WarnLevel:
		l = zap.WarnLevel
	case ErrorLevel:
		l = zap.ErrorLevel
	case DpanicLevel:
		l = zap.DPanicLevel
	case PanicLevel:
		l = zap.PanicLevel
	case FatalLevel:
		l = zap.FatalLevel
	default:
		l = zap.DebugLevel
	}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(i time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(i.Local().Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zap.CombineWriteSyncers(w),
		l,
	).With(fields) //自带node 信息
	//大于error增加堆栈信息
	return zap.New(core).WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
}
