package main

import (
	"fmt"
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
		fmt.Fprintln(os.Stderr, "create local mirrors: ", err)
	}

	err = gitops.SyncMirrors(repos, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(HELP_MSG)
		return
	}

	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// parse commands
	switch os.Args[1] {
	case "init":

	case "run":
		synchronizeRepos()

	case "daemon":
		done := make(chan struct{})
		go func() {
			defer func() { done <- struct{}{} }()
			for {
				timer := time.NewTimer(1 * time.Hour)
				synchronizeRepos()
				<-timer.C
			}
		}()
		<-done
	}

}
