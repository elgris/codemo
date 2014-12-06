package ideat

import (
	"codemo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ngerakines/ginpongo2"
	"go/ast"
	"go/parser"
	"go/token"
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

		f, err := parseSrc(form.Source)
		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(200, f)
		// TODO: return JSON
	})

	return nil
}

func parseSrc(src string) (f *ast.File, err error) {
	fset := token.NewFileSet()

	fullSrc := fmt.Sprintf("package codemotst\n func main() {\n %s \n }", src)
	fmt.Println(fullSrc)
	return parser.ParseFile(fset, "", fullSrc, 0)
}

func convertAst(f *ast.File)
