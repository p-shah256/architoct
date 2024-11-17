package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize app
	app := NewApp()

	// TODO: eventually use hexagonal arch
	// 1. create store
	// 2. create a service that works with that store ()
	// 3. use them here
	//
	// Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = app.templates

	// Routes
	e.GET("/", app.handleHome)
	e.GET("/story/:id", app.handleGetStory)
	e.POST("/vote/:id", app.handleVote)

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}
