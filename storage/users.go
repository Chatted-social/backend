package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		ByUsername(string) (User, error)
		ExistsByUsername(string) (bool, error)
		Create(user User) (int, error)
	}

	Users struct {
		*sqlx.DB
	}

	User struct {
		CreatedAt         time.Time `sq:"created_at" json:"created_at"`
		UpdatedAt         time.Time `sq:"updated_at" json:"updated_at"`
		ID                int       `sq:"id" json:"-"`
		Email             string    `sq:"email" json:"email"`
		Username          string    `sq:"username" json:"username"`
		FirstName         string    `sq:"first_name" json:"first_name"`
		LastName          string    `sq:"last_name" json:"last_name"`
		EncryptedPassword string    `sq:"password" json:"-" db:"password"`
	}
)

func (db *Users) ByUsername(username string) (user User, err error) {
	const q = "SELECT * FROM users WHERE username = $1"
	err = db.Get(&user, q, username)
	return user, err
}

func (db *Users) ExistsByUsername(username string) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM USERS WHERE username = $1)"
	row := db.QueryRow(q, username)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

func (db *Users) Create(user User) (int, error) {
	const q = `INSERT INTO users (email, username, first_name, last_name, password) VALUES
    ($1, $2, $3, $4, $5) RETURNING id`

	rows, err := db.Query(q, user.Email, user.Username, user.FirstName, user.LastName, user.EncryptedPassword)
	if err != nil {
		return 0, err
	}

	var id int
	rows.Next()

	err = rows.Scan(&id)

	return 0, err
}
