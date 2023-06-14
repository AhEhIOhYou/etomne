package logger

import (
	"log"
	"os"
)

const (
	Error = iota
	Info
	Debug
)

var tagText = map[int]string{
	Error: "ERROR",
	Info:  "INFO",
	Debug: "DEBUG",
}

func TagText(tag int) string {
	return tagText[tag]
}

func WriteLog(tag int, msg string) {
	loggerPath := os.Getenv("LOGS_PATH")

	if _, err := os.Stat(loggerPath); os.IsNotExist(err) {
		os.MkdirAll(loggerPath, os.ModePerm)
	}
	f, err := os.OpenFile(loggerPath + "/main_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	defer f.Close()

	switch tag {
	case 0:
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
		newLog.Output(2, "\t"+msg)
	case 1:
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime)
		newLog.Output(2, "\t"+msg)
	case 2:
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
		newLog.Output(2, msg)
	}
}

func WriteLogToFile(tag int, msg string, fileName string) {
	loggerPath := os.Getenv("LOGS_PATH")

	f, err := os.OpenFile(loggerPath + "/"+fileName+"_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	defer f.Close()

	newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
	newLog.Output(2, msg)
}
