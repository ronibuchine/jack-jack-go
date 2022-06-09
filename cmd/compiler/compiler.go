package main

import (
	"bufio"
	"flag"
	"fmt"
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
		tokenXmlFile, err := os.Create(filepath.Join("build", name+"_TOKENS.xml"))
		if err != nil {
			log.Fatal(err)
		}
		defer tokenXmlFile.Close()
		w := bufio.NewWriter(tokenXmlFile)
		err = fe.TokenToXML(tokens, w)
		if err != nil {
			log.Print(fmt.Sprint("Could not write token XML for ", file.Name()))
		}
		w.Flush()
	}
	ts := fe.TS{Tokens: tokens, File: jackFile}
	ast := fe.Parse(&ts)
	if astXML {
		astXmlFile, err := os.Create(filepath.Join("build", name+"_AST.xml"))
		if err != nil {
			log.Fatal(err)
		}
		defer astXmlFile.Close()
		w := bufio.NewWriter(astXmlFile)
		err = fe.NodeToXML(ast, w)
		if err != nil {
			log.Print(fmt.Sprint("Could not write ast XML for ", file.Name()))
		}
		w.Flush()
	}
	vmFile, err := os.Create(filepath.Join("build", name+".vm"))
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

var (
	tokenXML, astXML bool
)

// takes command line args and returns all files inside
func expandDirectories(jackFiles []string) []string {
	files := make([]string, 0)
	for _, arg := range jackFiles {
		fileInfo, err := os.Stat(arg)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
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
			files = append(files, arg)
		}
	}
	return files
}

// the main function for the compiler, the entry point to the program
func main() {
	flag.BoolVar(&tokenXML, "tokens", false, "output tokens in xml")
	flag.BoolVar(&astXML, "ast", false, "output ast in xml")
	// xmlHeader := flag.Bool("xml-header", false, "output header when writing xml")
	flag.Parse()

	// lol this function is so not intuitive...
	jackFiles := expandDirectories(flag.Args())
	if _, err := os.Stat("build"); !os.IsNotExist(err) {
		os.RemoveAll("build")
	}
	err := os.Mkdir("build", os.ModePerm)
	if err != nil {
		log.Fatal("Could not create build/ folder")
	}

	//  multi threading in GO is luvly
	var wg sync.WaitGroup // primitive used for waiting on threads
	for _, jackFile := range jackFiles {
		wg.Add(1)            // tell it how many go routines we are waiting on
		jackFile := jackFile // redeclare before passing to go routine
		go tokenizeAndParse(jackFile, &wg)
	}
	wg.Wait()
    
    be.Translate("build") // I don't think this outputs where we want it. Maybe this should just return a string?
}
