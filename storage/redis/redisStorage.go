package redisStorage

import (
	"dkubanyi/urlShortener/encoding"
	"dkubanyi/urlShortener/storage"
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
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
	existingLink, err := r.db.Get("srclnk:" + url).Result()

	if existingLink != "" {
		return existingLink, err
	}

	set := false
	encoded := ""

	for i := 1; i < 50; i++ {
		key := rand.Uint32()
		encoder := encoding.Base62Encoder{}
		encoded = encoder.Encode(int64(key))
		set, err = r.db.SetNX("lnk:"+encoded, url, 0).Result()

		if set {
			break
		}
	}

	if !set {
		err = errors.New("failed to create new link")
	}

	if err != nil {
		return "", err
	}

	err = r.db.Set("srclnk:"+url, encoded, 0).Err()

	return encoded, err
}

func (r *redisStorage) Load(code string) (string, error) {
	return r.db.Get("lnk:" + code).Result()
}

func (r *redisStorage) Close() error {
	return r.db.Close()
}
