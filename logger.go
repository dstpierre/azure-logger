package azurelogger

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const prefix = "al_"

var logfile *os.File

// Start changes output of log to be a file
func Start() error {
	logdir := getApplicationDirectory()
	if _, err := os.Stat(logdir); os.IsNotExist(err) {
		if err := os.Mkdir(logdir, 0666); err != nil {
			return err
		}
	}

	// getting the instance id
	filename := os.Getenv("WEBSITE_INSTANCE_ID")
	if len(filename) <= 6 {
		filename = "local.log"
	} else {
		filename = fmt.Sprintf("%s_%s.log", filename[0:6], time.Now().Format(time.RFC3339))
	}

	f, err := os.OpenFile(path.Join(logdir, prefix+filename), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	logfile = f

	log.SetOutput(logfile)

	purgeFiles()

	return nil
}

// Stop closes the log file
func Stop() {
	if logfile != nil {
		logfile.Close()
	}
}

func getApplicationDirectory() string {
	home := os.Getenv("HOME")
	if len(home) == 0 {
		home = "./"
	} else {
		home += "/LogFiles/Application"
	}

	return home
}

func purgeFiles() {
	filepath.Walk(getApplicationDirectory(), func(p string, f os.FileInfo, err error) error {
		if strings.HasPrefix(f.Name(), prefix) && strings.HasSuffix(f.Name(), ".log") {
			if time.Since(f.ModTime()).Hours() > 24 {
				os.Remove(p)
			}
		}
		return nil
	})
}
