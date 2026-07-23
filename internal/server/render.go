package server

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

var templates *template.Template

func init() {
	templates = template.Must(templates.ParseGlob("./webpages/*.html"))
}

func servePage(c *gin.Context, page string, data any) {
	if err := templates.ExecuteTemplate(c.Writer, page+".html", data); err != nil {
		fmt.Println(err)
	}
}
