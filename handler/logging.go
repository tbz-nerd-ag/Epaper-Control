package handler

import (
	"Control/types"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	logger := slog.New(slog.NewJSONHandler(willi, &slog.HandlerOptions{
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Stop()
	willi.Close()
	logger.Info("Logging stop!")
}
