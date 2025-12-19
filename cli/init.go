package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ConfFile struct {
	Name        string            `json:"name"`
	Author      string            `json:"author"`
	Description string            `json:"desc"`
	Status      string            `json:"status"`
	Tags        []string          `json:"tags"`
	Links       map[string]string `json:"links"`
}

func CreateId(length int16) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func InitProject(arguments []string) {
	reader := bufio.NewReader(os.Stdin)
	if len(arguments) < 1 {
		ShowHelp([]string{"init"})
		return
	}

	name := arguments[0]
	err := errors.New("")

	if name != "." {
		err := os.Mkdir(arguments[0], 0755)
		if err != nil {
			fmt.Printf("Error while creating %s folder : %v \n", arguments[0], err)
		}
	} else {
		file, err := os.Open(name)
		if err != nil {
			fmt.Printf("Error while opening %s folder : %v \n", arguments[0], err)
		}

		name = file.Name()
	}

	config := ConfFile{
		Name:        name,
		Author:      "",
		Description: "",
		Status:      "Idea",
		Tags:        []string{},
		Links:       make(map[string]string),
	}

	fmt.Printf("Give a short description of %s: ", name)
	config.Description, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error while reading description: %v \n", err)
	}

	fmt.Printf("Who's the author? ")
	config.Author, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error while reading author name: %v \n", err)
	}

	content, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Error while creating config file: %v \n", err)
	}

	err = os.Mkdir(".diskhub", 0755)
	if err != nil {
		fmt.Printf("Error while creating diskhub folder : %v \n", err)
	}

	file, err := os.Create(".diskhub/conf.json")
	if err != nil {
		fmt.Printf("Error while creating config file : %v \n", err)
	}
	_, err = file.Write(content)
	if err != nil {
		fmt.Printf("Error while writing in config file : %v \n", err)
	}
	file.Close()

	file, err = os.Create(".diskhub/.id")
	if err != nil {
		fmt.Printf("Error while creating id file : %v \n", err)
	}
	id, err := CreateId(16)
	if err != nil {
		fmt.Printf("Error while creating id for project: %v \n", err)
	}
	_, err = file.WriteString(id)
	if err != nil {
		fmt.Printf("Error while writing in id file : %v \n", err)
	}
	file.Close()

	file, err = os.Create(".diskhub/.ignore")
	if err != nil {
		fmt.Printf("Error while creating id file : %v \n", err)
	}
	if err != nil {
		fmt.Printf("Error while creating id for project: %v \n", err)
	}
	_, err = file.WriteString("use_gitignore = true")
	if err != nil {
		fmt.Printf("Error while writing in id file : %v \n", err)
	}
	file.Close()
}
