package setup

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"log"
	"{{{mytemplate}}}/config/env"
	"os"
	"time"
)

var logLevel = map[string]zerolog.Level{
	"TRACE": zerolog.TraceLevel,
	"OFF":   zerolog.Disabled,
	"DEBUG": zerolog.DebugLevel,
	"INFO":  zerolog.InfoLevel,
	"WARN":  zerolog.WarnLevel,
	"ERROR": zerolog.ErrorLevel,
	"FATAL": zerolog.FatalLevel,
	"PANIC": zerolog.PanicLevel,
}

var timeFormat = map[string]string{
	"UNIX":      zerolog.TimeFormatUnix,
	"UNIXMS":    zerolog.TimeFormatUnixMs,
	"UNIXMICRO": zerolog.TimeFormatUnixMicro,
	"UNIXNANO":  zerolog.TimeFormatUnixNano,
	"RFC3339":   time.RFC3339,
	"RFC1123":   time.RFC1123,
}

// Return new logger instance with logfile for closing when shutdown server
func NewLogger(env *env.Env) (zerolog.Logger, io.WriteCloser) {
	var logger zerolog.Logger
	var location io.WriteCloser
	var level zerolog.Level
	var format string

	if env.Log.LogLocation == "STD" {
		location = os.Stderr
	} else {
		logFile, err := os.OpenFile(env.Log.LogLocation, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Couldn't set log location %s : %v", env.Log.LogLocation, err)
		}
		location = logFile
	}

	if !env.App.Production {
		writer := zerolog.ConsoleWriter{Out: location}
		logger = zerolog.New(writer)
	} else {
		logger = zerolog.New(location)
	}

	if lvl, ok := logLevel[env.Log.LogLevel]; ok {
		level = lvl
	} else {
		level = zerolog.DebugLevel
	}

	if fm, ok := timeFormat[env.Log.TimeFieldFormat]; ok {
		format = fm
	} else {
		format = timeFormat["UNIX"]
	}

	loggerS := logger.Sample(&zerolog.LevelSampler{
		DebugSampler: &zerolog.BurstSampler{
			Burst:       env.Log.DebugBurstLvl,
			Period:      env.Log.DebugBurstPeriod,
			NextSampler: &zerolog.BasicSampler{N: env.Log.DebugN},
		},
		InfoSampler: &zerolog.BurstSampler{
			Burst:  env.Log.InfoBurstLvl,
			Period: env.Log.InfoBurstPertiod,
			NextSampler: &zerolog.BasicSampler{
				N: env.Log.DebugN,
			},
		},

		WarnSampler: &zerolog.BurstSampler{
			Burst:  env.Log.WarnBurstLvl,
			Period: env.Log.WarnBurstPertiod,
			NextSampler: &zerolog.BasicSampler{
				N: env.Log.WarnN,
			},
		},
	}).With().Timestamp().Logger()
	logger = loggerS

	zerolog.SetGlobalLevel(level)
	zerolog.TimestampFieldName = env.Log.TimeFieldName
	zerolog.MessageFieldName = env.Log.MessageFieldName
	zerolog.ErrorFieldName = env.Log.ErrorFieldName
	zerolog.TimeFieldFormat = format
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return logger, location

}
