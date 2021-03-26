package storage

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	PostStorage interface {
		Create(p Post) (post Post, err error)
		Update(p Post) (post Post, err error)
		Delete(p Post) error
		UserPosts(id, limit, offset int) ([]Post, error)
		ByID(id int) (Post, error)
	}

	Posts struct {
		*sqlx.DB
	}

	Post struct {
		ID        int       `sq:"id" json:"id"`
		OwnerID   int       `sq:"owner_id" json:"owner_id"`
		Title     string    `sq:"title" json:"title"`
		Body      string    `sq:"body" json:"body"`
		CreatedAt time.Time `sq:"created_at" json:"created_at"`
		UpdatedAt time.Time `sq:"updated_at" json:"updated_at"`
		Updated   bool      `sq:"updated" json:"updated"`
	}
)

func (db *Posts) ByID(id int) (Post, error) {
	const q = "SELECT * FROM posts WHERE id = $1 LIMIT 1"

	var post Post
	err := db.Get(&post, q, id)

	if err != nil {
		return Post{}, errors.New("post not found")
	}

	return post, nil

}

func (db *Posts) Create(p Post) (post Post, err error) {
	const q = "INSERT INTO posts (owner_id, title, body) VALUES ($1, $2, $3) RETURNING *"

	err = db.Get(&post, q, p.OwnerID, p.Title, p.Body)

	if err != nil {
		return Post{}, err
	}

	return post, nil

}

func (db *Posts) Delete(p Post) error {
	const q = "DELETE FROM posts WHERE id = $1 AND owner_id = $2"

	_, err := db.Exec(q, p.ID, p.OwnerID)

	if err != nil {
		return err
	}

	return nil

}

func (db *Posts) Update(p Post) (post Post, err error) {
	const q = "UPDATE posts SET title = $1, body = $2, updated = true, updated_at = now() WHERE id = $3 AND owner_id = $4 RETURNING *"
	return post, db.Get(&post, q, p.Title, p.Body, p.ID, p.OwnerID)
}

func (db *Posts) UserPosts(id, limit, offset int) ([]Post, error) {
	const q = "SELECT * FROM posts WHERE owner_id = $1 LIMIT $2 OFFSET $3"

	var posts []Post

	rows, err := db.Query(q, id, limit, offset)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.OwnerID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.Updated)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}
