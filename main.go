package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
	"github.com/raphgl/syngit/gitops"
)

const HELP_MSG = `usage: syngit <command>

Commands:
    init        create a configuration file
    run         synchronize configured repositories
    daemon      run syngit as a daemon, making it run every x amount of time
`

var cfg *config.Config

func synchronizeRepos() {
	repos := clients.GetRepos(cfg)

	err := gitops.CreateLocalMirrors(repos, cfg)
	if err != nil {
		slog.Error("Creating local mirrors", err)
		return
	}

	err = gitops.SyncMirrors(repos, cfg)
	if err != nil {
		slog.Error("Synchronizing mirrors", err)
		return
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(HELP_MSG)
		return
	}

	// setup logging
	logPath, err := config.GetLogFilePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logFile.Close()

	multiOut := io.MultiWriter(logFile, os.Stdout)
	logger := slog.New(slog.NewTextHandler(multiOut, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	// load config
	cfg, err = config.LoadConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// parse commands
	switch os.Args[1] {
	case "init":
		fmt.Println("TODO: this functionality has not yet been implemented.")

	case "run":
		synchronizeRepos()

	case "daemon":
		for {
			timer := time.NewTimer(1 * time.Hour)
			synchronizeRepos()
			<-timer.C
		}

	default:
		fmt.Printf(HELP_MSG)
	}

}
