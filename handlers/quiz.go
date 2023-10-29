package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/shair/db"
	"github.com/shair/entities"
	"github.com/shair/util"
)

const amountOfChoices = 4

func NewQuizOptions(c echo.Context) error {
	return c.Render(http.StatusOK, "newquiz.html", nil)
}

func CreateNewQuiz(c echo.Context) error {
	amountOfQuestions, err := strconv.Atoi(c.FormValue("questioncount"))

	if err != nil {
		return err
	}

	quiz := entities.Quiz{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Token:       util.GenerateSecureToken(20),
	}
	db.Database.Create(&quiz)

	questionsIter := util.MakeIntArray(amountOfQuestions)
	choicesIter := util.MakeIntArray(amountOfChoices)

	return c.Render(http.StatusOK, "newquizquestions.html", echo.Map{
		"Token":         quiz.Token,
		"QuestionsIter": questionsIter,
		"ChoicesIter":   choicesIter,
		"QuestionCount": amountOfQuestions,
	})
}

func CreateNewQuizQuestions(c echo.Context) error {
	token := c.FormValue("token")
	questionCount, err := strconv.Atoi(c.FormValue("questioncount"))
	if err != nil {
		return err
	}

	for questionIndex := 0; questionIndex < questionCount; questionIndex++ {
		var choices pq.StringArray
		for choiceIndex := 0; choiceIndex < amountOfChoices; choiceIndex++ {
			choice := c.FormValue(fmt.Sprintf("choice_%d_%d", questionIndex, choiceIndex))
			if !util.IsBlank(choice) {
				choices = append(choices, choice)
			}
		}

		question := entities.Question{
			QuizToken:       token,
			Question:        c.FormValue(fmt.Sprintf("question_%d", questionIndex)),
			CorrectResponse: c.FormValue(fmt.Sprintf("correct_%d", questionIndex)),
			PossibleChoices: choices,
		}
		db.Database.Create(&question)
	}

	return util.CreateSuccessResult(c, "quiz", token)
}

func GetQuiz(c echo.Context) error {
	_ = c.Param("token")

	return nil
}

func SubmitQuizResponse(c echo.Context) error {
	return nil
}
