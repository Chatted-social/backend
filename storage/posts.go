package storage

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type (
	PostStorage interface {
		Create(p Post) (post Post, err error)
		Update(p Post) (post Post, err error)
		Delete(p Post) error

		PostsIn(ids []int) (posts []Post, err error)
		ByIDs(OwnerID, PostID int) (Post, error)
		UserPosts(id, limit, offset int) (posts []Post, _ error)

		Like(UserID int, p Post) error
		Unlike(UserID int, p Post) error
		Liked(UserID int, p Post) (bool, error)
	}

	Posts struct {
		*sqlx.DB
	}

	Post struct {
		ID        int           `sq:"id" json:"-" db:"id"`
		PostID    int           `sq:"post_id" json:"post_id" db:"post_id"`
		OwnerID   int           `sq:"owner_id" json:"owner_id" db:"owner_id"`
		FromID    int           `sq:"from_id" json:"from_id" db:"from_id"`
		Body      string        `sq:"body" json:"body" db:"body"`
		CreatedAt time.Time     `sq:"created_at" json:"created_at" db:"created_at"`
		UpdatedAt time.Time     `sq:"updated_at" json:"updated_at" db:"updated_at"`
		LikesIDs  pq.Int64Array `sq:"likes_ids" json:"likes_ids" db:"likes_ids"`
		Updated   bool          `sq:"updated" json:"updated" db:"updated"`
	}
)

func (db *Posts) ByIDs(OwnerID, PostID int) (Post, error) {
	const q = "SELECT * FROM posts WHERE owner_id = $1 AND post_id = $2"

	var post Post
	err := db.Get(&post, q, OwnerID, PostID)

	if err != nil {
		return Post{}, errors.New("post not found")
	}

	return post, nil

}

func (db *Posts) Create(p Post) (post Post, err error) {
	const q = "INSERT INTO posts (post_id, owner_id, from_id, body) VALUES ($1, $2, $3, $4) RETURNING *"

	err = db.Get(&post, q, p.PostID, p.OwnerID, p.FromID, p.Body)

	if err != nil {
		return Post{}, err
	}

	return post, nil

}

func (db *Posts) Delete(p Post) error {
	const q = "DELETE FROM posts WHERE post_id = $1 AND owner_id = $2 AND from_id = $3"

	_, err := db.Exec(q, p.PostID, p.OwnerID, p.FromID)

	if err != nil {
		return err
	}

	return nil

}

func (db *Posts) Update(p Post) (post Post, err error) {
	const q = "UPDATE posts SET body = $1, updated = true WHERE post_id = $2 AND owner_id = $3 RETURNING *"
	return post, db.Get(&post, q, p.Body, p.PostID, p.OwnerID)
}

func (db *Posts) PostsIn(ids []int) (posts []Post, err error) {
	const q = `
		SELECT
		*
		FROM posts
		WHERE post_id IN (?)`

	query, args, err := sqlx.In(q, ids)
	if err != nil {
		return nil, err
	}

	return posts, db.Select(&posts, db.Rebind(query), args...)
}

func (db *Posts) UserPosts(id, limit, offset int) (posts []Post, _ error) {
	const q = "SELECT * FROM posts WHERE owner_id = $1 LIMIT $2 OFFSET $3"
	return posts, db.Select(&posts, q, id, limit, offset)
}

func (db *Posts) Like(UserID int, p Post) error {
	const q = "UPDATE posts SET likes_ids = array_append(likes_ids, $1) WHERE owner_id = $2 AND post_id = $3"

	return db.QueryRow(q, UserID, p.OwnerID, p.PostID).Err()

}

func (db *Posts) Unlike(UserID int, p Post) error {
	const q = "UPDATE posts SET likes_ids = array_remove(likes_ids, $1) WHERE owner_id = $2 AND post_id = $3"

	return db.QueryRow(q, UserID, p.OwnerID, p.PostID).Err()

}

func (db *Posts) Liked(UserID int, p Post) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM posts WHERE owner_id = $1 AND post_id = $2 AND likes_ids = ANY($3))"

	var exists bool
	return exists, db.QueryRow(q, p.OwnerID, p.PostID, UserID).Scan(&exists)

}
