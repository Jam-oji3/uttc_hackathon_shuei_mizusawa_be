package dao

import (
	"database/sql"
	"db/model"
)

type UserDAO struct {
	DB *sql.DB
}

func (dao *UserDAO) FindByName(name string) ([]model.User, error) {
	rows, err := dao.DB.Query("SELECT id, name, age FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (dao *UserDAO) Insert(user model.User) error {
	_, err := dao.DB.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", user.Id, user.Name, user.Age)
	return err
}
