package logging

import (
	"log"
	"net/http"
	"os"
	"time"
)

//var logger *log.Logger

/*
Logging represents logging interface
*/
type Logging interface {
	Httplog(http.HandlerFunc) http.HandlerFunc
	Printlog(string, string)
	PrintFatal(string, error)
}

type logging struct {
	logger *log.Logger
}

/*
New logging object
*/
func New(longEnv string) Logging {
	logger := log.New(os.Stdout, longEnv, log.LstdFlags|log.Lshortfile)
	return &logging{logger}
}

/*
Httplog handles how long it takes for a request to process
*/
func (l logging) Httplog(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		defer l.logger.Printf("%s request processed in %s\n", request.URL.Path, time.Now().Sub(startTime))
		next(response, request)
	}
}

/*
Printlog prints log message
*/
func (l logging) Printlog(logType string, logMessage string) {
	l.logger.Printf("%s : %s\n", logType, logMessage)
}

/*
PrintFatal prints a message and exits
*/
func (l logging) PrintFatal(logMessage string, err error) {
	l.logger.Fatalf("%v: %v", logMessage, err)
}
