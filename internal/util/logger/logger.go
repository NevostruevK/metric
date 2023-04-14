package logger

import (
	"log"
	"os"
)

func LogError(name, obj string, err error) {
	logger := log.New(os.Stdout, name, log.LstdFlags)
	if err != nil {
		logger.Printf("ERROR : %s returned the error : %v", obj, err)
	} else {
		logger.Printf("ERROR : %s ", obj)
	}
}
