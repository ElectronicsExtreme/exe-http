package exeserver

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const ()

var (
	errorChannel       chan string
	requestChannel     chan string
	transactionChannel chan string
)

type LogInfo struct {
	API         string
	Path        string
	QueryString string
	RefCode     string
	Body        string
	HTTPStatus  int
	Method      string
	ReqURL      string
	ch          chan string
}

func StartLogger(path string) {
	errorChannel = make(chan string)
	requestChannel = make(chan string)
	transactionChannel = make(chan string)
	go logListener(errorChannel, path+"/error")
	go logListener(requestChannel, path+"/request")
	go logListener(transactionChannel, path+"/transaction")
}

func logListener(ch chan string, path string) {
	err := createDirIfNotExist(path)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		logMsg := <-ch
		currentTime := time.Now()
		filename := fmt.Sprintf("%v/%04d-%02d.log", path, currentTime.Year(), currentTime.Month())
		outfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			log.Println("can't open file", filename)
			log.Println(logMsg)
			continue
		}
		logger := log.New(outfile, "", log.LstdFlags)
		logger.Println(logMsg)
	}
}

func (self *LogInfo) Write() {
	_, path, lineNumber, _ := runtime.Caller(1)
	paths := strings.Split(path, "/")
	filename := fmt.Sprintf("%v(%v)", paths[len(paths)-1], lineNumber)
	logString := fmt.Sprintf("RefCode:%s API:%s File:%s", self.RefCode, self.API, filename)
	if self.Path != "" {
		logString += " Path:" + self.Path
	}
	if self.QueryString != "" {
		logString += " QueryString:" + self.QueryString
	}
	if self.ReqURL != "" {
		logString += " ReqURL:" + self.ReqURL
	}
	if self.Method != "" {
		logString += " Method:" + self.Method
	}
	if self.HTTPStatus != 0 {
		logString += fmt.Sprintf(" HTTPStatus:%v", self.HTTPStatus)
	}
	logString += " Body:" + self.Body
	self.ch <- logString
}

func createDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
