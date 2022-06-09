package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	fe "jack-jack-go/lib/syntaxAnalyzer"

	// "jack-jack-go/lib/util"
	co "jack-jack-go/lib/vmwriter"

	be "jack-jack-go/lib/vmtranslator"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// tokenizes and parses a single jackfile, optionally writing the xml files based on flags
func tokenizeAndParse(jackFile string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(jackFile)
	if err != nil {
		log.Print(fmt.Sprint("Could not open file ", file))
	}
	name := strings.TrimSuffix(filepath.Base(file.Name()), ".jack")
	reader := bufio.NewReader(file)
	tokens := fe.Tokenize(reader)
	if tokenXML {
		writeOptionalFile(name, tokens)
	}
	ts := fe.TS{Tokens: tokens, File: jackFile}
	ast := fe.Parse(&ts)
	if astXML {
		writeOptionalFile(name, ast)
	}
	vmFile, err := os.Create(filepath.Join(buildDirName, name+".vm"))
	if err != nil {
		log.Fatal(err)
	}
	defer vmFile.Close()
	w := bufio.NewWriter(vmFile)
	compiler := co.NewJackCompiler(ast, name, w)
	err = compiler.CompileClass()
	if err != nil {
		log.Print(err)
	}
	w.Flush()
}

func writeOptionalFile(name string, data interface{}) {
	optionFile, err := os.Create(filepath.Join(buildDirName, name+"_TOKENS.xml"))
	if err != nil {
		log.Fatal(err)
	}
	defer optionFile.Close()
	w := bufio.NewWriter(optionFile)
	switch d := data.(type) {
	case []fe.Token:
		err = fe.TokenToXML(d, w)
	case *fe.Node:
		err = fe.NodeToXML(d, w)
	}
	if err != nil {
		log.Print(fmt.Sprint("Could not write token XML for ", name, ".jack"))
	}
	w.Flush()
}

var (
	tokenXML, astXML bool
    buildDirName string
)

// takes command line args and returns all files inside
func expandDirectories(jackFiles []string) (files []string, parent string) {
	for _, arg := range jackFiles {
		fileInfo, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			parent = fileInfo.Name()
			filepath.Walk(arg, func(path string, info os.FileInfo, _ error) error {
				if info.IsDir() {
					// ensures that only goes 1 layer down
					return nil
				}
				if strings.HasSuffix(path, ".jack") {
					files = append(files, path)
				}
				return nil
			})
		} else {
			parent, _ = filepath.Split(fileInfo.Name())
			files = append(files, arg)
		}
	}
	return files, parent
}

// the main function for the compiler, the entry point to the program
func main() {
	flag.BoolVar(&tokenXML, "tokens", false, "output tokens in xml")
	flag.BoolVar(&astXML, "ast", false, "output ast in xml")
	// xmlHeader := flag.Bool("xml-header", false, "output header when writing xml")
	flag.Parse()

	// lol this function is so not intuitive...
	jackFiles, parent := expandDirectories(flag.Args())
	buildDirName = filepath.Join(parent, "build")
	if _, err := os.Stat(buildDirName); !os.IsNotExist(err) {
		os.RemoveAll(buildDirName)
	}
	err := os.Mkdir(buildDirName, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create build/ folder")
	}

	fileInfo, err := os.Open("os")
	if err != nil {
		log.Fatal(err)
	}
	files, err := fileInfo.ReadDir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		srcName := filepath.Join("os", file.Name())
		src, err := os.Open(srcName)
		if err != nil {
			log.Fatal(err)
		}
		destName := filepath.Join(buildDirName, file.Name())
		dest, err := os.Create(destName)
		if err != nil {
			log.Fatal(err)
		}
		defer dest.Close()
		_, err = io.Copy(dest, src)
		if err != nil {
			log.Fatal(err)
		}
	}

	//  multi threading in GO is luvly
	var wg sync.WaitGroup // primitive used for waiting on threads
	for _, jackFile := range jackFiles {
		wg.Add(1)            // tell it how many go routines we are waiting on
		jackFile := jackFile // redeclare before passing to go routine
		go tokenizeAndParse(jackFile, &wg)
	}
	wg.Wait()

	be.Translate(buildDirName) // I don't think this outputs where we want it. Maybe this should just return a string?
}
