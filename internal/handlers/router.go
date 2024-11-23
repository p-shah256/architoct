package handlers

// NOTE: DESIGN: entry point into our service
// these are very specifically HTMX handlers
// if we need pure api we need other handler

import (
	// "log/slog"

	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"

	"architoct/internal/service"
	"architoct/internal/types"
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
	e.Renderer = app.Templates

	// home/get ops
	e.GET("/", app.handleGetHome)
	e.GET("/story/:id", app.handleGetStory)

	// user ops
	e.POST("/user/init", app.handleUser)
	e.POST("/user/name", app.handleUser)

	// story ops
	e.POST("/upvote/story/:storyid/user/:userid", app.handlePostSVote)
	e.POST("/comment/story/:storyid/user/:userid", app.handlePostScomment)

	// comment ops
	e.POST("/comment/comment/:commentid/user/:userid", app.handlePostCcomment)
	e.POST("/upvote/comment/:commentid/user/:userid", app.handlePostCVote)
	e.GET("/replies/comment/:commentid", app.handleGetCreplies)
}

// TODO: verify if infinite scroll works fine
// GET HANDLERS ///////////////////////////////////////////////////////////////
func (app *htmxHandler) handleGetHome(c echo.Context) error {
	page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}
	stories, err := app.service.GetHomeFeed(c.Request().Context(), page)
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

// TODO: implement this
func (app *htmxHandler) handleGetCreplies(c echo.Context) error {
	return nil
}

// POST HANDLERS //////////////////////////////////////////////////////////////
func (app *htmxHandler) handlePostSVote(c echo.Context) error {
	updatedStory, err := app.service.Upvote(c.Request().Context(), service.TypeStory, c.Param("storyid"), c.Param("userid"))
	if err != nil {
		return err
	}
	return c.Render(200, "singleStoryBlock", updatedStory)
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
	slog.Info("got request for a comment on a comment")
	return app.service.Comment(c.Request().Context(), c.Param("commentid"), c.Param("userid"), c.FormValue("body"), service.TypeComment)
}

func (app *htmxHandler) handlePostStory(c echo.Context) error {
	return app.service.NewStory(c.Request().Context(), c.Param("userid"), c.FormValue("body"), c.FormValue("title"))
}

func (app *htmxHandler) handleUser(c echo.Context) error{
	userId := c.Request().Header.Get("X-User-ID")
	// if not present it will be empty string .. indicating it is a create
	// if present indicates update username
    userName := c.Request().Header.Get("X-Display-Name")
	err := app.service.User(c.Request().Context(), userId, userName)
	if err != nil {
		switch err{
		// TODO: handle this send a message
		case types.ErrUsernameTaken:
			return nil
		}
	}
	return nil
}
