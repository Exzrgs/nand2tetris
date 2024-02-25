package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	C_PUSH       = iota
	C_POP        = iota
	C_ARITHMETIC = iota

	C_ADD = "add"
	C_SUB = "sub"
	C_NEG = "neg"
	C_EQ  = "eq"
	C_GT  = "gt"
	C_LT  = "lt"
	C_AND = "and"
	C_OR  = "or"
	C_NOT = "not"

	C_ARGUMENT = "argument"
	C_LOCAL    = "local"
	C_STATIC   = "static"
	C_CONSTANT = "constant"
	C_THIS     = "this"
	C_THAT     = "that"
	C_POINTER  = "pointer"
	C_TEMP     = "temp"

	S_ARGUMENT = "ARG"
	S_LOCAL    = "LCL"
	S_THIS     = "THIS"
	S_THAT     = "THAT"
	S_POINTER  = "3"
	S_TEMP     = "5"

	STORE_M_TO_D = "D=M"
	STORE_D_TO_M = "M=D"
	STORE_A_TO_D = "D=A"
	PUSH_D       = "@SP\nA=M\nM=D\n@SP\nM=M+1"
	POP_D        = "@SP\nM=M-1\nA=M\nD=M"
	POP_A        = "@SP\nM=M-1\nA=M\nA=M"
	ADD_DA_TO_D  = "D=D+A"
	SUB_DA_TO_D  = "D=D-A"
	SUB_AD_TO_D  = "D=A-D"
	NEG_D        = "D=-D"
	AND_DA_TO_D  = "D=D&A"
	OR_DA_TO_D   = "D=D|A"
	NOT_D        = "D=!D"
	GET_TMP_TO_A = "@tmp\nA=M"
)

var (
	segmentMap map[string]string = map[string]string{
		C_ARGUMENT: S_ARGUMENT,
		C_LOCAL:    S_LOCAL,
		C_THIS:     S_THIS,
		C_THAT:     S_THAT,
		C_POINTER:  S_POINTER,
		C_TEMP:     S_TEMP,
	}

	POP_DA = POP_D + "\n" + POP_A
)

func setACommand(num string) string {
	return "@" + num
}

func storeIndexAddressToTmp(segment string, index string) string {
	return fmt.Sprintf("@%s\nD=M\n@%s\nD=D+A\n@tmp\nM=D", segment, index)
}

func getIndexAddress(segment string, index string) string {
	segmentInt, _ := strconv.Atoi(segment)
	indexInt, _ := strconv.Atoi(index)

	return fmt.Sprintf("@%d", segmentInt+indexInt)
}

func getComparisonCommand(comparison string) string {
	labe1YES := "L" + uuid.New().String()
	labelNO := "L" + uuid.New().String()
	return fmt.Sprintf("@%s\n%s\nD=0\n@%s\n0;JMP\n(%s)\nD=-1\n(%s)", labe1YES, comparison, labelNO, labe1YES, labelNO)
}

func getStaticAddressCommand(fileName string, index string) string {
	return "@" + fileName + "." + index
}

func formatText(text string) string {
	text = strings.ReplaceAll(text, "\n", "")
	index := strings.Index(text, "//")
	if index != -1 {
		text = text[:index]
	}
	return text
}

func getCommandType(command string) (int, error) {
	words := strings.Fields(command)

	switch words[0] {
	case "push":
		return C_PUSH, nil
	case "pop":
		return C_POP, nil
	case C_ADD, C_SUB, C_NEG, C_EQ, C_GT, C_LT, C_AND, C_OR, C_NOT:
		return C_ARITHMETIC, nil
	default:
		return -1, fmt.Errorf("invalid command: %s", command)
	}
}

func parsePush(command string, fileName string) ([]string, error) {
	words := strings.Fields(command)

	if len(words) != 3 {
		return nil, fmt.Errorf("invalid push command %s", command)
	}

	segment := words[1]
	index := words[2]

	var address string
	switch segment {
	case C_ARGUMENT, C_LOCAL, C_THIS, C_THAT:
		segmentSymbol := segmentMap[segment]
		address = storeIndexAddressToTmp(segmentSymbol, index)
		tmp := GET_TMP_TO_A
		store := STORE_M_TO_D
		push := PUSH_D
		return []string{address, tmp, store, push}, nil
	case C_POINTER, C_TEMP:
		segmentSymbol := segmentMap[segment]
		address = getIndexAddress(segmentSymbol, index)
	case C_STATIC:
		address = getStaticAddressCommand(fileName, index)
	case C_CONSTANT:
		return []string{setACommand(index), STORE_A_TO_D, PUSH_D}, nil
	default:
		return nil, fmt.Errorf("invalid segment %s", segment)
	}

	store := STORE_M_TO_D
	push := PUSH_D
	return []string{address, store, push}, nil
}

func parsePop(command string, name string) ([]string, error) {
	words := strings.Fields(command)

	if len(words) != 3 {
		return nil, fmt.Errorf("invalid push command %s", command)
	}

	segment := words[1]
	index := words[2]

	var address string
	switch segment {
	case C_ARGUMENT, C_LOCAL, C_THIS, C_THAT:
		segmentSymbol := segmentMap[segment]
		address = storeIndexAddressToTmp(segmentSymbol, index)
		pop := POP_D
		tmp := GET_TMP_TO_A
		store := STORE_D_TO_M
		return []string{address, pop, tmp, store}, nil
	case C_POINTER, C_TEMP:
		segmentSymbol := segmentMap[segment]
		pop := POP_D
		address := getIndexAddress(segmentSymbol, index)
		store := STORE_D_TO_M
		return []string{pop, address, store}, nil
	case C_STATIC:
		pop := POP_D
		address = getStaticAddressCommand(name, index)
		store := STORE_D_TO_M
		return []string{pop, address, store}, nil
	case C_CONSTANT:
		return []string{POP_D}, nil
	default:
		return nil, fmt.Errorf("invalid segment %s", segment)
	}
}

func parseArithmetic(command string) ([]string, error) {
	switch command {
	case C_ADD:
		return []string{POP_DA, ADD_DA_TO_D, PUSH_D}, nil
	case C_SUB:
		return []string{POP_DA, SUB_AD_TO_D, PUSH_D}, nil
	case C_NEG:
		return []string{POP_D, NEG_D, PUSH_D}, nil
	case C_EQ:
		EQ_D := getComparisonCommand("D;JEQ")
		return []string{POP_DA, SUB_AD_TO_D, EQ_D, PUSH_D}, nil
	case C_GT:
		GT_D := getComparisonCommand("D;JGT")
		return []string{POP_DA, SUB_AD_TO_D, GT_D, PUSH_D}, nil
	case C_LT:
		LT_D := getComparisonCommand("D;JLT")
		return []string{POP_DA, SUB_AD_TO_D, LT_D, PUSH_D}, nil
	case C_AND:
		return []string{POP_DA, AND_DA_TO_D, PUSH_D}, nil
	case C_OR:
		return []string{POP_DA, OR_DA_TO_D, PUSH_D}, nil
	case C_NOT:
		return []string{POP_D, NOT_D, PUSH_D}, nil
	default:
		return nil, fmt.Errorf("invalid arithmetic command %s", command)
	}
}
