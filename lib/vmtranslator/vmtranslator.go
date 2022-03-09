package vmtranslator

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
	This function will accept a path to a dir or file.
	if it is a file then it will translate the file to an asm and output the file in the current dir
	if it is a directory it will output an asm for each file in the dir
*/
func Translate(path string) {
	output_file_name := strings.Split(path, ".")[0] + ".asm"
	seeker := 0

	output, err := os.Create(output_file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	// find all vm files within directory
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal("Failed to read directory")
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".vm" {
				temp := filepath.Join(path, file.Name())
				Translate(temp)
			}
		}

	} else { //input is a singular vm file
		input, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer input.Close()

		translationUnit := fileBaseNameNoExt(input)
		parsedCommands := ParseFile(input, translationUnit)

		hack := "// Code Generated from " + translationUnit + ".vm\nPowered by GO (TM)\n"
		for _, command := range parsedCommands {
			hackCommand, err := TranslateCommand(command)
			if err != nil {
				log.Fatal(err)
			}
			hack += hackCommand
		}

		seeker, err = output.WriteAt([]byte(hack), int64(seeker))
		if seeker, err = output.WriteAt([]byte(hack), int64(seeker)); err != nil {
			log.Fatal("There was a fatal error building the asm file.")
		}

	}
}

func fileBaseNameNoExt(path *os.File) string {
	return strings.TrimSuffix(filepath.Base(path.Name()), ".vm")
}
