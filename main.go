package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// TOKENS

type tokenType uint32

const (
	ttEmpty tokenType = iota // used for initialization and errors
	ttBlank
	ttIdentifier
	ttSemicolon
	ttBracketOpen
	ttBracketClose
	ttEqual
	ttPlus
	ttIntLit
	ttLet
	ttU32
	ttCounter
)

type token struct {
	t    tokenType
	data string
}

// NODES

type varType uint32

const (
	vtU32 varType = iota
)

type exprType uint32

const (
	etSum exprType = iota
)

type nodeType uint32

const (
	ntEmpty nodeType = iota
	ntVarDecl
	ntExpr
	ntIntLit
	ntCounter
)

type node struct {
	t nodeType

	varDecl struct {
		t     varType
		name  string
		value *node
	}

	expr struct {
		t exprType
		arg1 *node
		arg2 *node
	}

	intLit struct {
		value int64 // NOTE: currently no support for signed literals
	}
}

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Fatal("incorrect number of arguments")
	}

	contents, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ts, err := tokenize(string(contents))
	if err != nil {
		log.Fatal(err)
	}

	printTokens(ts)

	ns, err := parse(ts)
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range ns {
		printNode(n)
	}
}

/*
const (
	ntEmpty nodeType = iota
	ntVarDecl
	ntExpr
	ntIntLit
	ntCounter
)

type node struct {
	t nodeType

	varDecl struct {
		t     varType
		name  string
		value *node
	}

	expr struct {
		t exprType
		arg1 *node
		arg2 *node
	}

	intLit struct {
		value int64 // NOTE: currently no support for signed literals
	}
}
*/

func printNode(n node) {
	switch n.t {
	case ntVarDecl:
		fmt.Println(n.varDecl.t, n.varDecl.name)
		printNode(*n.varDecl.value)
	
	case ntExpr:
		fmt.Println(n.expr.t)
		printNode(*n.expr.arg1)
		printNode(*n.expr.arg2)

	case ntIntLit:
		fmt.Println(n.intLit.value)

	case ntEmpty:
		fmt.Println("[empty]")

	default:
		panic("unknown node type")
	}
}

func parse(ts []token) ([]node, error) {
	var ns []node

	for len(ts) > 0 {
		n, left, err := chopNode(ts)
		if err != nil {
			return nil, err
		}

		if n.t != ntEmpty {
			ns = append(ns, n)
		}

		ts = left
	}

	return ns, nil
}

func chopNode(ts []token) (node, []token, error) {
	n := node{
		t: ntEmpty,
	}

	// fmt.Println(ts)

	switch {
	case len(ts) == 0:
		panic("no tokens passed")

	case len(ts) >= 3 && ts[1].t == ttPlus:
		n.t = ntExpr;
		n.expr.t = etSum;

		arg1, left, err := chopNode(ts[0:1])
		if err != nil {
			return n, nil, err
		}
		ts = ts[2:]
		arg2, left, err := chopNode(ts)
		if err != nil {
			return n, nil, err
		}

		n.expr.arg1 = &arg1
		n.expr.arg2 = &arg2
		
		return n, left, nil

	case ts[0].t == ttSemicolon:
		n.t = ntEmpty
		return n, ts[1:], nil

	case ts[0].t == ttIntLit:
		n.t = ntIntLit
		value, err := strconv.ParseInt(ts[0].data, 10, 64)
		if err != nil {
			return n, nil, err
		}
		n.intLit.value = value
		return n, ts[1:], nil

	case ts[0].t == ttEqual:
		ts = ts[1:]

		assignVal, left, err := chopNode(ts)

		if err != nil {
			return n, nil, err
		}
		if assignVal.t == ntEmpty {
			err := fmt.Errorf("assigning to an empty value")
			return n, nil, err
		}

		return assignVal, left, nil

	case ts[0].t == ttLet:
		if len(ts) < 4 {
			err := fmt.Errorf("expecting at least 4 tokens in variable declaration")
			return n, nil, err
		}

		n.t = ntVarDecl
		ts = ts[1:]

		switch ts[0].t {
		case ttU32:
			n.varDecl.t = vtU32
		default:
			err := fmt.Errorf("invalid type in variable declaration")
			return n, nil, err
		}
		ts = ts[1:]

		if ts[0].t != ttIdentifier {
			err := fmt.Errorf("variable name is not an identifier")
			return n, nil, err
		}
		n.varDecl.name = ts[0].data
		ts = ts[1:]

		varValue, left, err := chopNode(ts)
		if err != nil {
			return n, nil, err
		}
		n.varDecl.value = &varValue

		return n, left, nil

	default:
		err := fmt.Errorf("unexpected token")
		return n, nil, err
	}
}

