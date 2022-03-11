package vmtranslator

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
	This function will accept a path to a dir or file.
	if it is a file then it will translate the file to an asm and output the file in the current dir
	if it is a directory it will output the asm with the name of the directory
*/
func Translate(path string) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.Create(strings.TrimSuffix(fileInfo.Name(), ".vm") + ".asm")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	CompileAllRegex()

	var hack string

	// find all vm files within directory
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal("Failed to read directory")
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".vm" {
				hack += singleFileTranslate(filepath.Join(path, file.Name()))
			}
		}

	} else { //input is a singular vm file
		hack += singleFileTranslate(fileInfo.Name())
	}

	hack += infiniteLoop

	if _, err = output.WriteString(hack); err != nil {
		log.Fatal("There was a fatal error building the asm file.")
	}
}

func singleFileTranslate(file string) (hack string) {
	input, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	hack += "// Code Generated from " + filepath.Base(input.Name()) + "\n// Powered by GO (TM)\n"

	translationUnit := strings.TrimSuffix(filepath.Base(input.Name()), ".vm")
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		if command := parseCommand(line, translationUnit); command != nil {
			hack += TranslateCommand(command)
		}
	}
	return
}
