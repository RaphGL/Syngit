package main

import "github.com/raphgl/syngit/clients"

func main() {
	repos, _ := clients.GetCodebergRepos()
	for _, r := range repos {
		clients.CloneRepo(&r)
	}
}
