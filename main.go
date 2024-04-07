package main

import (
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo"

	"example.com/habits"
)

type TestData struct {
	Name  string
	Phone string
	Email string
}
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func main() {
	e := echo.New()
	e.Static("/dist", "dist")
	e.Debug = true

	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e.Renderer = t

	habitList := make([]habits.Habit, 0)
	activity := habits.Activity{Name: "tester", Description: "asdf"}

	habit := habits.Habit{
		Activity: activity,
		Date:     time.Now(),
	}

	habitList = append(habitList, habit)
	habitList = append(habitList, habit)

	e.GET("/", func(e echo.Context) error {
		return e.Render(http.StatusOK, "index", &habitList)
	})

	resStruct := TestData{
		Name:  "asdfasdf",
		Phone: "123424",
		Email: "asdf@gmail.com",
	}
	e.GET("/get-info", func(e echo.Context) error {
		return e.Render(http.StatusOK, "name_card", &resStruct)
	})
	e.POST("/clicked", func(e echo.Context) error {
		return e.Render(http.StatusOK, "habits", &habitList)
	})
	
	e.Logger.Fatal(e.Start(":8080"))

}
