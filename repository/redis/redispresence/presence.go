package redispresence

import (
	"context"
	"fmt"
	"gapp/pkg/richerror"
	"time"
)

func (d DB) Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error {
	const op = richerror.Op("redispresence.Upsert")
	_, err := d.adapter.Client().Set(ctx, key, timestamp, expTime).Result()
	if err != nil {
		fmt.Println("UpsertPresence3 err", err.Error())

		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
