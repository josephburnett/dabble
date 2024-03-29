package main

import (
	"bufio"
	"dabble/core"
	"dabble/eval"
	"dabble/lexer"
	"dabble/parser"
	"fmt"
	"io"
	"os"
	"os/user"
)

// Based on Monkey repl.go.

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Dabble programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	Start(os.Stdin, os.Stdout)
}

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program, err := p.ParseProgram()
		if err != nil {
			io.WriteString(out, err.Error())
			io.WriteString(out, "\n")
			continue
		}

		eval.BeginTrace()
		evaluated := eval.Eval(core.Env, program)
		trace := eval.EndTrace()
		if evaluated != nil {
			io.WriteString(out, trace)
			io.WriteString(out, "\n")
			io.WriteString(out, evaluated.String())
			io.WriteString(out, "\n")
		}
	}
}
