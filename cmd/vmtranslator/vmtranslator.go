package main

import (
	"strings"
)

func normalizeSpaces(line string) (normalizedLine string) {
	temp := strings.Split(line, " ")
	var final []string
	for _, word := range temp {
		if word != "" {
			final = append(final, word)
		}
	}
	normalizedLine = strings.Join(final, " ")
	return
}

func main() {

}
