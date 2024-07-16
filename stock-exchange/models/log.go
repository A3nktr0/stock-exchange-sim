package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// type Typelog string

// const namelogfile = "log.txt"

var path = ""

var logfile *os.File

func Create_log_file() {
	namelogfile := strings.Split(os.Args[1], "/")[len(strings.Split(os.Args[1], "/"))-1] + ".log"
	namefile := "../log/" + time.Now().Format("01.02.2006")
	// time := time.Now().Format("01.02.2006_15h04:05_")
	if _, err := os.Stat(namefile); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(namefile, os.ModePerm)
		Check_err(err)
	}
	// path = namefile + "/" + time + namelogfile
	path = namefile + "/" + namelogfile
	newlogfile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	Check_err(err)
	logfile = newlogfile
}

func Inser_log(message string) {
	// message = strings.Split(message, "\n")[0] + "\n"
	fmt.Print(message)
	_, err := logfile.WriteString(message)
	Check_err(err)
}

// func GetHistory() (string, error) {
// 	data, err := os.ReadFile(path)
// 	return string(data), err
// }

func Check_err(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
