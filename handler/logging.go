package handler

import (
	"Control/types"
	"io"
	"log/slog"
	"os"

	"github.com/robfig/cron/v3"
	"gopkg.in/natefinch/lumberjack.v2"
)

func LoggingHandler() {
	willi := &lumberjack.Logger{
		Filename:   types.Config.Log_Filename,
		MaxAge:     types.Config.Log_Max_Age,
		MaxBackups: types.Config.Log_Max_Backups,
		Compress:   types.Config.Log_compress,
	}

	multiwriter := io.MultiWriter(willi, os.Stdout)

	logger := slog.New(slog.NewJSONHandler(multiwriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		if err := willi.Rotate(); err != nil {
			logger.Error("log rotation error", "err", err)
			return
		}
		logger.Info("log rotated")
	})

	c.Start()
}
