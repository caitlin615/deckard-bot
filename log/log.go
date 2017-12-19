// Package log serves as a light wrapper around logrus.
package log

import (
	"github.com/handwritingio/deckard-bot/config"

	"github.com/sirupsen/logrus"
)

func init() {
	if config.RuntimeEnv == "development" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// Alias these functions. More should be added here as needed.
var (
	Error   = logrus.Error
	Errorf  = logrus.Errorf
	Errorln = logrus.Errorln
	Fatal   = logrus.Fatal
	Fatalf  = logrus.Fatalf
	Fatalln = logrus.Fatalln
	Panic   = logrus.Panic
	Panicf  = logrus.Panicf
	Panicln = logrus.Panicln
	Print   = logrus.Print
	Printf  = logrus.Printf
	Println = logrus.Println
	Debug   = logrus.Debug
	Debugf  = logrus.Debugf
	Debugln = logrus.Debugln
	Warn    = logrus.Warn
	Warnf   = logrus.Warnf
	Warnln  = logrus.Warnln
	Info    = logrus.Info
	Infof   = logrus.Infof
	Infoln  = logrus.Infoln
)

// Fields is an alias for logrus.Fields.
type Fields logrus.Fields

// WithFields is an alias for logrus.WithFields.
func WithFields(f Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(f))
}
