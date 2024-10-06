package matchingservice

import (
	"fmt"
	"gapp/entity"
	"gapp/param"
	"gapp/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWaitingList(userID uint, category entity.Category) error
}
type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}
type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}
func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (
	param.AddToWaitingListResponse, error) {
	const op = richerror.Op("matchingservice.AddToWaitingList")
	// add user to the waiting list for the given category if not exist
	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	if err != nil {
		return param.AddToWaitingListResponse{},
			richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}
	return param.AddToWaitingListResponse{Timeout: s.config.WaitingTimeout}, nil
}

func (s Service) MatchWaitedUsers(req param.MatchWaitedUsersRequest) (param.MatchWaitedUsersResponse, error) {
	fmt.Println("MatchWaitedUsers", time.Now())
	return param.MatchWaitedUsersResponse{}, nil
}
