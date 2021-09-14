package clock

import (
	"fmt"
	"time"
)

type Location struct {
	Id       string
	Timezone string
	Time     time.Time
}

func (l *Location) String() string {
	return fmt.Sprintf("%s\t%s\t\t%s", l.Id, l.Timezone, l.Time.Format(time.RFC3339))
}
