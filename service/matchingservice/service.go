package matchingservice

import (
	"context"
	"fmt"
	"gapp/entity"
	"gapp/param"
	"gapp/pkg/protobufencoder"
	"gapp/pkg/richerror"
	"gapp/pkg/timestamp"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

type Publisher interface {
	Publish(event entity.Event, payload string)
}

// TODO - add ctx to ~all repo and use-case methods if needed
type Repo interface {
	AddToWaitingList(userID uint, category entity.Category) error
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	RemoveUsersFromWaitingList(category entity.Category, userIDs []uint)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, request param.GetPresenceRequest) (param.GetPresenceResponse, error)
}

type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config         Config
	repo           Repo
	presenceClient PresenceClient
	pub            Publisher
}

func New(config Config, repo Repo, presenceClient PresenceClient, pub Publisher) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient, pub: pub}
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

func (s Service) MatchWaitedUsers(ctx context.Context, _ param.MatchWaitedUsersRequest) (param.MatchWaitedUsersResponse, error) {
	const op = richerror.Op("matchingservice.MatchWaitedUsers")

	var wg sync.WaitGroup
	for _, category := range entity.CategoryList() {
		wg.Add(1)
		go s.match(ctx, category, &wg)
	}

	wg.Wait()

	return param.MatchWaitedUsersResponse{}, nil
}

func (s Service) match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	const op = richerror.Op("matchingservice.match")

	defer wg.Done()

	list, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		log.Errorf("GetWaitingListByCategory, err: %v\n", err)
		return
	}

	userIDs := make([]uint, 0)
	for _, l := range list {
		userIDs = append(userIDs, l.UserID)
	}

	if len(userIDs) < 2 {
		return
	}

	presenceList, err := s.presenceClient.GetPresence(ctx, param.GetPresenceRequest{UserIDs: userIDs})
	if err != nil {
		// TODO - log error
		// TODO - update metrics
		log.Errorf("GetWaitingListByCategory presenceClient.GetPresence, err: %v\n", err)
		return
	}

	presenceUserIDs := make([]uint, len(list))
	for _, l := range presenceList.Items {
		presenceUserIDs = append(presenceUserIDs, l.UserID)
	}

	var toBeRemovedUsers = make([]uint, 0)

	var finalList = make([]entity.WaitingMember, 0)
	for _, l := range list {
		lastOnlineTimestamp, ok := getPresenceItem(presenceList, l.UserID)
		// TODO - add 20 & 300 to config
		if ok && lastOnlineTimestamp > timestamp.Add(-20*time.Second) &&
			l.Timestamp > timestamp.Add(-300*time.Second) {
			finalList = append(finalList, l)
		} else {
			// remove from waiting list
			toBeRemovedUsers = append(toBeRemovedUsers, l.UserID)
		}
	}

	go s.repo.RemoveUsersFromWaitingList(category, toBeRemovedUsers)

	matchedUsersToBeRemoved := make([]uint, 0)
	for i := 0; i < len(finalList)-1; i = i + 2 {
		mu := entity.MatchedUsers{
			Category: category,
			UserIDs:  []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		fmt.Println("mu", mu)

		go s.pub.Publish(entity.MatchingUsersMatchedEvent,
			protobufencoder.EncodeEvent(entity.MatchingUsersMatchedEvent, mu))
		// publish a new event for mu

		// remove mu users from waiting list
		matchedUsersToBeRemoved = append(matchedUsersToBeRemoved, mu.UserIDs...)
	}

	go s.repo.RemoveUsersFromWaitingList(category, matchedUsersToBeRemoved)
}

func getPresenceItem(presenceList param.GetPresenceResponse, userID uint) (int64, bool) {
	for _, item := range presenceList.Items {
		if item.UserID == userID {
			return item.Timestamp, true
		}
	}

	return 0, false
}
