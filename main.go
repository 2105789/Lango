package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var hadRuntimeError bool
var interpreter = NewInterpreter()

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: Lango [script.lango]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		path := os.Args[1]
		if filepath.Ext(path) != ".lango" {
			fmt.Println("Error: Script must have '.lango' extension")
			os.Exit(1)
		}
		runFile(path)
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	run(string(bytes))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "q" || line == "quit" {
			fmt.Println("Exiting...")
			break
		}
		run(line)
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	tokenPtrs := make([]*Token, len(tokens))
	for i := range tokens {
		tokenPtrs[i] = &tokens[i]
	}

	parser := NewParser(tokenPtrs)
	expression, err := parser.Parse()
	if err != nil {
		fmt.Println("Error during parsing:", err)
		return
	}

	interpreter.Interpret(expression)
}

func printEnvironment(env *Environment) {
	if env == nil {
		fmt.Println("Environment is nil")
		return
	}

	for key, value := range env.values {
		fmt.Printf("  %s: %v\n", key, value)
	}

	if env.enclosing != nil {
		fmt.Println("  Enclosing:")
		printEnvironment(env.enclosing)
	}
}

func runtimeError(err error) {
	fmt.Println(err.Error())
	hadRuntimeError = true
}

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
}

func parseError(token *Token, message string) {
	if token.Type == EOF {
		fmt.Printf("[line %d] Error at end: %s\n", token.Line, message)
	} else {
		fmt.Printf("[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, message)
	}
}
