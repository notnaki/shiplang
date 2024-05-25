package main

import (
	"fmt"
	"os"
	"shiplang/src/lexer"
	"shiplang/src/parser"
	"shiplang/src/runtime"
)

func main() {

	bytes, err := os.ReadFile("examples/00.sp")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokens := lexer.Tokenize(string(bytes))
	ast := parser.Parse(tokens)

	// litter.Dump(ast)
	env := runtime.NewEnv(nil)
	runtime.Evaluate(ast, env)
	// litter.Dump(res)

	// litter.Dump(env)
}
