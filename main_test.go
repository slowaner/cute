package cute

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

const (
	EnvKeyAllureResultsPath  = "ALLURE_OUTPUT_PATH"   // Indicates the path to the results print folder
	EnvKeyAllureOutputFolder = "ALLURE_OUTPUT_FOLDER" // Indicates the name of the folder to print the results.
)

var (
	testServerAddress  = ""
	testServerHost     = ""
	testServerHostName = ""
	testServerPort     = ""
)

func TestMain(m *testing.M) {
	r := http.NewServeMux()
	r.HandleFunc("/with_body", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testServer := httptest.NewServer(r)
	defer testServer.Close()

	testServerAddress = testServer.URL

	u, err := url.Parse(testServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	testServerHost = u.Host
	testServerHostName = u.Hostname()
	testServerPort = u.Port()

	err = os.Setenv("ALLURE_ISSUE_PATTERN", testServerAddress+"/issue/%s")
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = os.Setenv("ALLURE_TESTCASE_PATTERN", testServerAddress+"/test_case/%s")
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = os.Setenv("ALLURE_LINK_TMS_PATTERN", testServerAddress+"/tms/%s")
	if err != nil {
		log.Fatalln(err)
		return
	}

	outPath := getResultPath()

	err = os.RemoveAll(outPath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	os.Exit(m.Run())
}

func getResultPath() string {
	resultsPathToOutput := os.Getenv(EnvKeyAllureResultsPath)
	outputFolderName := getOutputFolderName()

	if resultsPathToOutput != "" {
		return filepath.Join(resultsPathToOutput, outputFolderName)
	}

	return outputFolderName
}

func getOutputFolderName() string {
	outputFolderName := os.Getenv(EnvKeyAllureOutputFolder)
	if outputFolderName != "" {
		return outputFolderName
	}

	return "allure-results"
}
