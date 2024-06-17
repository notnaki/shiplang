package main

import (
	"encoding/gob"
	"flag"
	"fmt"

	"log"
	"os"
	"shiplang/src/ast"
	"shiplang/src/lexer"
	"shiplang/src/parser"
	"shiplang/src/runtime"
	"strings"

	"github.com/sanity-io/litter"
)

type StmtList []ast.Stmt

type encodedAST struct {
	ast.Stmt
	Type string   `json:"type"` // Add a field to store the actual AST type
	Body StmtList `json:"body"` // Use StmtList for the body
}

func saveAST(filename string, AST ast.Stmt) error {
	// Open the file for writing in binary mode
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Create a new encoder and encode the AST with a custom type for Body
	encoder := gob.NewEncoder(file)

	var encAst encodedAST = encodedAST{Type: "BlockStmt", Body: AST.(ast.BlockStmt).Body}
	encAst.Stmt = AST
	err = encoder.Encode(encAst)

	if err != nil {
		return fmt.Errorf("error encoding AST: %v", err)
	}

	return nil
}

func loadAST(filename string) (ast.Stmt, error) {
	// Open the file for reading in binary mode
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a new decoder and decode into a temporary struct
	decoder := gob.NewDecoder(file)
	var encodedAST encodedAST
	err = decoder.Decode(&encodedAST)
	if err != nil {
		return nil, fmt.Errorf("error decoding AST: %v", err)
	}

	return encodedAST.Stmt, nil
}

func main() {

	dumpTokens := flag.Bool("tokens", false, "Dump generated tokens")
	dumpAST := flag.Bool("ast", false, "Dump generated AST")
	dumpEnv := flag.Bool("env", false, "Dump generated environment")
	makeRunnable := flag.Bool("runnable", false, "Creates a runnable AST")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: go run main.go [--tokens] [--ast] [--env] [--runnable] <filename>")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	if strings.HasSuffix(filename, ".spr") {
		// Load and run the serialized AST
		loadedAst, err := loadAST(filename)
		if err != nil {
			fmt.Printf("Error loading file: %v\n", err)
			os.Exit(1)
		}
		env := runtime.NewEnv(nil)
		runtime.Evaluate(loadedAst, env)
	} else {
		// Process the source file and optionally create a runnable AST
		bytes, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		tokens := lexer.Tokenize(string(bytes))
		parsedAst := parser.Parse(tokens)
		env := runtime.NewEnv(nil)

		if *makeRunnable {
			sprFilename := strings.TrimSuffix(filename, ".sp") + ".spr"
			err := saveAST(sprFilename, parsedAst)
			if err != nil {
				log.Fatal("Error creating .spr file:", err)
			}
			fmt.Println("Created runnable AST:", sprFilename)
		} else {
			runtime.Evaluate(parsedAst, env)
		}

		if *dumpTokens {
			fmt.Println("--------------- Generated Tokens ---------------")
			litter.Dump(tokens)
		}

		if *dumpAST {
			fmt.Println("--------------- Generated AST ---------------")
			litter.Dump(parsedAst)
		}

		if *dumpEnv {
			fmt.Println("--------------- Generated Env ---------------")
			litter.Dump(env)
		}
	}
}
