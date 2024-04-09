package habits

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Activity struct {
	Name        string
	Description string
}

type Habit struct {
	ID       int
	Activity Activity
	Date     time.Time
}

func (a Activity) String() string {
	return fmt.Sprintf("%v (%v)", strings.ToUpper(a.Name), a.Description)
}
func (h Habit) String() string {
	return fmt.Sprintf("Activity: %v, Date: %v, ID: %v", h.Activity, h.Date.Format("02-Jan-2006"), h.ID)
}

func Remove(habitlist []Habit, id int) []Habit {
	idx := slices.IndexFunc(habitlist, func(h Habit) bool { return h.ID == id })
	fmt.Printf("idx: %v", idx)
	return slices.Delete(habitlist, idx, idx+1)
}
