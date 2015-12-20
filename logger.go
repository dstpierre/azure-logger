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

const (
	prefix         = "al_"
	defaultMaxSize = 200 * 1024
)

var active = false
var d time.Duration
var maxsize int64
var retentionDuration time.Duration
var logfile *os.File
var t *time.Ticker

// Start the outpug log to file with default options
func Start() error {
	return StartWithOptions(time.Hour, defaultMaxSize, 24*time.Hour)
}

// StartWithOptions can be used to change default option
func StartWithOptions(md time.Duration, mfs int64, rd time.Duration) error {
	active = true

	d = md
	maxsize = mfs
	retentionDuration = rd

	purgeFiles()

	err := swapFile()
	go monitor()

	return err
}

// Stop closes the log file
func Stop() {
	active = false

	log.SetOutput(os.Stderr)

	if t != nil {
		t.Stop()
		t = nil
	}

	if logfile != nil {
		logfile.Close()
		logfile = nil
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

func swapFile() error {
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
		filename = fmt.Sprintf("%s_%s.log", filename[0:6], time.Now().Format("2006-Jan-2-15-04-05"))
	}

	f, err := os.Create(path.Join(logdir, prefix+filename))
	if err != nil {
		return err
	}

	logfile = f

	log.SetOutput(logfile)

	return nil
}

func rollover() {
	if logfile == nil {
		return
	}

	s, err := logfile.Stat()
	if err != nil || s == nil {
		Stop()
	}

	kb := s.Size() / 1024
	if kb > maxsize {
		log.SetOutput(os.Stderr)

		logfile.Close()
		logfile = nil

		swapFile()
	}
}

func purgeFiles() {
	filepath.Walk(getApplicationDirectory(), func(p string, f os.FileInfo, err error) error {
		if f != nil && strings.HasPrefix(f.Name(), prefix) && strings.HasSuffix(f.Name(), ".log") {
			if time.Since(f.ModTime()) > retentionDuration {
				os.Remove(p)
			}
		}
		return nil
	})
}

func monitor() {
	t = time.NewTicker(d)
	for {
		<-t.C
		if active {
			rollover()
		}
		purgeFiles()
	}
}
