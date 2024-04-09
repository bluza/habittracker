package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"

	"example.com/habits"
)

const DB_FILE_NAME = "foo.sqlite"

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

func generateTestData() []habits.Habit {

	habitList := make([]habits.Habit, 0)

	habit := habits.Habit{
		Activity: habits.Activity{Name: "journaling", Description: "write something about your day"},
		Date:     time.Now(),
	}
	habit1 := habits.Habit{
		Activity: habits.Activity{Name: "gym", Description: "build a the body you want"},
		Date:     time.Now().Add(-24 * time.Hour),
	}

	h3 := habits.Habit{Activity: habits.Activity{Name: "meditation", Description: "asdf"}, Date: time.Now()}
	habitList = append(habitList, habit)
	habitList = append(habitList, habit1)
	habitList = append(habitList, h3)
	return habitList
}
func createTable(db *sql.DB) {

	sqlStmt := `
	create table if not exists habit (
		id integer primary key,
		name text,
		description text
		timestamp int
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
func insertHabit(db *sql.DB, habit habits.Habit) {
	stmt, err := db.Prepare("INSERT INTO habit(name, description, timestamp) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(habit.Activity.Name, habit.Activity.Description, habit.Date.Unix())
	if err != nil {
		log.Fatal(err)
	}
}

func loadHabits(db *sql.DB) []habits.Habit {

	rows, err := db.Query("select id, name from habit")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	habitList := make([]habits.Habit, 0)
	for rows.Next() {
		var tmpHabit habits.Habit
		err := rows.Scan(&tmpHabit.ID, &tmpHabit.Activity.Name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tmpHabit)
		habitList = append(habitList, tmpHabit)
	}
	return habitList
}
func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sql.Open("sqlite3", "./"+DB_FILE_NAME+"?cache=shared")
	db.SetMaxOpenConns(2)
	if err != nil {
		log.Fatal(err)
		log.Printf("Created DB file %v", DB_FILE_NAME)
	}
	defer db.Close()

	createTable(db)
	for _, habit := range generateTestData() {
		insertHabit(db, habit)
	}

	habitList := loadHabits(db)
	nextHabitID := len(habitList) - 1

	e := echo.New()
	e.Static("/dist", "dist")
	e.Debug = true
	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e.Renderer = t

	e.GET("/", func(e echo.Context) error {
		return e.Render(http.StatusOK, "index", &habitList)
	})
	
	e.POST("/clicked", func(e echo.Context) error {
		return e.Render(http.StatusOK, "habits", nil)
	})
	e.POST("add", func(e echo.Context) error {
		nextHabitID += 1
		Name := e.FormValue("Name")
		Description := e.FormValue("Description")
		habitList = append(
			habitList,
			habits.Habit{
				Date: time.Now(),
				ID:   nextHabitID,
				Activity: habits.Activity{
					Name:        Name,
					Description: Description,
				},
			})
		return e.Render(http.StatusOK, "habits", &habitList)
	})

	e.DELETE("/habit/:id", func(e echo.Context) error {
		id, err := strconv.Atoi(e.Param("id"))
		if err == nil {
			habitList = habits.Remove(habitList, id)
		}
		return e.Render(http.StatusOK, "habits", &habitList)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
