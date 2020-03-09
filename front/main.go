package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
)

func testGet(c echo.Context) error {
	return c.Render(http.StatusOK, "signin", "world")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/", testGet)

	e.Logger.Fatal(e.Start(":7070"))
}
