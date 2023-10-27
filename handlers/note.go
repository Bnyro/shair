package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shair/db"
	"github.com/shair/entities"
)

func ListNotes(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth")
	}
	exampleNote := entities.Note{
		UserID: user.ID,
	}

	var notes []entities.Note
	db.Database.Where(&exampleNote).Find(&notes)

	return c.Render(http.StatusOK, "notes.html", notes)
}

func NewNote(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/auth")
	}

	note := entities.Note{
		Title:  c.FormValue("title"),
		Body:   c.FormValue("body"),
		UserID: user.ID,
	}
	db.Database.Create(&note)

	return c.Redirect(http.StatusTemporaryRedirect, "/notes/")
}

func DeleteNote(c echo.Context) error {
	user, err := getUserByCookie(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	note := entities.Note{
		UserID: user.ID,
	}
	note.ID = uint(id)

	if db.Database.Where(&note).Delete(&note).RowsAffected == 0 {
		return errors.New("Note not found or owned by this user!")
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/notes/")
}
