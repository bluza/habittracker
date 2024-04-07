package habits

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Activity struct {
	Name        string
	Description string
}

type Habit struct {
	Activity Activity
	Date     time.Time
}

func (a Activity) String() string {
	return fmt.Sprintf("%v (%v)", strings.ToUpper(a.Name), a.Description)
}
func (h Habit) String() string {
	return fmt.Sprintf("Activity: %v, Date: %v.", h.Activity, h.Date.Format("02-Jan-2006"))
}

func (h *Habit) save() error {
	filename := h.Activity.Name + h.Date.Format("02-Jan-2006") + ".json"
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0600)
}

func load(filename string) (*Habit, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var habit Habit
	json.Unmarshal(bytes, &habit)
	return &habit, nil
}

// func main() {

// 	a1 := Activity{"meditieren", "tägliche kurze meditations session"}
// 	h1 := Habit{a1, time.Now()}

// 	h1.save()

// 	h2, err := load("meditieren07-Apr-2024.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(h2)
// }
