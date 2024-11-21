package handlers

// NOTE: DESIGN: entry point into our service
// these are very specifically HTMX handlers
// if we need pure api we need other handler

import (
	// "log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"architoct/internal/service"
)

type PageData struct {
	PageType string
	Data     interface{}
}

const (
	PageTypeHome    = "home"
	PageTypeStory   = "story"
	PageTypeProfile = "profile"
)

type htmxHandler struct {
	Templates *TemplateRenderer
	service   *service.ArchitoctService
}

func NewHtmxHandler(s *service.ArchitoctService) *htmxHandler {
	tmpl := NewTemplates()
	return &htmxHandler{
		Templates: tmpl,
		service:   s,
	}
}

func (app *htmxHandler) SetupRoutes(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Renderer = app.Templates

	e.GET("/", app.handleGetHome)
	e.GET("/story/:id", app.handleGetStory)

	// e.POST("/story/user/:userid", app.postStory)

	e.POST("/upvote/story/:storyid/user/:userid", app.handlePostSVote)
	e.POST("/upvote/comment/:commentid/user/:userid", app.handlePostCVote)

	e.POST("/comment/story/:storyid/user/:userid", app.handlePostScomment)
	e.POST("/comment/comment/:commentid/user/:userid", app.handlePostCcomment)
}

// GET HANDLERS ///////////////////////////////////////////////////////////////
func (app *htmxHandler) handleGetHome(c echo.Context) error {
	stories, err := app.service.GetHomeFeed(c.Request().Context())
	if err != nil {
		return err
	}
	return c.Render(200, "baseLayout", PageData{
		PageType: PageTypeHome,
		Data:     stories,
	})
}

func (app *htmxHandler) handleGetStory(c echo.Context) error {
	story, err := app.service.GetStoryPage(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(200, "storyContent", story)
	} else {
		return c.Render(200, "baseLayout", PageData{
			PageType: PageTypeStory,
			Data:     story,
		})
	}
}

// POST HANDLERS //////////////////////////////////////////////////////////////
func (app *htmxHandler) handlePostSVote(c echo.Context) error {
	updatedStory, err := app.service.Upvote(c.Request().Context(), service.TypeStory, c.Param("storyid"), c.Param("userid"))
	if err != nil {
		return err
	}
	return c.Render(200, "storyUpvoteMarker", updatedStory)
}

func (app *htmxHandler) handlePostCVote(c echo.Context) error {
	updatedComment, err := app.service.Upvote(c.Request().Context(), service.TypeComment, c.Param("commentid"), c.Param("userid"))
	if err != nil {
		return err
	}
	return c.Render(200, "commentUpvoteMarker", updatedComment)
}

// not returning an HTMX here cause let them refresh
// maybe for story we could just return sucess message or something
func (app *htmxHandler) handlePostScomment(c echo.Context) error {
	return app.service.Comment(c.Request().Context(), c.Param("storyid"), c.Param("userid"), c.FormValue("body"), service.TypeStory)
}

func (app *htmxHandler) handlePostCcomment(c echo.Context) error {
	return app.service.Comment(c.Request().Context(), c.Param("commentid"), c.Param("userid"), c.FormValue("body"), service.TypeComment)
}

func (app *htmxHandler) handlePostStory(c echo.Context) error {
	return app.service.NewStory(c.Request().Context(), c.Param("userid"), c.FormValue("body"), c.FormValue("title"))
}
