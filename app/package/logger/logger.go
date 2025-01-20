package logger

import (
	"log/slog"
	"os"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

// type CtxKey string

// var LogKey CtxKey = "logger"

type Logger struct {
	*slog.Logger

	logFile *os.File
}

func New(env, logFolder string) Logger {
	var logger Logger

	switch env {
	case envLocal:
		logger = Logger{
			Logger: slog.New(NewHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})),
			logFile: nil,
		}

	case envDev:
		logger = Logger{
			Logger: slog.New(NewHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			})),
			logFile: nil,
		}

	case envProd:
		logFile, err := os.OpenFile(getLogFileName(logFolder), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		logger = Logger{
			Logger: slog.New(NewHandler(logFile, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})),
			logFile: logFile,
		}
	}

	return logger
}

// func ContextWithLogger(ctx context.Context, env, logFolder string) context.Context {
// 	return context.WithValue(ctx, LogKey, New(env, logFolder))
// }

// func FromContext(ctx context.Context) Logger {
// 	return ctx.Value(LogKey).(Logger)
// }

func (l Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func (l Logger) Error(msg string, err error) {
	if err != nil {
		l.Logger.Error(msg, slog.Attr{
			Key:   "err",
			Value: slog.StringValue(err.Error()),
		})
	} else {
		l.Logger.Error(msg)
	}
}

func getLogFileName(logFolder string) string {
	return logFolder + time.Now().Format("02/01 15:04:05.000") + ".log"
}
