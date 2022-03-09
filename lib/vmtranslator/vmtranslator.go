package vmtranslator

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var seeker = 0 // Keeps track of where Translator is up to in output file

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

	// find all vm files within directory
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal("Failed to read directory")
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".vm" {
				Translate(filepath.Join(path, file.Name()))
			}
		}

	} else { //input is a singular vm file
		input, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer input.Close()

		parsedCommands := ParseFile(input, strings.TrimSuffix(filepath.Base(input.Name()), ".vm"))

		hack := "// Code Generated from " + input.Name() + "\n// Powered by GO (TM)\n"
		for _, command := range parsedCommands {
			hackCommand, err := TranslateCommand(command)
			if err != nil {
				log.Fatal(err)
			}
			hack += hackCommand
		}
		hack += infiniteLoop

		var bytes int
		if bytes, err = output.WriteAt([]byte(hack), int64(seeker)); err != nil {
			log.Fatal("There was a fatal error building the asm file.")
		}
		seeker += bytes

	}
}
