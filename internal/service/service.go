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
	"github.com/matoous/go-nanoid/v2"
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
func (architcot *ArchitoctService) GetHomeFeed(ctx context.Context, page int64, userid string) ([]types.Story, error) {
	stories, err := architcot.storyStore.GetRecent(ctx, 20, page)
	if err != nil {
		return nil, err
	}

	for i := range stories {
		stories[i].SetUserSpecificData(userid)
		formatStory(&stories[i])
	}
	return stories, nil
}

func (architcot *ArchitoctService) GetCommentReplies(ctx context.Context, commentid string, userid string) ([]types.Comment, error) {
	requestedComment, err := architcot.commentStore.GetById(ctx, commentid)
	if err != nil {
		return nil, types.ErrCommentNotFound
	}

    comments := make([]types.Comment, 0, len(requestedComment.Replies))
	for i := range requestedComment.Replies {
        comment, err := architcot.commentStore.GetById(ctx, requestedComment.Replies[i])
        if err != nil { // if one comment fails, log and igonre
			logger.L.Error().
				Err(err).
				Caller().
				Str("fetching comment", comment.Replies[i]).
				Str("commentid", commentid).
				Msg("failed to fetch comment")
            continue
        }
		comment.SetUserSpecificData(userid)
        formatComment(comment)
        comments = append(comments, *comment)
	}

	logger.Debug().Any("returning comments", comments).Msg("From getcommentreplies service")
	return comments, nil
}

func (architect *ArchitoctService) GetStoryPage(ctx context.Context, id string, userid string) (types.StoryPage, error) {
    requestedStory, err := architect.storyStore.GetByID(ctx, id)
    if err != nil {
        return types.StoryPage{}, err
    }
    comments := make([]types.Comment, 0, len(requestedStory.Replies))
    for i := range requestedStory.Replies {
        comment, err := architect.commentStore.GetById(ctx, requestedStory.Replies[i])
        if err != nil { // if one comment fails, log and igonre
			logger.L.Error().
				Err(err).
				Caller().
				Str("fetching comment", requestedStory.Replies[i]).
				Str("story", id).
				Msg("failed to fetch comment")
            continue
        }
        formatComment(comment)
		comment.SetUserSpecificData(userid)
        comments = append(comments, *comment)
    }
	requestedStory.SetUserSpecificData(userid)
    formatStory(requestedStory)
	logger.Debug().Any("returning story", *requestedStory).Msg("From getstorypage service")
	logger.Debug().Any("returning comments", comments).Msg("From getstorypage service")
    return types.StoryPage{
        Story:    *requestedStory,
        Comments: comments,
    }, nil
}

// POST ////////////////////////////////////////////////////////////////////////
func (architcot *ArchitoctService) Upvote(ctx context.Context, contentType ContentType, id string, userid string) (any, error) {
	logger.Debug().Any("comment", contentType).Str("id", id).Str("userid", userid).Msg("upvoting")
	if userid == "" {
		return nil, types.ErrUserNotFound
	}
	if contentType == TypeComment {
		updatedComment, err := architcot.commentStore.ToggleUpvote(ctx, id, userid)
		updatedComment.SetUserSpecificData(userid)
		return updatedComment, err
	} else {
		updatedStory, err := architcot.storyStore.ToggleUpvote(ctx, id, userid)
		updatedStory.SetUserSpecificData(userid)
		return updatedStory, err
	}
}

// TODO: can we optimise this ? options:
// 1. maybe send a postid with the request.. but that makes the api incosistent
// 2. let db layer handle adding postid to the comment... honestly the postid is not even required here
// how expensive is this extra visit to DB? in term of latency claude says ~0.1-1ms for local
func (architcot *ArchitoctService) Comment(ctx context.Context, parentid string, userid string, body string, contentType ContentType) (*types.Comment, error) {
	var commentID string
	var err error  // Declare err here since it's used throughout
	if userid == "" {
		return nil, types.ErrUserNotFound
	}
	if body == "" {
		return nil, types.ErrCommentNotPosted
	}
	// 1. get parentID
	// 2. create a comment with postid
	// 3. add new comment id to replies array of parent
	if contentType == TypeComment {
		logger.Debug().Msg("commenting on a comment")
		parentComment, err := architcot.commentStore.GetById(ctx, parentid)
		if err != nil {
			return nil, types.ErrCommentNotFound
		}
		logger.Debug().Str("got comment id", parentid).Msg("commenting on a comment")
		commentID, err = architcot.commentStore.InsertToDB(ctx, &types.Comment{
			PostID:    parentComment.PostID,
			Body:      body,
			UserID:    userid,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, types.ErrCommentNotPosted
		}
		logger.Debug().Str("new comment created", commentID).Msg("commenting on a comment")
		err = architcot.commentStore.AddToRepliesArray(ctx, parentid, commentID)
		if err != nil {
			return nil, types.ErrCommentNotPosted
		}
		err = architcot.storyStore.AddCommentCount(ctx, parentid)
		if err != nil {
			return nil, types.ErrCommentNotPosted
		}
	} else {
		commentID, err = architcot.commentStore.InsertToDB(ctx, &types.Comment{
			PostID:    parentid,
			Body:      body,
			UserID:    userid,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
		logger.Debug().Str("adding comment to post", parentid).Msg("service.. ading to post")
		err = architcot.storyStore.AddToRepliesArray(ctx, parentid, commentID)
if err != nil {
			return nil, types.ErrCommentNotPosted
		}
		 err = architcot.storyStore.AddCommentCount(ctx, parentid)
if err != nil {
			return nil, types.ErrCommentNotPosted
		}

	}
	logger.Debug().Str("got commentID", commentID).Msg("at service adding comment")
	newComment, err := architcot.commentStore.GetById(ctx, commentID)
	return newComment, err
}

func (architcot *ArchitoctService) NewStory(ctx context.Context, userid string, body string, title string) (*types.Story, error) {
    story, err := architcot.storyStore.Create(ctx, &types.Story{
        ID: gonanoid.Must(4),
        CreatedAt: time.Now(),
        Body: body,
        UserID: userid,
        Title: title,
    })
    return story, err
}

func (architcot *ArchitoctService) User(ctx context.Context, userid string) error {
	logger.Debug().Str("userid", userid).Msg("serviceUser")
	_,  err := architcot.userStore.Create(ctx, userid)
	return err
}
