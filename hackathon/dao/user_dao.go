package dao

import (
	"database/sql"
	"errors"
	"hackathon/model"
)

type UserDAO struct {
	DB *sql.DB
}

func (dao *UserDAO) FindById(id string) (*model.User, error) {
	row := dao.DB.QueryRow(`
	SELECT id, username, display_name, email, bio, icon_url, created_at, updated_at
	FROM user
	WHERE id = ?`, id)

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

func (dao *UserDAO) FindByUserName(userName string) (*model.User, error) {
	row := dao.DB.QueryRow(`
	SELECT id, username, display_name, email, bio, icon_url, created_at, updated_at
	FROM user
	WHERE username = ?`, userName)

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

// 変更処理はトランザクション対応
func (dao *UserDAO) Insert(tx *sql.Tx, user *model.User) error {
	_, err := tx.Exec(`
		INSERT INTO user (id, username, display_name, email, bio, icon_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Id, user.UserName, user.DisplayName, user.Email, user.Bio, user.IconURL, user.CreatedAt, user.UpdatedAt)
	return err
}

func (dao *UserDAO) Update(tx *sql.Tx, user *model.User) error {
	res, err := tx.Exec(`
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

func (dao *UserDAO) Delete(tx *sql.Tx, id string) error {
	res, err := tx.Exec(`
	DELETE FROM user 
	WHERE id = ?`, id)
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
