package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// <- golang module base
	"github.com/dfense/appchallenge/service"
	log "github.com/sirupsen/logrus" // common logger library
)

// TODO in lieu of a more elaborate *.toml config, etcd, consul configuration
const (
	DefaultPort     = ":8082"
	shutdownSeconds = 5
)

// main() - entry point into program
func main() {

	// print welcome
	printHeader()

	// parse command line
	logLevelString := flag.String("loglevel", "info", "loglevels")
	logOutput := flag.String("logoutput", "stdout", `output destination for log info ["file" or "stdout"]`)
	flag.Parse()

	setupLoggerDefaults(*logLevelString, *logOutput)

	// create a stop channel, and setup a signal handler
	stop := make(chan os.Signal, 1)                                    // <- simple channel to handle interrupt signal
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // <- send stop channel a notify event

	// see if environment variable for port was set. if not use default
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = DefaultPort
	}

	// create new http server and endpoints
	srv := service.NewHttpService(addr)

	// block on channel until a signal is received
	<-stop
	log.Info("server going to shut down")
	// create a context to server shutdown call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*shutdownSeconds)
	defer cancel()

	// call shutdown with new created context
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}

	log.Info("Thank you for the opportunity to take your test at apple cloud !!")

}

func printHeader() {
	fmt.Println("\n\n===============================================================")
	fmt.Println("  http service to retrieve chuck norris jokes and replace names")
	fmt.Println("  appchallenge test for John Frailey   ")
	fmt.Println("===============================================================")
	fmt.Println("")
}

// <--- everything here down is just setting up log preferences --->

// setupLoggerDefaults - set preferred logrus defaults
// TODO add file, network, rollover, etc and isolate into log.go
func setupLoggerDefaults(logLevelString, logoutput string) {

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	log.Infof("setting log level to %s", logLevelString)
	loglevel, err := stringToLevel(logLevelString)
	if err != nil {
		// TODO print all level options or refer to api doc
		log.Error(err)
	}

	// only SetOutput if specified for file. else stdout default
	if logoutput == "file" {
		var filename string = "logfile.log"
		// Create the log file if doesn't exist. And append to it if it already exists.
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Error("can not open file for logging, using stdout instead")
		} else {
			log.SetOutput(f)
		}
	}
	log.SetLevel(loglevel)
}

//convert string value of log level to the actual logrus.Level type
func stringToLevel(level string) (log.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return log.DebugLevel, nil
	case "info":
		return log.InfoLevel, nil
	case "error":
		return log.ErrorLevel, nil
	case "warning", "warn":
		return log.WarnLevel, nil
	case "fatal":
		return log.FatalLevel, nil
	case "panic":
		return log.PanicLevel, nil
	default:
		return log.InfoLevel, errors.New("level given wasn't valid")
	}
}
