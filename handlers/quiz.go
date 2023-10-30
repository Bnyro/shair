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
			CorrectAnswer:   c.FormValue(fmt.Sprintf("correct_%d", questionIndex)),
			PossibleAnswers: choices,
		}
		db.Database.Create(&question)
	}

	return util.CreateSuccessResult(c, "quiz", token)
}

func getQuizByToken(token string) (entities.Quiz, []entities.Question) {
	var quiz entities.Quiz
	var questions []entities.Question
	db.Database.Where("token = ?", token).First(&quiz)
	db.Database.Where("quiz_token = ?", token).Find(&questions)
	return quiz, questions
}

func GetQuiz(c echo.Context) error {
	quiz, questions := getQuizByToken(c.Param("token"))

	return c.Render(http.StatusOK, "solvequiz.html", echo.Map{
		"Quiz":      quiz,
		"Questions": questions,
	})
}

func SubmitQuizResponse(c echo.Context) error {
	quiz, questions := getQuizByToken(c.Param("token"))

	var solutions []echo.Map
	correctCount := 0
	totalCount := len(questions)

	for questionIndex, question := range questions {
		providedAnswer := c.FormValue(fmt.Sprintf("response_%d", questionIndex))

		isCorrect := providedAnswer == question.CorrectAnswer
		if isCorrect {
			correctCount++
		}

		solutions = append(solutions, echo.Map{
			"Question":  question,
			"IsCorrect": isCorrect,
			"Provided":  providedAnswer,
		})
	}

	return c.Render(http.StatusOK, "quizsolution.html", echo.Map{
		"Quiz":         quiz,
		"Solutions":    solutions,
		"CorrectCount": correctCount,
		"TotalCount":   totalCount,
	})
}
