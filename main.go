package main

// DESIGN: main just puts things together and starts... THAT's literally it.
import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"

	"architoct/internal/handlers"
	"architoct/internal/service"
	"architoct/internal/store/mongos"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// CONNECT TO DB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:password123@localhost:27017"))
	if err != nil {
		slog.Error("Cannot connect to DB", "error", err)
		return
	}
	// Ensure the client closes when the function returns
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			slog.Error("Failed to disconnect", "error", err)
		}
	}()

	db := client.Database("mvp_db")

	// init stores
	c := mongos.NewCommentStore(db)
	s := mongos.NewStoryStore(db)
	u := mongos.NewUserStore(db)

	// init service
	service := service.NewArchitoctService(s, c, u)

	// Initialize handler
	handler := handlers.NewHtmxHandler(service)
	e := echo.New()
	e.Renderer = handler.Templates
	e.Static("/views", "views")
	handler.SetupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}
