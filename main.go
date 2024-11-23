package main

// DESIGN: main just puts things together and starts... THAT's literally it.
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"architoct/internal/handlers"
	"architoct/internal/service"
	"architoct/internal/store/mongos"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupLogger(logPath string) error {
    // Create log directory if it doesn't exist
    if err := os.MkdirAll(logPath, 0755); err != nil {
        return fmt.Errorf("failed to create log directory: %v", err)
    }
    consoleWriter := zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: "15:04:05",
        NoColor:    false,
    }

    // File output (for backup)
    logFile, err := os.OpenFile(
        filepath.Join(logPath, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))),
        os.O_CREATE|os.O_APPEND|os.O_WRONLY,
        0644,
    )
    if err != nil {
        return fmt.Errorf("failed to open log file: %v", err)
    }

    // Multi-writer setup
    multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
    log.Logger = zerolog.New(multi).With().Timestamp().Logger()

    return nil
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()     // 1. Record start time
        req := c.Request()      // 2. Get request info
        res := c.Response()     // 3. Get response object
        err := next(c)          // 4. Call next middleware/handler

        // 5. After handler finishes, do logging
        if err == nil && res.Status < 400 {
            // Minimal logging for successful requests
			log.Info().                    // 1. Create log entry
				Int("status", res.Status). // 2. Add field
				Str("method", req.Method). // 3. Add another field
                Str("path", req.URL.Path).
                Dur("latency", time.Since(start)).
                Send()
        } else {
            // Detailed logging for errors
            errLogger := log.Error()
            if err != nil {
                errLogger = errLogger.Err(err)
            }

            errLogger.
                Int("status", res.Status).
                Str("method", req.Method).
                Str("path", req.URL.Path).
                Str("user_agent", req.UserAgent()).
                Str("remote_ip", c.RealIP()).
                Dur("latency", time.Since(start)).
                Interface("headers", req.Header).
                Send()
        }

        return err
    }
}

func main() {
	// CONNECT TO DB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://admin:password123@localhost:27017"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		// slog.Error("Cannot connect to DB", "error", err)
		return
	}
	defer client.Disconnect(context.Background())

	if err := SetupLogger("./logs"); err != nil {
		fmt.Printf("Error setting up logger: %v\n", err)
		os.Exit(1)
	}

	db := client.Database("mvp_db")

	// init stores
	c := mongos.NewCommentStore(db)
	s := mongos.NewStoryStore(db)
	u := mongos.NewUserStore(db)

	// init service
	// DESIGN: abstract this away so any dbClient can be used in service
	service := service.NewArchitoctService(s, c, u, client)
	// Initialize handler
	handler := handlers.NewHtmxHandler(service)

	e := echo.New()

	// Middleware for logging both info and errors
    e.Use(LoggingMiddleware)

	e.Renderer = handler.Templates
	// e.Static("/views", "views")
	e.Static("/assets", "views/assets")
	handler.SetupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}
