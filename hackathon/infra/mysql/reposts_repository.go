package mysql

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
)

type RepostsRepository struct{}

var _ repository.RepostsRepository = (*RepostsRepository)(nil)

func NewRepostsRepository() *RepostsRepository {
	return &RepostsRepository{}
}

func (r *RepostsRepository) InsertRepost(ctx context.Context, dbtx repository.DBTX, repost *model.Repost) error {
	_, err := dbtx.ExecContext(ctx, `
		INSERT INTO reposts (id, user_id, post_id, created_at)
		VALUES (?, ?, ?, ?)
	`, repost.Id, repost.UserId, repost.PostId, repost.CreatedAt)
	return err
}

func (r *RepostsRepository) DeleteRepost(ctx context.Context, dbtx repository.DBTX, userId string, postId string) error {
	_, err := dbtx.ExecContext(ctx, `
		DELETE FROM reposts
		WHERE user_id = ? AND post_id = ?
	`, userId, postId)
	return err
}
