package service

// DESIGN:
// this layer does not care if the handler was HTMX or REACT or just a cli..
// it does not care about rendering or templates or anything
// also it does not care if its mongo or something else.. so it has to be passed into it
//
// this is the real buisness logic
//
// also this cannot return any ptrs cause HTMX needs actual data
//
// IDEALLY this also should have no idea about the DB but ... its just very complex to do it otherwise to handle transactions
// altho we use nosql if we do not need atomicity but I need it.

import (
	"architoct/internal/logger"
	"architoct/internal/store/mongos"
	"architoct/internal/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// ENUM TO KEEP things simple///////////////////////////////////////////////////
type ContentType string

const (
	TypeStory   ContentType = "STORY"
	TypeComment ContentType = "COMMENT"
)

type ArchitoctService struct {
	storyStore   *mongos.StoryStore
	commentStore *mongos.CommentStore
	userStore    *mongos.UserStore
	dbClient 	*mongo.Client
}

// DESIGN: abstract these stores..
func NewArchitoctService(s *mongos.StoryStore, c *mongos.CommentStore, u *mongos.UserStore, dbClient *mongo.Client) *ArchitoctService {
	return &ArchitoctService{
		storyStore:   s,
		commentStore: c,
		userStore:    u,
		dbClient: dbClient,
	}
}

// PAGE GETS //////////////////////////////////////////////////////////////////////
// TODO: add ability for infinite scrolling
func (architcot *ArchitoctService) GetHomeFeed(ctx context.Context, page int64) ([]types.Story, error) {
	stories, err := architcot.storyStore.GetRecent(ctx, 20, page)
	if err != nil {
		return nil, err
	}

	for i := range stories {
		formatStory(&stories[i])
	}
	return stories, nil
}

func (architect *ArchitoctService) GetStoryPage(ctx context.Context, id string) (types.StoryPage, error) {
    requestedStory, err := architect.storyStore.GetByID(ctx, id)
    if err != nil {
        return types.StoryPage{}, err
    }

    comments := make([]types.Comment, 0, len(requestedStory.Replies))
    for i := range requestedStory.Replies {
        comment, err := architect.commentStore.GetById(ctx, requestedStory.Replies[i])
		// if one comment fails, log and igonre
        if err != nil {
			logger.L.Error().
				Str("error", "failed to fetch comment for storyId").
				Str("story", id).
				Str("comment", requestedStory.Replies[i])
        }
        formatComment(comment)
        comments = append(comments, *comment)
    }

    formatStory(requestedStory)
    return types.StoryPage{
        Story:    *requestedStory,
        Comments: comments,
    }, nil
}

// POST ////////////////////////////////////////////////////////////////////////
func (architcot *ArchitoctService) Upvote(ctx context.Context, contentType ContentType, id string, userid string) (any, error) {
	logger.Debug().Any("comment", contentType).Str("id", id).Str("userid", userid).Msg("upvoting")
	if contentType == TypeComment {
		updatedComment, err := architcot.commentStore.ToggleUpvote(ctx, id, userid)
		return updatedComment, err
	} else {
		updatedStory, err := architcot.storyStore.ToggleUpvote(ctx, id, userid)
		return updatedStory, err
	}
}

// TODO: can we optimise this ? options:
// 1. maybe send a postid with the request.. but that makes the api incosistent
// 2. let db layer handle adding postid to the comment... honestly the postid is not even required here
// how expensive is this extra visit to DB? in term of latency claude says ~0.1-1ms for local
func (architcot *ArchitoctService) Comment(ctx context.Context, parentid string, userid string, body string, contentType ContentType) error {
	if contentType == TypeComment {
		parentComment, err := architcot.commentStore.GetById(ctx, parentid)
		if err != nil {
			return types.ErrCommentNotFound
		}
		commentID, err := architcot.commentStore.Create(ctx, &types.Comment{
			PostID:    parentComment.PostID,
			Body:      body,
			UserID:    userid,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return types.ErrCommentNotPosted
		}
		return architcot.commentStore.AddReply(ctx, parentid, commentID)
	} else {
		commentID, err := architcot.commentStore.Create(ctx, &types.Comment{
			PostID:    parentid,
			Body:      body,
			UserID:    userid,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		return architcot.storyStore.AddComment(ctx, parentid, commentID)
	}
}

func (architcot *ArchitoctService) NewStory(ctx context.Context, userid string, body string, title string) error {
	return architcot.storyStore.Create(ctx, &types.Story{
		CreatedAt: time.Now(),
		Body: body,
		UserID: userid,
		Title: title,
	})
}

func (architcot *ArchitoctService) User(ctx context.Context, userid string) error {
	logger.Debug().Str("userid", userid)
	_,  err := architcot.userStore.Create(ctx, userid)
	return err
}
