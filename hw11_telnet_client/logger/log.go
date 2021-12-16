package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Zap() *zap.SugaredLogger {
	if logger == nil {
		logger = initLog()
	}
	return logger
}

func init() {
	_ = Zap()
}

func Finalize() {
	if logger != nil {
		Zap().Sync()
	}
}

func initLog() *zap.SugaredLogger {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
