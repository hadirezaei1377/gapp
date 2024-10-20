package broker

import "gapp/entity"

type Publisher interface {
	Publish(event entity.Event, payload string)
}