func printTokens(ts []token) {
	if ttCounter != 11 {
		panic("not all tokens implemented")
	}

	for _, t := range ts {
		switch t.t {
		case ttEmpty:
			fmt.Printf("[EMPTY] ")
		case ttBlank:
			fmt.Printf("[BLANK] ")
		case ttIdentifier:
			fmt.Printf("%s ", t.data)
		case ttSemicolon:
			fmt.Printf("[;]\n")
		case ttBracketOpen:
			fmt.Printf("[(] ")
		case ttBracketClose:
			fmt.Printf("[)] ")
		case ttEqual:
			fmt.Printf("[=] ")
		case ttPlus:
			fmt.Printf("[+] ")
		case ttIntLit:
			fmt.Printf("<%s> ", t.data)
		case ttLet:
			fmt.Printf("[let] ")
		case ttU32:
			fmt.Printf("[u32] ")
		default:
			panic("unrecognized token")
		}
	}
}

func tokenize(s string) ([]token, error) {
	if ttCounter != 11 {
		panic("not all tokens implemented")
	}

	var ts []token

	for len(s) > 0 {
		t, left, err := chopToken(s)
		if err != nil {
			return nil, err
		}

		if t.t != ttBlank {
			ts = append(ts, t)
		}

		s = left
	}

	return ts, nil
}

func chopToken(s string) (token, string, error) {
	if ttCounter != 11 {
		panic("not all tokens implemented")
	}

	t := token{
		t:    ttEmpty,
		data: "",
	}

	switch {
	// single character tokens
	case s[0] == ';':
		t.t = ttSemicolon
		return t, s[1:], nil
	case s[0] == '(':
		t.t = ttBracketOpen
		return t, s[1:], nil
	case s[0] == ')':
		t.t = ttBracketClose
		return t, s[1:], nil
	case s[0] == '=':
		t.t = ttEqual
		return t, s[1:], nil
	case s[0] == '+':
		t.t = ttPlus
		return t, s[1:], nil

	// keywords
	case strings.HasPrefix(s, "let "):
		t.t = ttLet
		return t, s[len("let "):], nil
	case strings.HasPrefix(s, "u32 "):
		t.t = ttU32
		return t, s[len("u32 "):], nil

	// literals
	// NOTE: no checks
	case s[0] >= '0' && s[0] <= '9':
		t.t = ttIntLit
		for len(s) > 0 && (s[0] >= '0' && s[0] <= '9') {
			t.data += string(s[0])
			s = s[1:]
		}
		return t, s, nil

	// identifier
	case s[0] >= 'a' && s[0] <= 'z':
		t.t = ttIdentifier
		for len(s) > 0 && ((s[0] >= 'a' && s[0] <= 'z') || s[0] == '_') {
			t.data += string(s[0])
			s = s[1:]
		}
		return t, s, nil

	// blank characters
	case s[0] == ' ' || s[0] == '\t' || s[0] == '\n':
		t.t = ttBlank
		return t, s[1:], nil

	default:
		err := fmt.Errorf("unrecognized token starting with %c", s[0])
		t.t = ttEmpty
		return t, "", err
	}
}
