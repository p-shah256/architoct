package main

// DESIGN: main just puts things together and starts... THAT's literally it.
import (
	"context"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"

	"architoct/internal/handlers"
	"architoct/internal/logger"
	"architoct/internal/service"
	"architoct/internal/store/mongos"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	if err := logger.SetupLogger("./logs"); err != nil {
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
    e.Use(logger.Middleware)

	e.Renderer = handler.Templates
	// e.Static("/views", "views")
	e.Static("/assets", "views/assets")
	handler.SetupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}
