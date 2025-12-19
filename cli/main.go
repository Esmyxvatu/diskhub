package main

import (
	"fmt"
	"os"
)

type Project struct {
	Name	string
	Path	string
	Desc	string
	Tags	[]string
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		ShowHelp([]string{})
		return
	}

	switch args[0] {
	case "help":
		ShowHelp(args[1:])
	case "list":
		ListProjects(args[1:])
	case "init":
		InitProject(args[1:])
	default:
		fmt.Printf("The command %s wasn't found\n", args[0])
		fmt.Println("Here's all command available : help, list, show")
		fmt.Printf("Type %s help for more information\n", os.Args[0])
	}
}
