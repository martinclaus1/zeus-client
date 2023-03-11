package config

import "time"

var Version string

var BuildTime string

func init() {
	if Version == "" {
		Version = "dev"
	}

	if BuildTime == "" {
		BuildTime = time.Now().Format(time.RFC1123)
	}
}
