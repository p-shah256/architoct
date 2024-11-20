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

    e.POST("/upvote/story/:storyid/user/:userid", app.handleSVote)
    e.POST("/comment/story/:storyid/user/:userid", app.handleComment)
    e.POST("/upvote/comment/:commentid/user/:userid", app.handleCVote)
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
    slog.Info("story comments", "number", len(story.Comments))
    return c.Render(200, "storyPage", story)
}

// POST HANDLERS //////////////////////////////////////////////////////////////
func (app *htmxHandler) handleSVote(c echo.Context) error {
    updatedStory, err := app.service.Upvote(c.Request().Context(), service.TypeStory, c.Param("storyid"), c.Param("userid"))
    if err != nil {
        return err
    }
    return c.Render(200, "storyUpvoteMarker", updatedStory)
}

func (app *htmxHandler) handleCVote(c echo.Context) error {
    updatedComment, err := app.service.Upvote(c.Request().Context(), service.TypeComment, c.Param("commentid"), c.Param("userid"))
    if err != nil {
        return err
    }
    return c.Render(200, "commentUpvoteMarker", updatedComment)
}

func (app *htmxHandler) handleComment(c echo.Context) error {
    body := c.FormValue("body")
    storyID := c.Param("storyid")
    userID := c.Param("userid")

    err := app.service.Comment(c.Request().Context(), storyID, userID, body, service.TypeStory)
    if err != nil {
        return err
    }
    return nil
}
