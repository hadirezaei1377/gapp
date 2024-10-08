package redismatching

import (
	"context"
	"fmt"
	"gapp/entity"
	"gapp/pkg/richerror"

	"gapp/pkg/timestamp"

	"github.com/redis/go-redis/v9"
)

// TODO - add to config in usecase layer...
const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const op = richerror.Op("redismatching.AddToWaitingList")

	_, err := d.adapter.Client().
		ZAdd(context.Background(),
			fmt.Sprintf("%s:%s", WaitingListPrefix, category),
			redis.Z{Score: float64(timestamp.Now()),
				Member: fmt.Sprintf("%d", userID),
			}).Result()
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
