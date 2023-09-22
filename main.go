package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shair/config"
	"github.com/shair/db"
	"github.com/shair/handlers"
)

type Template struct{}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	theme := "catppuccin"
	themeCookie, err := c.Cookie("Theme")
	if err == nil {
		theme = themeCookie.Value
	}
	themeTmpl := "templates/" + theme + ".html"

	tmpl, err := template.ParseFiles("templates/base.html", themeTmpl, "templates/"+name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return tmpl.Execute(w, data)
}

func main() {
	db.Init()

	config.Init()

	router := echo.New()
	router.Use(middleware.CORS())
	router.Renderer = &Template{}

	router.Static("/static", "static")
	router.Static("/files", "files")

	router.GET("/", handlers.HomeHandler)
	router.POST("/", handlers.HomeHandler)

	router.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "ok",
		})
	})

	user := router.Group("/user")
	user.POST("/register", handlers.RegisterUser)
	user.POST("/login", handlers.LoginUser)
	user.POST("/logout", handlers.LogoutUser)
	user.POST("/delete", handlers.DeleteUser)

	paste := router.Group("/paste")
	paste.GET("/", handlers.NewPaste)
	paste.POST("/", handlers.NewPaste)
	paste.GET("/:id", handlers.GetPaste)

	upload := router.Group("/upload")
	upload.GET("/", handlers.NewUpload)
	upload.POST("/", handlers.NewUpload)
	upload.GET("/:id", handlers.GetUpload)

	notes := router.Group("/notes")
	notes.GET("/", handlers.ListNotes)
	notes.POST("/", handlers.ListNotes)
	notes.POST("/new", handlers.NewNote)
	notes.POST("/delete/:id", handlers.DeleteNote)

	router.Logger.Fatal(router.Start(":3000"))
}
