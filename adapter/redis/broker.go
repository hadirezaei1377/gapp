package redis

import (
	"context"
	"gapp/entity"
	"time"

	"github.com/labstack/gommon/log"
)

func (a Adapter) Publish(event entity.Event, payload string) {
	// TODO - add 1 to config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	if err := a.client.Publish(ctx, string(event), payload).Err(); err != nil {
		log.Errorf("publish err: %v\n", err)
		// TODO - log
		// TODO - update metrics
	}

	// TODO - update metrics
}
