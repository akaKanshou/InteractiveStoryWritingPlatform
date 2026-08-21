package server

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

var templates, templates2 *template.Template

func init() {
	templates = template.Must(templates.ParseGlob("./webpages2/*.html"))
	templates2 = template.Must(templates.ParseGlob("./webpages2/*.html"))
}

func servePage(c *gin.Context, page string, data any) {
	if err := templates.ExecuteTemplate(c.Writer, page+".html", data); err != nil {
		fmt.Println(err)
	}
}

func servePage2(c *gin.Context, page string, data any) {
	if err := templates.ExecuteTemplate(c.Writer, page+".html", data); err != nil {
		fmt.Println(err)
	}
}
