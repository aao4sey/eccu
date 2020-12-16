package services

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

func SetLogFilter(debug bool) {
	if debug {
		filter := &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"DEBUG", "INFO", "ERROR"},
			MinLevel: logutils.LogLevel("DEBUG"),
			Writer:   os.Stderr,
		}
		log.SetOutput(filter)
	} else {
		filter := &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"DEBUG", "INFO", "ERROR"},
			MinLevel: logutils.LogLevel("INFO"),
			Writer:   os.Stderr,
		}
		log.SetOutput(filter)
	}
}
