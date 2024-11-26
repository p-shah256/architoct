package logger

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

var L zerolog.Logger

func SetupLogger(logPath string) error {
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if err := os.MkdirAll(logPath, 0755); err != nil {
        return fmt.Errorf("failed to create log directory: %v", err)
    }

    consoleWriter := zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: "15:04:05",
        NoColor:    false,
    }

    logFile, err := os.OpenFile(
        filepath.Join(logPath, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))),
        os.O_CREATE|os.O_APPEND|os.O_WRONLY,
        0644,
    )
    if err != nil {
        return fmt.Errorf("failed to open log file: %v", err)
    }

    multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

    L = zerolog.New(multi).With().Timestamp().Logger()

    log.Logger = L

    return nil
}

func Debug() *zerolog.Event {
    return L.Debug().Caller(1)
}
