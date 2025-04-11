package log

import (
	"log"
	"os"

	"github.com/OutOfStack/game-library/internal/appconf"
	gelf "github.com/snovichkov/zap-gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns new zap logger instance
func New(cfg appconf.Cfg) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleWriter := zapcore.Lock(os.Stderr)
	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), consoleWriter, zap.InfoLevel),
	}

	host, _ := os.Hostname()
	gelfCore, err := gelf.NewCore(gelf.Addr(cfg.Graylog.Address), gelf.Host(host))
	if err != nil {
		log.Printf("can't create gelf core: %v", err)
	}
	if gelfCore != nil {
		cores = append(cores, gelfCore)
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.WithCaller(false)).With(zap.String("service", appconf.ServiceName))

	return logger
}
