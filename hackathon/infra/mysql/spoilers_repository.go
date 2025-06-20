package mysql

import (
	"context"
	"hackathon/model"
	"hackathon/repository"
)

type SpoilersRepository struct{}

var _ repository.SpoilersRepository = (*SpoilersRepository)(nil)

func NewSpoilersRepository() *SpoilersRepository {
	return &SpoilersRepository{}
}

func (r *SpoilersRepository) InsertSpoiler(ctx context.Context, dbtx repository.DBTX, spoiler *model.Spoiler) error {
	_, err := dbtx.ExecContext(ctx, `
		INSERT INTO spoilers (id, post_id, keyword)
		VALUES (?, ?, ?)
	`, spoiler.Id, spoiler.PostId, spoiler.Keyword)
	return err
}
