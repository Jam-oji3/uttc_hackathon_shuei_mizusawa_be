package mysql

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
)

type UserRepository struct {
	DB *sql.DB
}

// UserRepositoryインターフェースを実装
var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*model.User, error) {
	row := r.DB.QueryRowContext(ctx, `
	SELECT id, username, display_name, email, bio, icon_url, created_at, updated_at
	FROM user WHERE id = ?`, id)

	var u model.User
	err := row.Scan(&u.Id, &u.UserName, &u.DisplayName, &u.Email, &u.Bio, &u.IconURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByUserName(ctx context.Context, userName string) (*model.User, error) {
	row := r.DB.QueryRowContext(ctx, `
	SELECT id, username, display_name, email, bio, icon_url, created_at, updated_at
	FROM user WHERE username = ?`, userName)

	var u model.User
	err := row.Scan(&u.Id, &u.UserName, &u.DisplayName, &u.Email, &u.Bio, &u.IconURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Insert(ctx context.Context, tx *sql.Tx, user *model.User) error {
	_, err := tx.ExecContext(ctx, `
	INSERT INTO user (id, username, display_name, email, bio, icon_url, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Id, user.UserName, user.DisplayName, user.Email, user.Bio, user.IconURL, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) Update(ctx context.Context, tx *sql.Tx, user *model.User) error {
	res, err := tx.ExecContext(ctx, `
	UPDATE user 
	SET display_name = ?, bio = ?, icon_url = ?, updated_at = ?
	WHERE id = ?`,
		user.DisplayName, user.Bio, user.IconURL, user.UpdatedAt, user.Id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	res, err := tx.ExecContext(ctx, `
	DELETE FROM user WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
