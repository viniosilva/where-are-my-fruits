package infra

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigLogger() *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()

	return logger.Sugar()
}

// Refers: https://github.com/uber-go/zap/blob/master/example_test.go
