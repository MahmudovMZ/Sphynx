package data

import (
	"wordGame/internal/models"
	db "wordGame/pkg"
)

func SignUp_user(username, password string) error {
	_, err := db.GetDB().Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, password)
	return err
}

func GetByUsersName(username string) ([]models.User, error) {
	rows, err := db.GetDB().Query(
		`SELECT id, username, password_hash, created_at FROM users WHERE username = $1`,
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.PasswordHash,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
