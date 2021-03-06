package log

import (
	"log"
	"theater/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var SLogger *zap.SugaredLogger

func init() {
	var err error
	var cfg zap.Config
	var ecfg zapcore.EncoderConfig

	env := config.GetRuntimeEnv()
	switch env {
	case "development":
		cfg = zap.NewDevelopmentConfig()
		ecfg = zap.NewDevelopmentEncoderConfig()
	case "production":
		cfg = zap.NewProductionConfig()
		ecfg = zap.NewProductionEncoderConfig()
	}

	ecfg.MessageKey = "message"
	ecfg.LevelKey = "level"
	ecfg.TimeKey = "time"
	ecfg.EncodeTime = zapcore.ISO8601TimeEncoder
	ecfg.CallerKey = "caller"

	cfg.EncoderConfig = ecfg
	cfg.OutputPaths = []string{"stdout", config.LogPath()}
	cfg.ErrorOutputPaths = []string{"stderr", config.LogPath()}
	cfg.Encoding = "json"
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	Logger, err = cfg.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	SLogger = Logger.Sugar()
}
