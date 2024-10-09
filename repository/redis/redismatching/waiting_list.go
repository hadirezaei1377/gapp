package redismatching

import (
	"context"
	"fmt"
	"gapp/entity"
	"gapp/pkg/richerror"
	"strconv"
	"time"

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

func (d DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const op = richerror.Op("redismatching.GetWaitingListByCategory")
	// TODO - add to config
	min := fmt.Sprintf("%d", timestamp.Add(-200000*time.Hour))
	max := strconv.Itoa(int(timestamp.Now()))
	list, err := d.adapter.Client().ZRangeByScoreWithScores(ctx, getCategoryKey(category), &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  0,
	}).Result()

	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	var result = make([]entity.WaitingMember, 0)
	for _, l := range list {
		userID, _ := strconv.Atoi(l.Member.(string))
		result = append(result, entity.WaitingMember{
			UserID:    uint(userID),
			Timestamp: int64(l.Score),
			Category:  category,
		})
	}
	return result, nil
}

func getCategoryKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", WaitingListPrefix, category)
}
