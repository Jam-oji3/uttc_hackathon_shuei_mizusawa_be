package mysql

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
)

type UserRepository struct {
}

// UserRepositoryインターフェースを実装
var _ repository.UserRepository = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindById(ctx context.Context, dbtx repository.DBTX, id string) (*model.User, error) {
	row := dbtx.QueryRowContext(ctx, `
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

func (r *UserRepository) FindByUserName(ctx context.Context, dbtx repository.DBTX, userName string) (*model.User, error) {
	row := dbtx.QueryRowContext(ctx, `
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

func (r *UserRepository) Insert(ctx context.Context, dbtx repository.DBTX, user *model.User) error {
	_, err := dbtx.ExecContext(ctx, `
	INSERT INTO user (id, username, display_name, email, bio, icon_url, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Id, user.UserName, user.DisplayName, user.Email, user.Bio, user.IconURL, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) Update(ctx context.Context, dbtx repository.DBTX, user *model.User) error {
	res, err := dbtx.ExecContext(ctx, `
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

func (r *UserRepository) Delete(ctx context.Context, dbtx repository.DBTX, id string) error {
	res, err := dbtx.ExecContext(ctx, `
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
