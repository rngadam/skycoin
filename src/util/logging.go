//
// logging helpers for cmd/mesh and others
//

// TODO: move other packages to use this

package util

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	logging "github.com/op/go-logging"
)

const (
	defaultLogFormat = "[%{module}:%{level}] %{message}"
)

// LogConfig logger configurations
type LogConfig struct {
	// for internal usage
	level logging.Level
	// Level convertes to level during initialization
	Level string
	// list of all modules
	Modules []string
	// format
	Format string
	// enable colors
	Colors bool
	// output
	Output io.Writer
}

// TODO:
// DefaultLogConfig vs (DevLogConfig + ProdLogConfig) ?

// DevLogConfig default development config for logging
func DevLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		level:   logging.DEBUG, // int
		Level:   "debug",       // string
		Modules: modules,
		Format:  defaultLogFormat,
		Colors:  true,
		Output:  os.Stdout,
	}
}

// ProdLogConfig Default production config for logging
func ProdLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		level:   logging.ERROR,
		Level:   "error",
		Modules: modules,
		Format:  defaultLogFormat,
		Colors:  false,
		Output:  os.Stdout,
	}
}

// convertes l.Level (string) to l.level (int)
// or panics if l.Level is invalid
func (l *LogConfig) initLevel() {
	level, err := logging.LogLevel(l.Level)
	if err != nil {
		log.Panicf("Invalid -log-level %s: %v", l.Level, err)
	}
	l.level = level
}

// InitLogger initialize logging using this LogConfig;
// it panics if l.Format is invalid or l.Level is invalid
func (l *LogConfig) InitLogger() {
	l.initLevel()

	format := logging.MustStringFormatter(l.Format)
	logging.SetFormatter(format)
	for _, s := range l.Modules {
		logging.SetLevel(l.level, s)
	}
	stdout := logging.NewLogBackend(l.Output, "", 0)
	stdout.Color = l.Colors
	logging.SetBackend(stdout)
}

// DisableLogging disables the logger completely
func DisableLogging() {
	logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
}
