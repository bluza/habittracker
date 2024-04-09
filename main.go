package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"

	"example.com/habits"
)

const DB_FILE_NAME = "foo.sqlite"

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
	CREATE TABLE IF NOT EXISTS habit (
		id integer primary key,
		name text,
		description text
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
func insertHabit(db *sql.DB, habit habits.Habit) {
	stmt, err := db.Prepare("INSERT INTO habit(name, description) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(habit.Activity.Name, habit.Activity.Description)
	if err != nil {
		log.Fatal(err)
	}
}

func loadHabits(db *sql.DB) []habits.Habit {

	rows, err := db.Query("SELECT id, name, description FROM habit")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	habitList := make([]habits.Habit, 0)
	for rows.Next() {
		var tmpHabit habits.Habit
		err := rows.Scan(&tmpHabit.ID, &tmpHabit.Activity.Name, &tmpHabit.Activity.Description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tmpHabit)
		habitList = append(habitList, tmpHabit)
	}
	return habitList
}
func deleteHabit(db *sql.DB, habitList []habits.Habit, idToDelete int) []habits.Habit {
	stmt, err := db.Prepare("DELETE FROM habit WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(idToDelete)
	if err != nil {
		log.Fatal(err)
	}
	return habits.Remove(habitList, idToDelete)
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sql.Open("sqlite3", "./"+DB_FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = os.Stat(DB_FILE_NAME)

	if err != nil {
		createTable(db)
		for _, habit := range generateTestData() {
			insertHabit(db, habit)
		}
	}

	habitList := loadHabits(db)

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
		newHabit := habits.Habit{
			Date: time.Now(),
			Activity: habits.Activity{
				Name:        e.FormValue("Name"),
				Description: e.FormValue("Description"),
			},
		}
		insertHabit(db, newHabit)
		habitList = append(habitList, newHabit)
		return e.Render(http.StatusOK, "habits", &habitList)
	})

	e.DELETE("/habit/:id", func(e echo.Context) error {
		id, err := strconv.Atoi(e.Param("id"))
		if err == nil {
			habitList = deleteHabit(db, habitList, id)
		}
		return e.Render(http.StatusOK, "habits", &habitList)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
