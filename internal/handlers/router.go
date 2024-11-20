package handlers

// NOTE: DESIGN: entry point into our service
// these are very specifically HTMX handlers
// if we need pure api we need other handler

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"architoct/internal/service"
)

type htmxHandler struct {
    Templates *TemplateRenderer
    service *service.ArchitoctService
}

func NewHtmxHandler(s *service.ArchitoctService) *htmxHandler {
    tmpl := NewTemplates()
    return &htmxHandler{
        Templates: tmpl,
        service: s,
    }
}

func (app *htmxHandler) SetupRoutes(e *echo.Echo) {
    e.Use(middleware.Logger())
    e.Renderer = app.Templates

    e.GET("/", app.handleHome)
    e.GET("/story/:id", app.handleStory)
    // e.POST("/story", app.postStory)

    e.POST("/upvote/story/:userid/:id", app.handleSVote)
    // e.POST("/upvote/comment/:id", app.handleCvote)
}

// Story handlers
func (app *htmxHandler) handleHome(c echo.Context) error {
    stories, err := app.service.HomeFeed(c.Request().Context())
    if err != nil {
        return err
    }
    return c.Render(200, "baseLayout", stories)
}

func (app *htmxHandler) handleStory(c echo.Context) error {
    story, err := app.service.StoryPage(c.Request().Context(), c.Param("id"))
    if err != nil {
        return err
    }
    slog.Info("returning story", "story", story)
    return c.Render(200, "storyPage", story)
}

func (app *htmxHandler) handleSVote(c echo.Context) error {
    updatedStory, err := app.service.Upvote(c.Request().Context(), false, c.Param("id"), c.Param("userid"))
    if err != nil {
        return err
    }
    return c.Render(200, "upvoteMarker", updatedStory)
}
