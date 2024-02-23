package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	A_COMMAND       = 0
	C_COMMAND       = 1
	SYMBOL          = 2
	INVALID_COMMAND = 3

	ASSIGN     = "1"
	NOT_ASSIGN = "0"

	CONDITION     = "1"
	NOT_CONDITION = "0"

	COMP_ZERO        = "0101010"
	COMP_ONE         = "0111111"
	COMP_MINUS_ONE   = "0111010"
	COMP_D           = "0001100"
	COMP_Y           = "110000"
	COMP_NOT_D       = "0001101"
	COMP_NOT_Y       = "110001"
	COMP_MINUS_D     = "0001111"
	COMP_MINUS_Y     = "110011"
	COMP_D_PLUS_ONE  = "0011111"
	COMP_Y_PLUS_ONE  = "110111"
	COMP_D_MINUS_ONE = "0001110"
	COMP_Y_MINUS_ONE = "110010"
	COMP_D_PLUS_Y    = "000010"
	COMP_D_MINUS_Y   = "010011"
	COMP_Y_MINUS_D   = "000111"
	COMP_D_AND_Y     = "000000"
	COMP_D_OR_Y      = "010101"

	COMP_A           = "0" + COMP_Y
	COMP_NOT_A       = "0" + COMP_NOT_Y
	COMP_MINUS_A     = "0" + COMP_MINUS_Y
	COMP_A_PLUS_ONE  = "0" + COMP_Y_PLUS_ONE
	COMP_A_MINUS_ONE = "0" + COMP_Y_MINUS_ONE
	COMP_D_PLUS_A    = "0" + COMP_D_PLUS_Y
	COMP_D_MINUS_A   = "0" + COMP_D_MINUS_Y
	COMP_A_MINUS_D   = "0" + COMP_Y_MINUS_D
	COMP_D_AND_A     = "0" + COMP_D_AND_Y
	COMP_D_OR_A      = "0" + COMP_D_OR_Y

	COMP_M           = "1" + COMP_Y
	COMP_NOT_M       = "1" + COMP_NOT_Y
	COMP_MINUS_M     = "1" + COMP_MINUS_Y
	COMP_M_PLUS_ONE  = "1" + COMP_Y_PLUS_ONE
	COMP_M_MINUS_ONE = "1" + COMP_Y_MINUS_ONE
	COMP_D_PLUS_M    = "1" + COMP_D_PLUS_Y
	COMP_D_MINUS_M   = "1" + COMP_D_MINUS_Y
	COMP_M_MINUS_D   = "1" + COMP_Y_MINUS_D
	COMP_D_AND_M     = "1" + COMP_D_AND_Y
	COMP_D_OR_M      = "1" + COMP_D_OR_Y
)

var (
	nowSymbolAddress = 10000
)

func formatText(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\n", "")
	index := strings.Index(text, "//")
	if index != -1 {
		text = text[:index]
	}
	return text
}

func getCommandType(command string) int {
	aCommand := regexp.MustCompile(`^@`)
	cCommand := regexp.MustCompile(`[=;]`)
	symbol := regexp.MustCompile(`^\(.*\)$`)

	if aCommand.MatchString(command) {
		return A_COMMAND
	}
	if cCommand.MatchString(command) {
		return C_COMMAND
	}
	if symbol.MatchString(command) {
		return SYMBOL
	}
	return INVALID_COMMAND
}

func getSymbolAddress(symbol string) int {
	address, exist := symbolAddress[symbol]
	if exist {
		return address
	} else {
		symbolAddress[symbol] = nowSymbolAddress
		nowSymbolAddress += 1
		return symbolAddress[symbol]
	}
}

func getLabelSymbol(command string) string {
	symbol := strings.Trim(command, "()")
	return symbol
}

func getValue(command string) int {
	valueString := command[1:]
	value, err := strconv.Atoi(valueString)
	if err != nil {
		value = getSymbolAddress(valueString)
	}

	return value
}

