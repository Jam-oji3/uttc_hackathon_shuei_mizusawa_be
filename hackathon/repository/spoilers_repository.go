package repository

import (
	"context"
	"hackathon/model"
)

type SpoilersRepository interface {
	InsertSpoiler(ctx context.Context, dbtx DBTX, spoiler *model.Spoiler) error
}
