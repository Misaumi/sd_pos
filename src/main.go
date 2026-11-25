package main

import "os"

func main() {
	f := buildFacade()
	if len(os.Args) > 1 && os.Args[1] == "demo" {
		RunDemo(f, os.Stdout)
		return
	}
	NewConsole(f, os.Stdin, os.Stdout).Run()
}
