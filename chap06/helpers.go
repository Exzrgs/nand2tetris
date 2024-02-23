package main

import (
	"path/filepath"
)

const (
	BINARY_A_COMMAND = "0"
	BINARY_C_COMMAND = "111"
)

func getCCommandBinary(command string) (string, error) {
	compMnemonic, destMnemonic, jumpMnemonic, err := getMnemonic(command)
	if err != nil {
		return "", err
	}

	comp, err := getComp(compMnemonic)
	if err != nil {
		return "", err
	}
	dest, err := getDest(destMnemonic)
	if err != nil {
		return "", err
	}
	jump, err := getJump(jumpMnemonic)
	if err != nil {
		return "", err
	}
	binary := BINARY_C_COMMAND + comp + dest + jump
	return binary, nil
}

func getACommandBinary(command string) (string, error) {
	value := getValue(command)
	valueString, err := formatValue(value)
	if err != nil {
		return "", err
	}
	binary := BINARY_A_COMMAND + valueString
	return binary, nil
}

func getBinaryFileName(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	core := base[0 : len(base)-len(ext)]
	binaryFileName := core + ".hack"
	return binaryFileName
}
