package main

import (
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (app *App) handleHome(c echo.Context) error {
	data := PageData{
		Stories: app.store.GetAllStories(),
	}
	return c.Render(200, "baseLayout", data)
}

func (app *App) handleGetStory(c echo.Context) error {
	idStr := c.Param("id")  // Changed from QueryParam to Param
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "Invalid story ID")
	}

	story, found := app.store.GetStory(id)
	if !found {
		return c.String(404, "Story not found")
	}

	slog.Info("sending this story", "story", story)
	return c.Render(200, "storyThreadPage", story)
}

func (app *App) handleVote(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "Invalid story ID")
	}

	if ok := app.store.IncrementVote(id); !ok {
		return c.String(404, "Story not found")
	}

	story, _ := app.store.GetStory(id)
	return c.JSON(200, story.VoteCount)
}
