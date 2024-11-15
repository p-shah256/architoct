package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize app
	app := NewApp()

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = app.templates

	// Routes
	e.GET("/", app.handleHome)
	e.GET("/story", app.handleGetStory)
	e.POST("/vote/:id", app.handleVote)

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}

func NewApp() *App {
	return &App{
		store:     NewStoryStore(),
		templates: newTemplates(),
	}
}
