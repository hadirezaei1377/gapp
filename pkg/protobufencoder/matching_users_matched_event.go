package protobufencoder

import (
	"encoding/base64"
	"gapp/contract/golang/matching"
	"gapp/entity"
	"gapp/pkg/slice"

	"google.golang.org/protobuf/proto"
)

func EncodeMatchingUsersMatchedEvent(mu entity.MatchedUsers) string {
	pbMu := matching.MatchedUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
	}

	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func DecodeMatchingUsersMatchedEvent(data string) entity.MatchedUsers {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// TODO - log
		// TODO - update metrics
		return entity.MatchedUsers{}
	}

	pbMu := matching.MatchedUsers{}
	if err := proto.Unmarshal(payload, &pbMu); err != nil {
		// TODO - log
		// TODO - update metrics
		return entity.MatchedUsers{}
	}

	return entity.MatchedUsers{
		Category: entity.Category(pbMu.Category),
		UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
	}

}
