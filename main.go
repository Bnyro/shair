package main

import (
	"fmt"
	"html/template"
	"io"

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

	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	tmpl, _ = tmpl.Parse(fmt.Sprintf("{{ define \"theme\" }}class=\"%s\"{{ end }}", theme))

	return tmpl.Execute(w, data)
}

func main() {
	db.Init()

	config.Init()

	go runWorker()

	router := echo.New()
	router.Use(middleware.CORS())
	router.Renderer = &Template{}

	router.Static("/static", "static")

	router.GET("/", handlers.Home)
	router.POST("/", handlers.Home)

	router.GET("/status", handlers.Status)
	router.POST("/theme", handlers.SetTheme)

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
	router.Static("/files", config.UploadDir)

	notes := router.Group("/notes")
	notes.GET("/", handlers.ListNotes)
	notes.POST("/", handlers.ListNotes)
	notes.POST("/new", handlers.NewNote)
	notes.POST("/delete/:id", handlers.DeleteNote)

	downloads := router.Group("/downloads")
	downloads.GET("/", handlers.Files)
	downloads.POST("/", handlers.Files)
	downloads.POST("/delete/:filename", handlers.DeleteFile)
	router.Static("/dl", config.DownloadDir)

	gallery := router.Group("/gallery")
	gallery.GET("/", handlers.Gallery)
	gallery.POST("/", handlers.Gallery)
	gallery.GET("/dia/", handlers.GalleryDia)
	gallery.POST("/delete/:filename", handlers.DeleteImage)
	router.Static("/images", config.GalleryDir)

	blog := router.Group("/blog")
	blog.GET("/", handlers.Blog)
	blog.POST("/", handlers.Blog)
	blog.GET("/:reference/", handlers.BlogPost)
	blog.POST("/new", handlers.NewBlogPost)
	blog.POST("/delete/:id", handlers.DeleteBlogPost)

	router.Logger.Fatal(router.Start(":3000"))
}
