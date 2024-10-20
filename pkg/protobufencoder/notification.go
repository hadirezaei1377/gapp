package protobufencoder

import (
	"encoding/base64"
	"gapp/contract/goproto/notification"
	"gapp/entity"

	"google.golang.org/protobuf/proto"
)

func EncodeNotification(mu entity.Notification) string {
	pbMu := notification.Notification{
		Type:    mu.EventType,
		Payload: mu.Payload,
	}

	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func DecodeNotification(data string) entity.Notification {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// TODO - log
		// TODO - update metrics
		return entity.Notification{}
	}

	pbMu := notification.Notification{}
	if err := proto.Unmarshal(payload, &pbMu); err != nil {
		// TODO - log
		// TODO - update metrics
		return entity.Notification{}
	}

	return entity.Notification{
		EventType: pbMu.Type,
		Payload:   pbMu.Payload,
	}

}
