package postgres

import (
	"chat/internal/entity"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *entity.User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID)
}

func (r *UserRepo) FindByUsername(username string) (*entity.User, error) {
	row := r.db.QueryRow(`SELECT id, username, password_hash FROM users WHERE username = $1`, username)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil && err != sql.ErrNoRows {

		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) FindByUserID(userID uint) (*entity.User, error) {
	row := r.db.QueryRow(`SELECT id, username, password_hash FROM users WHERE id = $1`, userID)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
