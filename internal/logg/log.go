package logg

import (
	"fmt"
	"log"
	"os"
)

var logger = log.Default()

func Errorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func Info(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}

func Fatal(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Fatal]: %s\n", msg)
	os.Exit(1)
}
