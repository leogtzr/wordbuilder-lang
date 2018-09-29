package main

import (
	// "fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"wordbuilder/evaluator"
	"wordbuilder/lexer"
	"wordbuilder/object"
	"wordbuilder/parser"
)

func main() {

	args := os.Args
	if len(args) <= 1 {
		log.Fatal("Wrong number of arguments ... ")
	}

	fileArg := os.Args[1]
	programFile, err := os.Open(fileArg)

	if err != nil {
		log.Fatal("error opening file ... ", err)
	}

	defer programFile.Close()

	out := os.Stdout

	programContent, _ := ioutil.ReadAll(programFile)

	// fmt.Println(string(programContent))

	l := lexer.New(string(programContent))
	p := parser.New(l)
	env := object.NewEnvironment()

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParseErrors(out, p.Errors())
		os.Exit(1)
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
