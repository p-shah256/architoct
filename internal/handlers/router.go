package handlers

// NOTE: DESIGN: entry point into our service
// these are very specifically HTMX handlers
// if we need pure api we need other handler

import (
	"fmt"
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
    e.POST("/comment/story/:storyid/user/:userid", app.handleComment)
    // e.POST("/upvote/comment/:id", app.handleCvote)
}

// GET HANDLERS ///////////////////////////////////////////////////////////////
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
    slog.Info("story comments", "number", story)
    return c.Render(200, "storyPage", story)
}

// POST HANDLERS //////////////////////////////////////////////////////////////
func (app *htmxHandler) handleSVote(c echo.Context) error {
    updatedStory, err := app.service.Upvote(c.Request().Context(), service.TypeStory, c.Param("id"), c.Param("userid"))
    if err != nil {
        return err
    }
    return c.Render(200, "upvoteMarker", updatedStory)
}

func (app *htmxHandler) handleComment(c echo.Context) error {
    body := c.FormValue("body")
    storyID := c.Param("storyid")
    userID := c.Param("userid")
    fmt.Printf("Received at router: body=%s, storyID=%s, userID=%s\n", body, storyID, userID)

    err := app.service.Comment(c.Request().Context(), storyID, userID, body, service.TypeStory)
    if err != nil {
        return err
    }
    return nil
}
