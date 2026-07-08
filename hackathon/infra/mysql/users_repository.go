package mysql

import (
	"context"
	"database/sql"
	"errors"
	"hackathon/model"
	"hackathon/repository"
)

type UsersRepository struct {
}

// UsersRepositoryインターフェースを実装
var _ repository.UsersRepository = (*UsersRepository)(nil)

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{}
}

func (r *UsersRepository) FindById(ctx context.Context, dbtx repository.DBTX, id string) (*model.User, error) {
	row := dbtx.QueryRowContext(ctx, `
	SELECT id, username, display_name, email, bio, icon_url, created_at, updated_at
	FROM users WHERE id = ?`, id)

	var u model.User
	err := row.Scan(&u.Id, &u.Username, &u.DisplayName, &u.Email, &u.Bio, &u.IconURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) FindByUserName(ctx context.Context, dbtx repository.DBTX, userName string) (*model.User, error) {
	row := dbtx.QueryRowContext(ctx, `
		SELECT id, username, display_name, email, bio, icon_url, created_at, updated_at
	FROM users WHERE username = ?`, userName)

	var u model.User
	err := row.Scan(&u.Id, &u.Username, &u.DisplayName, &u.Email, &u.Bio, &u.IconURL, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) FindProfileByUsername(ctx context.Context, dbtx repository.DBTX, username string, viewerId string) (*model.UserProfile, error) {
	row := dbtx.QueryRowContext(ctx, `
	SELECT 
	    u.id,
		u.username,
		u.display_name,
		u.bio,
		u.icon_url,
		u.created_at,
		(SELECT COUNT(*) FROM follows WHERE follower_id = u.id) AS following_count,
		(SELECT COUNT(*) FROM follows WHERE followed_id = u.id) AS follower_count,
		(SELECT COUNT(*) FROM posts WHERE user_id = u.id) AS post_count,
		EXISTS (
			SELECT 1 FROM follows WHERE follower_id = ? AND followed_id = u.id
		) AS is_following
	FROM users u
	WHERE u.username = ?
	`, viewerId, username) // viewerIdが?1, usernameが?2の順番

	var profile model.UserProfile
	var isFollowing bool

	err := row.Scan(
		&profile.Id,
		&profile.Username,
		&profile.DisplayName,
		&profile.Bio,
		&profile.IconURL,
		&profile.CreatedAt,
		&profile.Stats.FollowingCount,
		&profile.Stats.FollowerCount,
		&profile.Stats.PostCount,
		&isFollowing,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	profile.IsFollowing = isFollowing // UserProfileにisFollowingのboolフィールドを追加しておく

	return &profile, nil
}

func (r *UsersRepository) Insert(ctx context.Context, dbtx repository.DBTX, user *model.User) error {
	_, err := dbtx.ExecContext(ctx, `
	INSERT INTO users (id, username, display_name, email, bio, icon_url, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Id, user.Username, user.DisplayName, user.Email, user.Bio, user.IconURL, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UsersRepository) Update(ctx context.Context, dbtx repository.DBTX, user *model.User) error {
	res, err := dbtx.ExecContext(ctx, `
	UPDATE users 
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

func (r *UsersRepository) Delete(ctx context.Context, dbtx repository.DBTX, id string) error {
	res, err := dbtx.ExecContext(ctx, `
	DELETE FROM users WHERE id = ?`, id)
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
