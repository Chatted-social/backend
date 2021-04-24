package storage

import (
	"context"
	"time"
	"unicode"

	"github.com/fatih/structs"
	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func init() {
	// structs is used with squirrel (sq)
	structs.DefaultTagName = "sq"
}

type DB struct {
	*sqlx.DB
	Users UsersStorage
	Posts PostStorage
}

type RedisCache struct {
	*redis.Client
}

func (r RedisCache) Get(key string) ([]byte, error) {
	return r.Client.Get(context.Background(), key).Bytes()
}

func (r RedisCache) Set(key string, val []byte, ttl time.Duration) error {
	return r.Client.Set(context.Background(), key, val, ttl).Err()
}

func (r RedisCache) Delete(key string) error {
	return r.Delete(key)
}

func (r RedisCache) Reset() error {
	return r.Reset()
}

func Open(url string) (*DB, error) {
	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("db", toSnakeCase)

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return &DB{
		DB:    db,
		Users: &Users{DB: db},
		Posts: &Posts{DB: db},
	}, nil
}

func NewRedisCache(opt *redis.Options) (*RedisCache, error) {
	client := redis.NewClient(opt)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &RedisCache{
		Client: client,
	}, nil
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) {
			prev := unicode.IsLower(runes[i-1])
			next := i+1 < length && unicode.IsLower(runes[i+1])

			if prev || next {
				out = append(out, '_')
			}
		}

		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
