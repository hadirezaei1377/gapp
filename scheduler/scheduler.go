package scheduler

import (
	"fmt"
	"sync"
	"time"

	"gapp/service/matchingservice"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
}

func New(matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

// long-running process
func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	s.sch.Every(3).Second().Do(s.MatchWaitedUsers)
	s.sch.StartAsync()
	<-done
	// wait to finish job
	fmt.Println("stop scheduler..")
	s.sch.Stop()
}

func (s Scheduler) MatchWaitedUsers() {
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
