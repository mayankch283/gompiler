package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func repl() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Minimal Compiler REPL!")
	fmt.Println("Type 'exit' to quit.")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		// Compile the input using the compiler function
		output, err := compiler(input)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Output:", output)
		}
	}

	fmt.Println("Goodbye!")
}

func main() {
	repl()
}
