package mysql

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
)

type FollowsRepository struct{}

var _ repository.FollowsRepository = (*FollowsRepository)(nil)

func NewFollowsRepository() *FollowsRepository {
	return &FollowsRepository{}
}

func (r *FollowsRepository) InsertFollow(ctx context.Context, dbtx repository.DBTX, follow *model.Follow) error {
	_, err := dbtx.ExecContext(ctx, `
		INSERT INTO follows (id, follower_id, followed_id, created_at)
		VALUES (?, ?, ?, ?)
	`, follow.Id, follow.FollowerId, follow.FollowedId, follow.CreatedAt)
	return err
}

func (r *FollowsRepository) DeleteFollow(ctx context.Context, dbtx repository.DBTX, followerId, followedId string) error {
	_, err := dbtx.ExecContext(ctx, `
		DELETE FROM follows
		WHERE follower_id = ? AND followed_id = ?
	`, followerId, followedId)
	return err
}
