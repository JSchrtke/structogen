package main

import (
	"fmt"
	"os"
)

func main() {
	templateBytes, err := os.ReadFile("./template.str")
	if err != nil {
		panic(err)
	}

	templateString := string(templateBytes)
	tokens := makeTokens(templateString)

	parsed, err := parseStructogram(tokens)
	if err != nil {
		panic(err)
	}
	parsedJson, err := parsed.ToJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%s", parsedJson))
}
