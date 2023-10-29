package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	QuizToken       string         `json:"token"`
	Question        string         `json:"question"`
	CorrectResponse string         `json:"correct"`
	PossibleChoices pq.StringArray `json:"choices" gorm:"type:text[]"`
}
