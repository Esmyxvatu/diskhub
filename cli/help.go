package main

import "fmt"

func ShowHelp(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println("============= Diskhub - CLI ============= \n")
		fmt.Println("Command Available :")
		fmt.Println("	- help (command)	: show this help message, or if given a command, more information about it")
		fmt.Println("	- list <dir> 		: get a list of all projects found in subdirectory of given path")
		fmt.Println("	- show <dir>		: tells you if a project is in given directory and if so, show you all information about")
		fmt.Println("	- init <name>		: create a dir of the given name, and add an example diskhub.toml")
		// fmt.Println("	- stats			: gives you stats about every projects you have, it takes some times to process every projects")

		return
	} else if len(arguments) > 1 {
		fmt.Println("Help command can only takes one (optionnal) parameter, a command name. every other parameters will be ignored\n")
	}

	switch arguments[0] {
	case "list":
		fmt.Println("The `list` command requires one parameter, `dir`. This parameter should only be a path to a directory.")
		fmt.Println("It would check for every subdirectory if there is a `diskhub.toml` file.")
		fmt.Println("If a `diskhub.toml` is found, it would print the project following the pattern under, and pass to the next subdirectory.")
		fmt.Println("Pattern:")
		fmt.Println("\033[1;3mProject Name\033[90m /path/to/project \033[0m")
		fmt.Println("Description of the project, generally a short line")
		fmt.Println("\033[96mEach Tags\033[0m")
	case "show":
		fmt.Println("The `show` command requires one parameter, `dir`. This parameter should only be a path to a directory.")
		fmt.Println("This command will check if a file named `diskhub.toml` is in the given path, and if so, show you all the informations about.")
		fmt.Println("If the file isn't found, it would just tell you to try after creating it.")
		fmt.Println("The information this commands gives you are:")
		fmt.Println("	- The name of the project and it's author")
		fmt.Println("	- The description of the project")
		fmt.Println("	- The tags and it's status")
		fmt.Println("	- If the project has a readme, his readme, else an indication that it doesn't have one")
	case "help":
		fmt.Println("The `help` command takes one optionnal parameters, the name of another commands.")
		fmt.Println("If a command name isn't passed as parameter, it will list every commands and give you a short description of it.")
		fmt.Println("Else it gives you more precision about what the command do, and sometimes even gives an exemple.")
	case "init":
		fmt.Println("The `init` command takes one parameter, the name of the new project.")
		fmt.Println("It will create a folder of the given name and then ask you some question.")
		fmt.Println("With the answer, it will create a `diskhub.toml` file, with basics information and some already configured with your answer.")
	default:
		fmt.Printf("The command %s wasn't found in the list of possibilities. Maybe it's a typo?\n", arguments[0])
		fmt.Println("Anyways, heres the list of supported commands: help, list, show")
	}
}
