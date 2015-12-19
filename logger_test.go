package azurelogger_test

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	logger "github.com/dstpierre/azure-logger"
)

func TestLogOutputChangeToFile(t *testing.T) {
	if _, err := os.Stat("test.log"); err == nil {
		os.Remove("al_local.log")
	}

	err := logger.Start()
	if err != nil {
		t.Fail()
		return
	}

	log.Println("working")

	logger.Stop()

	if _, err := os.Stat("al_local.log"); os.IsNotExist(err) {
		t.Fail()
	}
}

func TestPurgeOldFiles(t *testing.T) {
	if _, err := os.Stat("al_local.log"); os.IsNotExist(err) {
		t.Fail()
		return
	}

	in, err := os.Open("al_local.log")
	if err != nil {
		t.Fail()
		return
	}
	defer in.Close()

	out, err := os.Create("al_old.log")
	if err != nil {
		t.Fail()
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		t.Fail()
	}

	err = out.Close()
	if err != nil {
		t.Fail()
		return
	}

	oldTime := time.Now().Add(-48 * time.Hour)
	os.Chtimes("al_old.log", oldTime, oldTime)

	err = logger.Start()
	if err != nil {
		t.Fail()
	}
	defer logger.Stop()

	log.Println("testing delete of old files")

	if _, err := os.Stat("al_old.log"); err == nil {
		t.Fail()
	}
}
