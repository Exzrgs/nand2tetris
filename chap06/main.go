package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var symbolAddress map[string]int = map[string]int{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

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

	commandList := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		command := formatText(text)

		commandType := getCommandType(command)
		switch commandType {
		case SYMBOL:
			symbol := getLabelSymbol(command)
			symbolAddress[symbol] = len(commandList)
			continue

		case NO_COMMAND:
			continue

		case INVALID_COMMAND:
			fmt.Printf("invalid command: %s\n", command)
			os.Exit(1)
		}

		commandList = append(commandList, command)
	}

	binaryTextList := make([]string, 0)

	for _, command := range commandList {
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
