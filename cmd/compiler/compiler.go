package main

import (
	"bufio"
	"flag"
	"fmt"
	fe "jack-jack-go/lib/syntaxAnalyzer"
	// be "jack-jack-go/lib/vmtranslator"
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
	reader := bufio.NewReader(file)
	tokens := fe.Tokenize(reader)
	if tokenXML {
		tokenXmlFile, err := os.Create("build/" + strings.TrimSuffix(file.Name(), ".jack") + "_TOKENS.xml")
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
	ast := fe.Parse(tokens)
	if astXML {
		astXmlFile, err := os.Create("build/" + strings.TrimSuffix(file.Name(), ".jack") + "_AST.xml")
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

}

var (
	tokenXML, astXML bool
)

// the main function for the compiler, the entry point to the program
func main() {
	flag.BoolVar(&tokenXML, "tokens", false, "output tokens in xml")
	flag.BoolVar(&astXML, "ast", false, "output ast in xml")
	// xmlHeader := flag.Bool("xml-header", false, "output header when writing xml")
	flag.Parse()

	args := flag.Args()
	if _, err := os.Stat("build"); !os.IsNotExist(err) {
		os.RemoveAll("build")
	}
	err := os.Mkdir("build", os.ModePerm)
	if err != nil {
		log.Fatal("Could not create build/ folder")
	}

	var jackFiles []string
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get cwd")
		}
		cwdFiles, err := os.ReadDir(cwd)
		if err != nil {
			log.Fatal("Failed to read directory")
		}
		for _, file := range cwdFiles {
			fileName := file.Name()
			if filepath.Ext(fileName) != ".jack" {
				continue
			}
			jackFiles = append(jackFiles, fileName)
		}
	} else {
		jackFiles = args
	}

    //  multi threading in GO is luvly
    var wg sync.WaitGroup // primitive used for waiting on threads
    for i := 0; i < len(jackFiles); i++ {
        wg.Add(1) // tell it how many go routines we are waiting on
        go tokenizeAndParse(jackFiles[i], &wg)
	}
    wg.Wait()
}














