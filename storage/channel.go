package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type (
	ChannelStorage interface {
		Create(Channel) (int, error)
		Update(Channel) error
		Delete(Channel) error

		ByUsername(username string) (Channel, error)

		ExistsByUsername(string) (bool, error)
		ExistsByID(int) (bool, error)

		UserPromote(int, int) error
		UserDemote(int, int) error
		UserBlock(int, int) error

		UserSubscribe(int, int) error
		UserUnsubscribe(int, int) error

		UserIsAdmin(int, int) (bool, error)
		UserIsOwner(int, int) (bool, error)
		UserIsSubscriber(int, int) (bool, error)
		UserIsBlocked(int, int) bool

		LatestPostID(int) int
	}

	Channels struct {
		*sqlx.DB
	}

	Channel struct {
		Title          string        `json:"title" sq:"title" db:"title"`
		Username       string        `json:"username" sq:"username" db:"username"`
		ID             int           `json:"id" sq:"id" db:"id"`
		OwnerID        int           `json:"owner_id" sq:"owner_id" db:"owner_id"`
		BlockedIDs     pq.Int64Array `json:"blocked_ids" sq:"blocked_ids" db:"blocked_ids"`
		SubscribersIDs pq.Int64Array `json:"subscribers_ids" sq:"subscribers_ids" db:"subscribers_ids"`
		AdminsIDs      pq.Int64Array `json:"admins_ids" sq:"admins_ids" db:"admins_ids"`
		CreatedAt      time.Time     `json:"created_at" sq:"created_at" db:"created_at"`
	}
)

func (db Channels) Create(c Channel) (id int, err error) {
	const q = "INSERT INTO channels (owner_id, username, title) VALUES ($1, $2, $3) RETURNING id"

	err = db.Get(&id, q, c.OwnerID, c.Username, c.Title)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (db Channels) Update(c Channel) error {
	const q = "UPDATE channels SET title = $1, username = $2 WHERE owner_id = $3 AND id = $4"

	_, err := db.Exec(q, c.Title, c.Username, c.OwnerID, c.ID)

	if err != nil {
		return err
	}

	return nil

}

func (db Channels) Delete(c Channel) error {
	const q = "DELETE from channels WHERE id = $1 AND owner_id = $2"

	_, err := db.Exec(q, c.ID, c.OwnerID)

	if err != nil {
		return err
	}

	return nil

}

func (db Channels) ByUsername(username string) (Channel, error) {
	const q = "SELECT * FROM channels WHERE username = $1"

	var channel Channel
	rows, err := db.Queryx(q, username)

	if err != nil {
		return Channel{}, err
	}

	for rows.Next() {
		err := rows.StructScan(&channel)

		if err != nil {
			return Channel{}, err
		}

	}

	return channel, nil

}

func (db Channels) ExistsByUsername(username string) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM channels WHERE username = $1)"
	row := db.QueryRow(q, username)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

func (db Channels) ExistsByID(id int) (bool, error) {
	const q = "SELECT EXISTS(SELECT * FROM channels WHERE id = $1)"
	row := db.QueryRow(q, id)

	var exists bool
	err := row.Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}

func (db Channels) UserIsAdmin(ChannelID, UserID int) (bool, error) {
	const q = "SELECT COUNT(*) FROM channels WHERE $1 = ANY (admins_ids) AND id = $2"

	row := db.QueryRow(q, UserID, ChannelID)

	var count int
	row.Scan(&count)

	return count == 1, nil

}

func (db Channels) UserIsOwner(CID, UID int) (bool, error) {
	const q = "SELECT COUNT(*) FROM channels WHERE owner_id = $1 AND id = $2"

	row := db.QueryRow(q, UID, CID)

	var count int
	row.Scan(&count)

	return count == 1, nil

}

func (db Channels) UserIsSubscriber(ChannelID, UserID int) (bool, error) {
	const q = "SELECT COUNT(*) FROM channels WHERE $1 = ANY (subscribers_ids) AND id = $2"

	row := db.QueryRow(q, UserID, ChannelID)

	var count int
	row.Scan(&count)

	return count == 1, nil

}

func (db Channels) UserIsBlocked(ChannelID, UserID int) bool {
	const q = "SELECT COUNT(*) FROM channels WHERE $1 = ANY (blocked_ids) AND id = $2"

	row := db.QueryRow(q, UserID, ChannelID)

	var count int
	row.Scan(&count)

	return count == 1

}

func (db Channels) UserPromote(ChannelID int, UserID int) error {
	const q = "UPDATE channels SET admins_ids = array_append(admins_ids, $1) WHERE id = $2 RETURNING id"
	return db.QueryRow(q, UserID, ChannelID).Err()
}

func (db Channels) UserDemote(ChannelID int, UserID int) error {
	const q = "UPDATE channels SET admins_ids = array_remove(admins_ids, $1) WHERE id = $2 RETURNING id"
	return db.QueryRow(q, UserID, ChannelID).Err()
}

func (db Channels) UserBlock(ChannelID, UserID int) error {
	const q = "UPDATE channels SET blocked_ids = array_append(blocked_ids, $1) WHERE id = $2"

	_, err := db.Exec(q, UserID, ChannelID)

	if err != nil {
		return err
	}

	return nil

}

func (db Channels) UserSubscribe(ChannelID, UserID int) error {
	const q = "UPDATE channels SET subscribers_ids = array_append(subscribers_ids, $1) WHERE id = $2"

	return db.QueryRow(q, UserID, ChannelID).Err()

}

func (db Channels) UserUnsubscribe(ChannelID, UserID int) error {
	const q = "UPDATE channels SET subscribers_ids = array_remove(subscribers_ids, $1) WHERE id = $2"

	return db.QueryRow(q, UserID, ChannelID).Err()
}

func (db Channels) LatestPostID(ChannelID int) int {
	const q = "SELECT post_id FROM posts WHERE owner_id = $1 ORDER BY post_id DESC LIMIT 1"

	row := db.QueryRow(q, ChannelID)

	var latest int
	row.Scan(&latest)

	return latest

}
