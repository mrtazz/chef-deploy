package version

import (
	"encoding/json"
	"expvar"
	"fmt"
)

var (
	version   = ""
	goversion = ""
)

func init() {
	expvar.Publish("version", GetVersion())
}

// Version struct to get current version from other tools
type Version struct {
	Version   string `json:"version"`
	Goversion string `json:"goversion"`
}

// GetDocoptVersionString returns the version string formatted to be used by
// the docopt parser
func GetDocoptVersionString(cmd string) string {
	return fmt.Sprintf("%s %s built with %s",
		cmd, version, goversion)
}

// GetVersion allows other packages to get version information
func GetVersion() Version {
	return Version{
		version,
		goversion,
	}
}

func (V Version) String() string {
	json, err := json.Marshal(V)
	if err != nil {
		return ""
	}
	return string(json)
}
