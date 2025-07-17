package main

import (
	"fmt"
	"log"
	"os"
	// "strconv"
	"strings"
)

// TOKENS

type tokenType uint32

const (
	ttEmpty tokenType = iota // used for initialization and errors
	ttBlank
	ttEOF
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

type nodeType uint32

const (
	ntEmpty nodeType = iota
	ntCounter
)

type node struct {
	t nodeType
}

// VARIABLES

type variableType uint32

const (
	vtU32 variableType = iota
)

type variable struct {
	t      variableType
	name   string
	offset uint32
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
}

func printTokens(ts []token) {
	if ttCounter != 12 {
		panic("not all tokens implemented")
	}

	for _, t := range ts {
		switch t.t {
		case ttEmpty:
			fmt.Printf("[EMPTY] ")
		case ttBlank:
			fmt.Printf("[BLANK] ")
		case ttEOF:
			fmt.Printf("[EOF]\n")
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
	if ttCounter != 12 {
		panic("not all tokens implemented")
	}

	var ts []token
	eof := false

	for !eof {
		t, left, err := chopToken(s)
		if err != nil {
			return nil, err
		}

		if t.t != ttBlank {
			ts = append(ts, t)
		}
		if t.t == ttEOF {
			eof = true
		}

		s = left
	}

	return ts, nil
}

func chopToken(s string) (token, string, error) {
	if ttCounter != 12 {
		panic("not all tokens implemented")
	}

	t := token{
		t:    ttEmpty,
		data: "",
	}

	switch {
	case len(s) == 0:
		t.t = ttEOF
		return t, "", nil

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
