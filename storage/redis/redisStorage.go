package redisStorage

import (
	"dkubanyi/urlShortener/storage"
	"errors"
	"github.com/go-redis/redis"
)

func New(host, port string) (storage.Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "", // no password set
		DB:       0,  // use default db
	})

	_, err := client.Ping().Result()

	return &redisStorage{client}, err
}

type redisStorage struct{ db *redis.Client }

func (r *redisStorage) Save(url string) (string, error) {
	return "", errors.New("not yet implemented")
}

func (r *redisStorage) Load(code string) (string, error) {
	return r.db.Get("lnk:" + code).Result()
}

func (r *redisStorage) Close() error {
	return r.db.Close()
}
