package entity

import "time"

type Game struct {
	ID          uint
	Category    Category
	QuestionIDs []uint
	PlayerIDs   []uint
	StartTime   time.Time
}

// player is user who is playing in specified game (determine by ID)
type Player struct {
	ID      uint
	UserID  uint
	GameID  uint
	Score   uint
	Answers []PlayerAnswer
}

type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Choice     PossibleAnswerChoice
}

func data() {
}
