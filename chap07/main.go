package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("command format is invalid")
		return
	}

	gotPath := args[1]

	f, err := os.Open(gotPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := f.Name()
	ext := filepath.Ext(gotPath)
	coreName := fileName[0 : len(fileName)-len(ext)]

	commandList := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		command := formatText(text)

		if command == "" {
			continue
		}

		commandType, err := getCommandType(command)
		if err != nil {
			fmt.Println(err)
			return
		}

		var commands []string
		switch commandType {
		case C_PUSH:
			commands, err = parsePush(command, coreName)
		case C_POP:
			commands, err = parsePop(command, coreName)
		case C_ARITHMETIC:
			commands, err = parseArithmetic(command)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		commandList = append(commandList, commands...)
	}

	newFileName := getOutputFileName(f.Name())
	newFile, err := os.Create(newFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cmd := range commandList {
		_, err := newFile.WriteString(cmd + "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
