package scheduler

import (
	"fmt"
	"sync"
	"time"

	"gapp/param"
	"gapp/service/matchingservice"
	"time"

	"github.com/go-co-op/gocron"
)

type Config struct {
	MatchWaitedUsersIntervalInSeconds int `koanf:"match_waited_users_interval_in_seconds"`
}
type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
	config   Config
}

func New(config Config, matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		config:   config,
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

// long-running process
func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	s.sch.Every(s.config.MatchWaitedUsersIntervalInSeconds).Second().Do(s.MatchWaitedUsers)
	s.sch.StartAsync()
	<-done
	// wait to finish job
	fmt.Println("stop scheduler..")
	s.sch.Stop()
}

func (s Scheduler) MatchWaitedUsers() {
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
