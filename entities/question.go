package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	QuizToken       string         `json:"token"`
	Question        string         `json:"question"`
	CorrectAnswer   string         `json:"correct"`
	PossibleAnswers pq.StringArray `json:"choices" gorm:"type:text[]"`
}
