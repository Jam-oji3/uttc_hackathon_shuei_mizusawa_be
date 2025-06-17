package mysql

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
)

type LikesRepository struct{}

var _ repository.LikesRepository = (*LikesRepository)(nil)

func NewLikesRepository() *LikesRepository {
	return &LikesRepository{}
}

func (r *LikesRepository) InsertLike(ctx context.Context, dbtx repository.DBTX, like *model.Like) error {
	_, err := dbtx.ExecContext(ctx, `
		INSERT INTO likes (id, user_id, post_id, created_at)
		VALUES (?, ?, ?, ?)
	`, like.Id, like.UserId, like.PostId, like.CreatedAt)
	return err
}

func (r *LikesRepository) DeleteLike(ctx context.Context, dbtx repository.DBTX, userId string, postId string) error {
	_, err := dbtx.ExecContext(ctx, `
		DELETE FROM likes
		WHERE user_id = ? AND post_id = ?
	`, userId, postId)
	return err
}
