package pkg

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func Measure(start time.Time, name string) {
	seconds := time.Since(start).Seconds()
	log.WithField("time", seconds).WithField("function", name).Debugln("Execution time tracked.")
}
