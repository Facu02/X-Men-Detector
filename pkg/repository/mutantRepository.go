package repository

import (
	"context"
	"x-menDetector/internal/models"

	"github.com/go-redis/redis/v8"
)

type MutanRepositoryImp struct {
	redisDb *redis.Client
}

func NewMutanRepository(redisDb *redis.Client) models.MutanRepository {
	return MutanRepositoryImp{redisDb: redisDb}
}

func (repo MutanRepositoryImp) IncrementCounter(counterName string) error {
	ctx := context.Background()

	return repo.redisDb.Incr(ctx, counterName).Err()
}

func (repo MutanRepositoryImp) GetCounter(counterName string) (string, error) {
	ctx := context.Background()

	value, err := repo.redisDb.Get(ctx, counterName).Result()
	if err == redis.Nil {
		return "0", nil
	} else if err != nil {
		return "", err
	}

	return value, nil
}
