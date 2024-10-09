package presenceservice

import (
	"context"
	"fmt"
	"gapp/param"
	"gapp/pkg/richerror"
	"time"
)

type Config struct {
	ExpirationTime time.Duration `koanf:"expiration_time"`
	Prefix         string        `koanf:"prefix"`
}

type Repo interface {
	Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error
}

type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}

func (s Service) Upsert(ctx context.Context, req param.UpsertPresenceRequest) (param.UpsertPresenceResponse, error) {
	const op = richerror.Op("presenceservice.Upsert")
	err := s.repo.Upsert(ctx,
		fmt.Sprintf("%s:%d", s.config.Prefix, req.UserID),
		req.Timestamp, s.config.ExpirationTime)
	if err != nil {
		fmt.Println("UpsertPresence2 err", err.Error())

		return param.UpsertPresenceResponse{}, richerror.New(op).WithErr(err)
	}

	return param.UpsertPresenceResponse{}, nil
}

func (s Service) GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error) {
	fmt.Println("req", request)
	// TODO - implement me
	return param.GetPresenceResponse{Items: []param.GetPresenceItem{
		{UserID: 1, Timestamp: 12452151},
		{UserID: 2, Timestamp: 124534551},
	}}, nil
}
