package protobufencoder

import (
	"encoding/base64"
	"gapp/contract/golang/matching"
	"gapp/entity"
	"gapp/pkg/slice"

	"google.golang.org/protobuf/proto"
)

func EncodeEvent(event entity.Event, data any) string {
	var payload []byte

	switch event {
	case entity.MatchingUsersMatchedEvent:
		mu, ok := data.(entity.MatchedUsers)
		if !ok {
			// TODO - log
			// TODO - update metrics
			return ""
		}

		pbMu := matching.MatchedUsers{
			Category: string(mu.Category),
			UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
		}

		var err error
		payload, err = proto.Marshal(&pbMu)
		if err != nil {
			// TODO - log error
			// TODO - update metrics
			return ""
		}
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func DecodeEvent(event entity.Event, data string) any {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// TODO - log
		// TODO - update metrics
		return nil
	}

	switch event {
	case entity.MatchingUsersMatchedEvent:
		pbMu := matching.MatchedUsers{}
		if err := proto.Unmarshal(payload, &pbMu); err != nil {
			// TODO - log
			// TODO - update metrics
			return nil
		}

		return entity.MatchedUsers{
			Category: entity.Category(pbMu.Category),
			UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
		}
	}

	return nil
}
