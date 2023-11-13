package main

import (
	"fmt"
	"os"

	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
	"github.com/raphgl/syngit/gitops"
)

var cfg config.Config

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	repos := clients.GetRepos(cfg)

	err = gitops.CreateLocalMirrors(repos, cfg)
	if err != nil {
        fmt.Fprintln(os.Stderr,"create local mirrors: ", err)
	}

	err = gitops.SyncMirrors(repos, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
