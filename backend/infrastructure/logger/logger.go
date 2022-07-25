package logger

import (
	"log"
	"os"
)

const (
	Error = 0
	Info  = 1
	Debug = 2
	Fatal = 3
)

var tagText = map[int]string{
	Error: "ERROR",
	Info:  "INFO",
	Debug: "DEBUG",
	Fatal: "FATAL",
}

func TagText(tag int) string {
	return tagText[tag]
}

func WriteLog(tag int, msg string) {

	f, err := os.OpenFile("logs/main_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	defer f.Close()

	switch tag {
	case 0:
		// Уровень Error - название файла и номер строки конкретного места ошибки
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
		newLog.Output(2, "\t"+msg)
	case 1:
		// Уровень Info - полезная информация
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime)
		newLog.Output(2, "\t"+msg)
	case 2:
		// Уровень Debug - временная отладочная информация
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
		newLog.Output(2, "\t"+msg)
	case 3:
		// Уровень Fatal - запись и экстренное завершение работы
		newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
		newLog.Output(2, "\t"+msg)
		os.Exit(420)
	}
}

func WriteLogToFile(tag int, msg string, fileName string) {

	f, err := os.OpenFile("logs/"+fileName+"_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	defer f.Close()

	newLog := log.New(f, TagText(tag)+"\t", log.Ldate|log.Ltime|log.Lshortfile)
	newLog.Output(2, "\t"+msg)
}
