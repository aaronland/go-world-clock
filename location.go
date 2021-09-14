package clock

import (
	"time"
)

// type Location contains structured data about the time in a given timezone
type Location struct {
	// The (string) Who's On First ID associated with the timezone
	Id string
	// The major/minor label for the timezone
	Timezone string
	// The time in that timezone
	Time time.Time
}
