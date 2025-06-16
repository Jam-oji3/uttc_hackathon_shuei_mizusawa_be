package usecase

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
	"hackathon/util"
	"net/url"
	"strings"
	"time"
)

type PostCreateUseCase struct {
	TxExecutor repository.TransactionExecutor
	PostRepo   repository.PostsRepository
	DB         *sql.DB
}

func NewPostCreateUseCase(txExecutor repository.TransactionExecutor, postRepo repository.PostsRepository, db *sql.DB) *PostCreateUseCase {
	return &PostCreateUseCase{
		TxExecutor: txExecutor,
		PostRepo:   postRepo,
		DB:         db,
	}
}

func (uc *PostCreateUseCase) Execute(ctx context.Context, userId string, text string, replyTo *string, repostRef *string, mediaType *string, mediaUrl *string) (*model.Post, error) {
	id := util.GenerateULID()
	now := time.Now()

	post := model.Post{
		Id:        id,
		UserId:    userId,
		Text:      text,
		ReplyTo:   replyTo,
		RepostRef: repostRef,
		MediaType: mediaType,
		MediaURL:  mediaUrl,
		CreatedAt: now,
	}

	// バリデーション
	if err := validatePostData(&post); err != nil {
		return nil, err
	}

	// トランザクション内で投稿作成
	_, err := uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.PostRepo.InsertPost(ctx, tx, post); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func validatePostData(post *model.Post) error {
	if post.Id == "" {
		return errors.New("post ID is required")
	}
	if post.UserId == "" {
		return errors.New("user ID is required")
	}
	if strings.TrimSpace(post.Text) == "" && post.MediaURL == nil {
		return errors.New("either text or media URL must be provided")
	}
	if len(post.Text) > 200 {
		return errors.New("text is too long (max 200 characters)")
	}

	if post.RepostRef != nil && post.ReplyTo != nil {
		return errors.New("repostRef and replyTo are mutually exclusive")
	}

	// MediaURLがあればURLとして妥当かチェック
	if post.MediaURL != nil {
		parsedURL, err := url.ParseRequestURI(*post.MediaURL)
		if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
			return errors.New("mediaURL must be a valid HTTP or HTTPS URL")
		}
	}

	// MediaTypeは空でもいいが、もしあれば簡単に制限
	if post.MediaType != nil {
		allowedMediaTypes := map[string]bool{
			"photo": true,
			"model": true,
		}
		if !allowedMediaTypes[*post.MediaType] {
			return errors.New("unsupported media type")
		}
	}

	return nil
}
