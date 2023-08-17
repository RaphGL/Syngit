package main

import (
	"fmt"
	"os"

	"github.com/raphgl/syngit/clients"
	"github.com/raphgl/syngit/config"
)

var cfg config.Config

func main() {
	cfg, err := config.LoadConfig("syngit.toml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	repos, _ := clients.GetCodebergRepos(cfg)
	for _, r := range repos {
		clients.CloneRepo(&r)
	}
}
