package handlers

import (
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
    // e.POST("/story", app.postStory)

    e.POST("/upvote/story/:id", app.handleSVote)
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

func (app *htmxHandler) handleSVote(c echo.Context) error {
    updatedStory, err := app.service.Upvote(c.Request().Context(), false, c.Param("id"))
    if err != nil {
        return err
    }
    return c.Render(200, "upvoteMarker", updatedStory)
}

// func (app *htmxHandler) handleCvote(c echo.Context) error {
//     status, newCount, err := app.service.Upvote(c.Request().Context(), true, c.Param("id"))
//     if err != nil {
//         return err
//     }
//     return c.Render(200, "baseLayout", stories)
// }


// func (app *htmxHandler) postStory(c echo.Context) error {
//     story, err := app.story.Create(
//         c.Request().Context(),
//         c.FormValue("title"),
//         c.FormValue("body"),
//     )
//     if err != nil {
//         return c.Render(400, "partials/error", err)
//     }
//     return c.Render(200, "partials/story", story)
// }

// func (app *htmxHandler) handleVote(c echo.Context) error {
//     err := app.story.Upvote(
//         c.Request().Context(),
//         c.Param("id"),
//     )
//     if err != nil {
//         return c.Render(400, "partials/error", err)
//     }
//     // Return just the updated vote count
//     count, err := app.story.GetVoteCount(c.Request().Context(), c.Param("id"))
//     if err != nil {
//         return err
//     }
//     return c.Render(200, "partials/vote-count", count)
// }
