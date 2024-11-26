package handlers

// NOTE: DESIGN: entry point into our service
// these are very specifically HTMX handlers
// if we need pure api we need other handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"architoct/internal/logger"
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
	PageTypeAbout = "about"
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
	e.POST("/user", app.handleUser)
	e.GET("/about", app.handleAbout)

	// story ops
	e.POST("/upvote/story/:storyid", app.handlePostSVote)
	e.POST("/story", app.handlePostStory)

	// comment ops
	e.POST("/comment/comment/:commentid", app.handlePostCcomment)
	e.POST("/upvote/comment/:commentid", app.handlePostCVote)
	e.POST("/comment/story/:storyid", app.handlePostScomment)
	e.GET("/comments/replies/:commentid", app.handleGetCreplies)
	e.GET("/load-editor/:commentid", app.handleGetEditor)
}

// TODO: verify if infinite scroll works fine
// GET HANDLERS ///////////////////////////////////////////////////////////////
func (app *htmxHandler) handleGetHome(c echo.Context) error {
	page, _ := strconv.ParseInt(c.QueryParam("page"), 10, 64)
	if page == 0 {
		page = 1
	}
	stories, err := app.service.GetHomeFeed(c.Request().Context(), page, c.Get("userID").(string))
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "baseLayout", PageData{
		PageType: PageTypeHome,
		Data:     stories,
	})
}


func (app *htmxHandler) handleAbout(c echo.Context) error {
if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "aboutContent", nil)
	} else {
		return c.Render(http.StatusOK, "baseLayout", PageData{
			PageType: PageTypeAbout,
			Data:     nil,
		})
	}
}

func (app *htmxHandler) handleGetStory(c echo.Context) error {
	story, err := app.service.GetStoryPage(c.Request().Context(), c.Param("id"), c.Get("userID").(string))
	if err != nil {
		return err
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "storyContent", story)
	} else {
		return c.Render(http.StatusOK, "baseLayout", PageData{
			PageType: PageTypeStory,
			Data:     story,
		})
	}
}

func (app *htmxHandler) handleGetEditor(c echo.Context) error {
    commentID := c.Param("id")
    return c.Render(http.StatusOK, "commentReplyForm", map[string]interface{}{
        "ID": commentID,
    })
}

func (app *htmxHandler) handleGetCreplies(c echo.Context) error {
	comments, err := app.service.GetCommentReplies(c.Request().Context(), c.Param("commentid"), c.Get("userID").(string))
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "comment-replies", comments)
}

// POST HANDLERS //////////////////////////////////////////////////////////////
func (app *htmxHandler) handlePostSVote(c echo.Context) error {
	cookie, err := c.Cookie("userID")
	logger.Debug().Str("userid", cookie.Value).Msg("GetHome")
	userID := cookie.Value
	updatedStory, err := app.service.Upvote(c.Request().Context(), service.TypeStory, c.Param("storyid"), userID)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "singleStoryBlock", updatedStory)
}

func (app *htmxHandler) handlePostCVote(c echo.Context) error {
	updatedComment, err := app.service.Upvote(c.Request().Context(), service.TypeComment, c.Param("commentid"), c.Get("userID").(string))
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "commentUpvoteMarker", updatedComment)
}

// not returning an HTMX here cause let them refresh
// maybe for story we could just return sucess message or something
func (app *htmxHandler) handlePostScomment(c echo.Context) error {
	newComment, err := app.service.Comment(c.Request().Context(), c.Param("storyid"), c.Get("userID").(string), c.FormValue("body"), service.TypeStory)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "singleComment", newComment)
}

func (app *htmxHandler) handlePostCcomment(c echo.Context) error {
	newComment, err := app.service.Comment(c.Request().Context(), c.Param("commentid"), c.Get("userID").(string), c.FormValue("body"), service.TypeComment)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "singleComment", newComment)
}

func (app *htmxHandler) handlePostStory(c echo.Context) error {
	userId := c.Request().Header.Get("X-User-ID")
	body := c.FormValue("body")
	title := c.FormValue("title")

	newStory, err := app.service.NewStory(c.Request().Context(), userId, body, title)
	if err != nil {
		return err
	}
	logger.Debug().Msg("Successfully created new story")

	return c.Render(http.StatusOK, "singleStoryBlock", newStory)
}

func (app *htmxHandler) handleUser(c echo.Context) (error) {
	userId := c.Request().Header.Get("X-User-ID")
	logger.Debug().Str("userID", userId).Msg("handleUser")
	err := app.service.User(c.Request().Context(), userId)
	if err != nil {
		switch err {
		case types.ErrUsernameTaken:
			return types.ErrUsernameTaken
		}
	}
	return nil
}
