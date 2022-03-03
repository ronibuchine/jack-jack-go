package targil0

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const OUTPUT_FILE_NAME string = "Tar0.asm"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Path unspecified")
	}
	dir := args[0]

	var outputFileName string
	if len(args) == 1 {
		outputFileName = OUTPUT_FILE_NAME
	} else {
		outputFileName = args[1]
	}

	ASMfile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}

	vmfiles := getVMFiles(dir)
	if len(vmfiles) == 0 {
		log.Fatal("No VM files found in given directory")
	}

	var buyAmount float64 = 0
	var sellAmount float64 = 0
	output := ""

	for _, vmfile := range vmfiles {
		vmfile, err := os.Open(dir + vmfile.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer vmfile.Close()

		scanner := bufio.NewScanner(vmfile)
		for scanner.Scan() {
			line := scanner.Text()
			words := strings.Fields(line)

			amount, err := strconv.ParseFloat(words[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			price, err := strconv.ParseFloat(words[3], 64)
			if err != nil {
				log.Fatal(err)
			}

			var lineOutput string
			var calcAmount float64
			switch words[0] {
			case "buy":
				lineOutput, calcAmount = buy(words[1], price, amount)
				buyAmount += calcAmount
			case "sell":
				lineOutput, calcAmount = sell(words[1], price, amount)
				sellAmount += calcAmount
			default:
				log.Fatal("ono")
			}

			output = output + "\n" + lineOutput
		}
	}

	total := fmt.Sprintf("TOTAL BUY: %.2f\nTOTAL SELL: %.2f", buyAmount, sellAmount)
	fmt.Print(total)
	ASMfile.WriteString(output + "\n" + total)
	ASMfile.Close()
}

func getVMFiles(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	n := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".vm" {
			files[n] = file
			n++
		}
	}
	return files[:n]
}

func buy(productName string, amount float64, price float64) (output string, totalAmount float64) {
	totalAmount = price * amount
	output = fmt.Sprintf("$$$ SELL $$$ %s\n%.2f", productName, totalAmount)
	return

}

func sell(productName string, amount float64, price float64) (output string, totalAmount float64) {
	totalAmount = price * amount
	output = fmt.Sprintf("### BUY ### %s\n%.2f", productName, totalAmount)
	return
}
