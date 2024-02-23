package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var symbolAddress map[string]int

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("command format is invalid")
		os.Exit(1)
	}

	filePath := args[1]
	ext := filepath.Ext(filePath)
	if ext != ".asm" {
		fmt.Printf("file ext must be *.asm but got %s\n", ext)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer file.Close()

	binaryTextList := make([]string, 0)
	symbolAddress = make(map[string]int, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		command := formatText(text)

		commandType := getCommandType(command)
		var binaryText string
		switch commandType {
		case A_COMMAND:
			binaryText, err = getACommandBinary(command)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case C_COMMAND:
			binaryText, err = getCCommandBinary(command)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case SYMBOL:
			symbol := getLabelSymbol(command)
			symbolAddress[symbol] = len(binaryTextList)
			continue

		case NO_COMMAND:
			continue

		case INVALID_COMMAND:
			fmt.Printf("invalid command: %s\n", command)
			os.Exit(1)
		}

		binaryTextList = append(binaryTextList, binaryText)
	}

	binaryFileName := getBinaryFileName(filePath)

	binaryFile, err := os.Create(binaryFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, txt := range binaryTextList {
		_, err := binaryFile.WriteString(txt + "\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
