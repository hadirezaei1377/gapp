package presence

import (
	"context"
	"gapp/contract/golang/presence"
	"gapp/param"
	"gapp/pkg/protobufmapper"
	"gapp/pkg/slice"

	"google.golang.org/grpc"
)

type Client struct {
	client presence.PresenceServiceClient
}

func New(conn *grpc.ClientConn) Client {

	return Client{
		client: presence.NewPresenceServiceClient(conn),
	}
}

func (c Client) GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error) {
	resp, err := c.client.GetPresence(ctx,
		&presence.GetPresenceRequest{
			UserIds: slice.MapFromUintToUint64(request.UserIDs),
		})
	if err != nil {

		return param.GetPresenceResponse{}, err
	}

	return protobufmapper.MapGetPresenceResponseFromProtobuf(resp), nil
}
