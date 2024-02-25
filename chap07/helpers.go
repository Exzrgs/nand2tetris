package main

import "path/filepath"

func getOutputFileName(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	core := base[0 : len(base)-len(ext)]
	newFileName := core + ".asm"
	return newFileName
}