func formatValue(value int) (string, error) {
	binary := strconv.FormatInt(int64(value), 2)
	if len(binary) > 15 {
		err := fmt.Errorf("value over flow")
		return "", err
	}

	if len(binary) < 15 {
		paddingLen := 15 - len(binary)
		binary = strings.Repeat("0", paddingLen) + binary
	}

	return binary, nil
}

func getMnemonic(command string) (string, string, string, error) {
	index := strings.Index(command, "=")
	if index != -1 && index != 0 && index != len(command)-1 {
		comp := command[index+1:]
		dest := command[:index]
		jump := ""
		return comp, dest, jump, nil
	}

	index = strings.Index(command, ";")
	if index != -1 && index != 0 && index != len(command)-1 {
		comp := command[:index]
		dest := ""
		jump := command[index+1:]
		return comp, dest, jump, nil
	}

	err := fmt.Errorf("invalid ACommand")
	return "", "", "", err
}

func getDest(mnemonic string) (string, error) {
	a := NOT_ASSIGN
	d := NOT_ASSIGN
	m := NOT_ASSIGN
	var err error
	err = nil

	switch mnemonic {
	case "M":
		m = ASSIGN
	case "D":
		d = ASSIGN
	case "MD":
		m = ASSIGN
		d = ASSIGN
	case "A":
		a = ASSIGN
	case "AM":
		a = ASSIGN
		m = ASSIGN
	case "AD":
		a = ASSIGN
		d = ASSIGN
	case "AMD":
		a = ASSIGN
		m = ASSIGN
		d = ASSIGN
	case "":
	default:
		err = fmt.Errorf("invalid before =")
	}

	return a + d + m, err
}

func getComp(mnemonic string) (string, error) {
	var comp string
	var err error
	err = nil

	switch mnemonic {
	case "0":
		comp = COMP_ZERO
	case "1":
		comp = COMP_ONE
	case "-1":
		comp = COMP_MINUS_ONE
	case "D":
		comp = COMP_D
	case "A":
		comp = COMP_A
	case "!D":
		comp = COMP_NOT_D
	case "!A":
		comp = COMP_NOT_A
	case "-D":
		comp = COMP_MINUS_D
	case "-A":
		comp = COMP_MINUS_A
	case "D+1":
		comp = COMP_D_PLUS_ONE
	case "A+1":
		comp = COMP_A_PLUS_ONE
	case "D-1":
		comp = COMP_D_MINUS_ONE
	case "A-1":
		comp = COMP_A_MINUS_ONE
	case "D+A":
		comp = COMP_D_PLUS_A
	case "D-A":
		comp = COMP_D_MINUS_A
	case "A-D":
		comp = COMP_A_MINUS_D
	case "D&A":
		comp = COMP_D_AND_A
	case "D|A":
		comp = COMP_D_OR_A
	case "M":
		comp = COMP_M
	case "!M":
		comp = COMP_NOT_M
	case "-M":
		comp = COMP_MINUS_M
	case "M+1":
		comp = COMP_M_PLUS_ONE
	case "M-1":
		comp = COMP_M_MINUS_ONE
	case "D+M":
		comp = COMP_D_PLUS_M
	case "D-M":
		comp = COMP_D_MINUS_M
	case "M-D":
		comp = COMP_M_MINUS_D
	case "D&M":
		comp = COMP_D_AND_M
	case "D|M":
		comp = COMP_D_OR_M
	default:
		err = fmt.Errorf("invalid comp mnemonic")
	}

	return comp, err
}

func getJump(mnemonic string) (string, error) {
	less := NOT_CONDITION
	equal := NOT_CONDITION
	more := NOT_CONDITION
	var err error
	err = nil

	switch mnemonic {
	case "JGT":
		more = CONDITION
	case "JEQ":
		equal = CONDITION
	case "JGE":
		equal = CONDITION
		more = CONDITION
	case "JLT":
		less = CONDITION
	case "JNE":
		more = CONDITION
		less = CONDITION
	case "JLE":
		less = CONDITION
		equal = CONDITION
	case "JMP":
		less = CONDITION
		equal = CONDITION
		more = CONDITION
	case "":
	default:
		err = fmt.Errorf("invalid after ;")
	}

	return less + equal + more, err
}
