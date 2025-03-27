package utils

import "fmt"

var BuildVersion string
var BuildRelease string
var BuildTime string

func GetBuildVersion() string {
	if BuildRelease == "" {
		return fmt.Sprintf("git-%v", BuildVersion)
	} else {
		return fmt.Sprintf("%v (git-%v)", BuildRelease, BuildVersion)
	}
}
