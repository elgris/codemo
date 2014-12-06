package ideat

import (
	"codemo/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ngerakines/ginpongo2"
)

type App struct {
	*gin.Engine
	workdir string
}

func NewApp() (*App, error) {
	a := &App{
		workdir: "../",
	}

	return a, a.Init()
}

func (a *App) Init() error {
	a.Engine = gin.Default()

	a.Use(ginpongo2.Pongo2())
	a.Static("/assets", "assets")

	a.GET("/", func(c *gin.Context) {
		c.Set("template", "../views/index.tpl.html")
		c.Set("data", map[string]interface{}{"message": "Hello World!"})
	})

	a.POST("src", func(c *gin.Context) {
		form := models.SourceForm{}

		c.BindWith(&form, binding.Form)

		f, err := models.ParseSrc(form.Source)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, f)
		// TODO: return JSON
	})

	return nil
}
