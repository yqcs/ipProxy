package log

import (
	"log"
	"os"
	"os/user"
	"time"
)

var (
	fileName   string
	folderPath string
)

func Pr(logPrefix string, text string, a ...interface{}) {
	userInfo, err := user.Current()
	folderPath = userInfo.HomeDir + `\log\`
	fileName = time.Now().Format(`2006-01-02 15-04-05`) + ".log"
	file := folderPath + fileName
	_, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, os.ModePerm)
	}

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetPrefix("[" + logPrefix + "] ")
	log.Println(text, a)
}
