package storage

import (
	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		ByUsername(string) (User, error)
		ExistsByUsername(string) (bool, error)
		ExistsByEmail(string) (bool, error)
		ExistsByID(int) (bool, error)
		Create(user User) error

		LatestPostID(int) int
	}

	Users struct {
		*sqlx.DB
	}

	User struct {
		ID                int    `sq:"id" json:"-"`
		Email             string `sq:"email" json:"email"`
		Username          string `sq:"username" json:"username"`
		FirstName         string `sq:"first_name" json:"first_name"`
		LastName          string `sq:"last_name" json:"last_name"`
		EncryptedPassword string `sq:"password" json:"-" db:"password"`
	}
)

func (db *Users) ByUsername(username string) (user User, err error) {
	const q = "SELECT * FROM users WHERE username = $1"
	err = db.Get(&user, q, username)
	return user, err
}

func (db *Users) ExistsByUsername(username string) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM users WHERE username = $1)"
	row := db.QueryRow(q, username)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

func (db *Users) ExistsByEmail(email string) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM users WHERE email = $1)"
	row := db.QueryRow(q, email)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

func (db *Users) ExistsByID(id int) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM users WHERE id = $1)"
	row := db.QueryRow(q, id)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

func (db *Users) LatestPostID(UserID int) int {
	const q = "SELECT post_id FROM posts WHERE owner_id = $1 ORDER BY post_id DESC LIMIT 1"

	row := db.QueryRow(q, UserID)

	var latest int
	row.Scan(&latest)

	return latest

}

func (db *Users) Create(user User) error {
	const q = "INSERT INTO users (email, username, first_name, last_name, password) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Exec(q, user.Email, user.Username, user.FirstName, user.LastName, user.EncryptedPassword)

	return err

}
