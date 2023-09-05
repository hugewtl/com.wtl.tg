package main

import (
	"log"
	"os"
)

func initLogFile(filename string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		logfile.Close()
	}()
	// multiWriter := io.MultiWriter(os.Stdout, logfile)

	log.SetOutput(logfile)

}
