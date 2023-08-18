package main

import (
	"fmt"
	"os"

	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

var cfg config.Config

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	repos := clients.GetRepos(cfg)
	for _, mirrors := range repos {
		for _, m := range mirrors {
			clients.CloneRepo(m, cfg)
		}
	}
}
