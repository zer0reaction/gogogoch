package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	// "strings"
)

type token struct {
	t    tokenType
	data string
}

type tokenType uint32

const (
	ttEmpty tokenType = iota
	ttSymbol
	ttEqual
	ttInteger
	ttSemicolon
	ttDog
	ttEOF
	ttCounter
)

type nodeType uint32

const (
	ntEmpty nodeType = iota
	ntRoot
	ntOpAssign
	ntOpPutByte
	ntVariable
	ntIntLiteral
	ntCounter
)

// NOTE: all variables are uint32
type node struct {
	t nodeType

	rootChildNodes []node
	opArgs         []node
	varName        string
	intLitValue    uint32
}

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Fatal("incorrect number of arguments")
	}

	tokens, err := tokenize(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// printTokens(tokens)

	root, err := parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	printAST(root)
}

func tokenize(s string) ([]token, error) {
	var ts []token

	if ttCounter != 7 {
		panic("tokenize: not all tokens implemented")
	}

	for len(s) > 0 {
		t := token{
			t:    ttEmpty,
			data: "",
		}

		switch {
		// equal
		case s[0] == '=':
			t.t = ttEqual
			s = s[1:]

		// semicolon
		case s[0] == ';':
			t.t = ttSemicolon
			s = s[1:]

		// put byte
		case s[0] == '@':
			t.t = ttDog
			s = s[1:]

		// integer
		// NOTE: no checks
		case s[0] >= '0' && s[0] <= '9':
			t.t = ttInteger
			for len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
				t.data += (string)(s[0])
				s = s[1:]
			}

		// symbol
		case s[0] >= 'a' && s[0] <= 'z':
			t.t = ttSymbol
			for len(s) > 0 && s[0] >= 'a' && s[0] <= 'z' {
				t.data += (string)(s[0])
				s = s[1:]
			}

		default:
			s = s[1:]
		}

		if t.t != ttEmpty {
			ts = append(ts, t)
		}
	}

	ts = append(ts, token{t: ttEOF, data: ""})

	return ts, nil
}

func printTokens(ts []token) {
	if ttCounter != 7 {
		panic("printTokens: not all tokens implemented")
	}

	for _, t := range ts {
		switch t.t {
		case ttEmpty:
			fmt.Print("[EMPTY] ")
		case ttSymbol:
			fmt.Printf("%s ", t.data)
		case ttEqual:
			fmt.Print("[=] ")
		case ttInteger:
			fmt.Printf("(integer %s) ", t.data)
		case ttDog:
			fmt.Print("[@] ")
		case ttSemicolon:
			fmt.Print("[;]\n")
		case ttEOF:
			fmt.Print("[EOF]\n")
		}
	}
}

func parse(ts []token) (node, error) {
	// fmt.Printf("%+v\n", ts)

	if ttCounter != 7 {
		panic("parse: not all tokens implemented")
	}
	if ntCounter != 6 {
		panic("parse: not all nodes implemented")
	}

	switch {
	case len(ts) == 1:
		switch {
		// variable
		case ts[0].t == ttSymbol:
			v := node{
				t:       ntVariable,
				varName: ts[0].data,
			}
			return v, nil

		// int literal
		case ts[0].t == ttInteger:
			val, err := strconv.ParseInt(ts[0].data, 10, 32)
			if err != nil {
				return node{t: ntEmpty}, err
			}

			il := node{
				t:           ntIntLiteral,
				intLitValue: uint32(val),
			}

			return il, nil

		default:
			return node{t: ntEmpty}, fmt.Errorf("unexpected token type while trying to parse 1 token")
		}
	case len(ts) == 2:
		switch {
		// put byte
		case ts[0].t == ttDog:
			op := node{
				t: ntOpPutByte,
			}

			arg, err := parse(ts[1:2])
			if err != nil {
				return node{t: ntEmpty}, fmt.Errorf("failed to parse arg of '@'")
			}

			op.opArgs = append(op.opArgs, arg)
			return op, nil

		default:
			return node{t: ntEmpty}, fmt.Errorf("unexpected token types while trying to parse 2 tokens")
		}

	case len(ts) == 3:
		switch {
		// var = int or var
		case ts[1].t == ttEqual:
			op := node{
				t: ntOpAssign,
			}

			left, err := parse(ts[0:1])
			if err != nil {
				return node{t: ntEmpty}, fmt.Errorf("failed to parse left side arg of '='")
			}
			if left.t != ntVariable {
				return node{t: ntEmpty}, fmt.Errorf("left side arg of '=' is not a variable")
			}

			right, err := parse(ts[2:3])
			if err != nil {
				return node{t: ntEmpty}, fmt.Errorf("failed to parse right side arg of '='")
			}

			op.opArgs = append(op.opArgs, left)
			op.opArgs = append(op.opArgs, right)
			return op, nil

		default:
			return node{t: ntEmpty}, fmt.Errorf("unexpected token types while trying to parse 3 tokens")
		}

	case len(ts) > 3:
		root := node{
			t: ntRoot,
		}

		var accumulated []token

		for len(ts) > 0 {
			if ts[0].t == ttSemicolon {
				n, err := parse(accumulated)
				if err != nil {
					return node{t: ntEmpty}, fmt.Errorf("failed to parse root")
				}

				root.rootChildNodes = append(root.rootChildNodes, n)
				accumulated = nil
				ts = ts[1:]
			} else {
				accumulated = append(accumulated, ts[0])
				ts = ts[1:]
			}
		}

		// EOF left
		if len(accumulated) != 1 {
			return node{t: ntEmpty}, fmt.Errorf("not all tokens parsed")
		}

		return root, nil

	default:
		fmt.Println(len(ts))
		panic("unexpected number of tokens")
	}
}

func printAST(n node) {
	if ntCounter != 6 {
		panic("printAST: not all nodes implemented")
	}

	switch n.t {
	case ntRoot:
		for _, child := range(n.rootChildNodes) {
			printAST(child)
		}

	case ntOpAssign:
		fmt.Print("[=]: ")
		for _, child := range(n.opArgs) {
			printAST(child)
		}
		fmt.Println("")

	case ntOpPutByte:
		fmt.Print("[@]: ")
		for _, child := range(n.opArgs) {
			printAST(child)
		}
		fmt.Println("")

	case ntVariable:
		fmt.Printf("%s ", n.varName)

	case ntIntLiteral:
		fmt.Printf("%v ", n.intLitValue)
	
	default:
		fmt.Print("UNRECOGNIZED ")
	}
}
