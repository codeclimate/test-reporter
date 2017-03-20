package version

import (
	"fmt"
)

// Version number of the command (1.x)
var Version = ""

// BuildVersion is the SHA of the git command the binary was built against (ed3dfu)
var BuildVersion = ""

// BuildTime is the time the binary was built
var BuildTime = ""

func FormattedVersion() string {
	if Version == "" {
		return "unknown"
	}
	return fmt.Sprintf("%s (%s @ %s)\n", Version, BuildVersion, BuildTime)
}
