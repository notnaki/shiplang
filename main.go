package main

import (
	"os"

	"shipgo/src/interpreter"
	"shipgo/src/lexer"
	"shipgo/src/parser"
)

func main() {

	bytes, _ := os.ReadFile("./examples/01.sp")
	tokens := lexer.Tokenize(string(bytes))

	ast := parser.Parse(tokens)
	env := interpreter.CreateGlobalEnv()
	// println("--------------- Generated Tokens ---------------")
	// litter.Dump(tokens)
	// println("--------------- Generated Ast ---------------")
	// litter.Dump(ast)
	interpreter.Evaluate(ast, env)

	// println("--------------- Generated Env ---------------")
	// litter.Dump(env)

	// println("--------------- Generated Result ---------------")

	// litter.Dump(res)
}
