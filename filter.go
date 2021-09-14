package clock

// type Filter provides definitions for filtering timezone results
type Filters struct {
	// Zero or more strings to test whether they are contained by a given timezone's longform (major/minor) label.
	Timezones []string
}
