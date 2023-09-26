package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shair/db"
	"github.com/shair/entities"
	"github.com/shair/util"
)

func Blog(c echo.Context) error {
	var posts []entities.BlogPost
	db.Database.Find(&posts)

	var renderedPosts []echo.Map
	for _, post := range posts {
		body := util.MdToHTML(util.FirstN(post.Body, 150))
		renderedPosts = append(renderedPosts, echo.Map{
			"ID":        post.ID,
			"Reference": post.Reference,
			"Title":     post.Title,
			"Body":      template.HTML(body),
		})
	}

	return c.Render(http.StatusOK, "blog.html", echo.Map{
		"Admin": isAdmin(c),
		"Posts": renderedPosts,
	})
}

func NewBlogPost(c echo.Context) error {
	if !isAdmin(c) {
		return errors.New("Not authorized!")
	}
	reference := strings.ToLower(strings.ReplaceAll(c.FormValue("title"), " ", "-"))

	post := entities.BlogPost{
		Title:     c.FormValue("title"),
		Reference: reference,
		Body:      c.FormValue("body"),
	}
	db.Database.Create(&post)

	return c.Redirect(http.StatusTemporaryRedirect, "/blog/")
}

func BlogPost(c echo.Context) error {
	var post entities.BlogPost

	if db.Database.Where("reference = ?", c.Param("reference")).Find(&post).RowsAffected == 0 {
		return errors.New("Post not found!")
	}

	return c.Render(http.StatusOK, "post.html", echo.Map{
		"Title":   post.Title,
		"Body":    template.HTML(util.MdToHTML(post.Body)),
		"Created": post.CreatedAt.Format(time.RFC822),
		"Words":   len(strings.Split(post.Body, " ")),
	})
}

func DeleteBlogPost(c echo.Context) error {
	if !isAdmin(c) {
		return errors.New("Not authorized!")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	post := entities.BlogPost{}
	post.ID = uint(id)

	if db.Database.Where(&post).Delete(&post).RowsAffected == 0 {
		return errors.New("Post not found!")
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/blog/")
}
