package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	fmt.Println("vim-go")
	log.Println("logger")
}
