package azurelogger_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	logger "github.com/dstpierre/azure-logger"
)

func TestLogOutputChangeToFile(t *testing.T) {
	if _, err := os.Stat("al_local.log"); err == nil {
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

func TestSimulateAzureWebApp(t *testing.T) {
	os.Setenv("HOME", "d:/todelete")
	os.Setenv("WEBSITE_INSTANCE_ID", "test1234567")

	filename := fmt.Sprintf("al_test12_%s.log", time.Now().Format("2006-Jan-2-15-04-05"))

	err := logger.Start()
	if err != nil {
		fmt.Printf("Unable to start logger: %s", err.Error())
		t.Fail()
		return
	}

	log.Println("from azure simulated environment")

	logger.Stop()

	if _, err := os.Stat("d:/todelete/LogFiles/Application/" + filename); os.IsNotExist(err) {
		t.Fail()
	}
}

func TestRollOver(t *testing.T) {
	err := os.RemoveAll("d:/todelete/LogFiles/Application")
	if err != nil {
		fmt.Println("Unable to remove all test log files")
		t.Fail()
		return
	}

	os.Setenv("HOME", "d:/todelete")
	os.Setenv("WEBSITE_INSTANCE_ID", "test1234567")

	err = logger.StartWithOptions(3*time.Second, 15*1024, 5*time.Minute)
	if err != nil {
		fmt.Printf("Unable to start logger: %s", err.Error())
		t.Fail()
		return
	}

	start := time.Now()
	i := 0
	for {
		if time.Since(start) > 40*time.Second {
			break
		}
		i++
		log.Printf("wirting %d\n", i)
	}

	logger.Stop()

	files, err := ioutil.ReadDir("d:/todelete/LogFiles/Application")
	if err != nil {
		fmt.Println("Unable to read files from application log directory")
		t.Fail()
		return
	}

	if len(files) < 2 {
		fmt.Printf("There's only %d files", len(files))
		t.Fail()
	}
}
