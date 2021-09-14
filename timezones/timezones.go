// package timezones provides methods for loading timezone information derived from Who's On First data.
package timezones

import (
	"embed"
)

//go:embed timezones.csv
var FS embed.FS
