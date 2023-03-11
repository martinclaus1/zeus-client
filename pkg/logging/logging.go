package logging

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"runtime"
)

// FormatterHook is a hook that writes logs of specified LogLevels with a formatter to specified Writer
type FormatterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

// Fire will be called when some logging function is called with current hook
// It will format log entry and write it to appropriate writer
func (hook *FormatterHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *FormatterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func SetupLogging(debugMode bool) {
	log.SetOutput(logrus.StandardLogger().Writer())
	file, err := os.OpenFile("zeus-client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Info("Failed to logger to file, using default stderr")
	}
	logrus.SetOutput(io.Discard) // Send all logs to nowhere by default
	logrus.SetLevel(logrus.DebugLevel)

	textFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", ""
		},
		DisableLevelTruncation: true,
	}
	logrus.AddHook(&FormatterHook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
		Formatter: textFormatter,
	})

	levels := []logrus.Level{logrus.InfoLevel}
	if debugMode {
		levels = append(levels, logrus.DebugLevel)
	}

	logrus.AddHook(&FormatterHook{
		Writer:    os.Stdout,
		LogLevels: levels,
		Formatter: textFormatter,
	})

	logrus.AddHook(&FormatterHook{
		Writer:    file,
		LogLevels: logrus.AllLevels,
		Formatter: &logrus.JSONFormatter{},
	})

	logrus.SetReportCaller(true)
}
