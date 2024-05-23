package main

import (
	"flag"
	"fmt"
	"os"

	// interpreter "shipgo/src/interpreter_old"
	"shipgo/src/interpreter"
	"shipgo/src/lexer"
	"shipgo/src/parser"

	"github.com/sanity-io/litter"
)

func main() {

	dumpTokens := flag.Bool("tokens", false, "Dump generated tokens")
	dumpAST := flag.Bool("ast", false, "Dump generated AST")
	dumpEnv := flag.Bool("env", false, "Dump generated environment")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: go run main.go [--tokens] [--ast] [--env] <filename>")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokens := lexer.Tokenize(string(bytes))
	ast := parser.Parse(tokens)
	env := interpreter.CreateGlobalEnv()
	interpreter.Evaluate(ast, env)

	if *dumpTokens {
		fmt.Println("--------------- Generated Tokens ---------------")
		litter.Dump(tokens)
	}

	if *dumpAST {
		fmt.Println("--------------- Generated AST ---------------")
		litter.Dump(ast)
	}

	if *dumpEnv {
		fmt.Println("--------------- Generated Env ---------------")
		litter.Dump(env)
	}
}
