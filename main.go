package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	ttIntegerLiteral
	ttCharLiteral
	ttSemicolon
	ttPutbyte
	ttEOF
	ttCounter
)

type nodeType uint32

const (
	ntEmpty nodeType = iota
	ntRoot
	ntOpAssign
	ntOpPutbyte
	ntVariable
	ntIntLiteral
	ntCharLiteral
	ntCounter
)

// NOTE: all variables are uint32
type node struct {
	t nodeType

	rootChildNodes []node
	opArgs         []node
	varName        string
	intLitValue    uint32
	charLitValue   uint8
}

type variable struct {
	name   string
	offset uint32
}

func main() {
	log.SetFlags(0)

	if len(os.Args) != 3 {
		log.Fatal("incorrect number of arguments")
	}

	contents, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	tokens, err := tokenize(string(contents))
	if err != nil {
		log.Fatal(err)
	}

	root, err := parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	code, err := codegen(root)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(os.Args[2], []byte(code), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func tokenize(s string) ([]token, error) {
	var ts []token

	if ttCounter != 8 {
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

		// putbyte
		case strings.HasPrefix(s, "putbyte "):
			t.t = ttPutbyte
			s = s[len("putbyte "):]

		// integer
		// NOTE: no checks
		case s[0] >= '0' && s[0] <= '9':
			t.t = ttIntegerLiteral
			for len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
				t.data += (string)(s[0])
				s = s[1:]
			}

		// char
		// NOTE: no escape characters or checks
		case len(s) >= 3 && s[0] == '\'' && s[2] == '\'':
			t.t = ttCharLiteral
			t.data = string(s[1])
			s = s[3:]

		// symbol
		case s[0] >= 'a' && s[0] <= 'z':
			t.t = ttSymbol
			for len(s) > 0 && (s[0] >= 'a' && s[0] <= 'z' || s[0] == '_') {
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
		case ttIntegerLiteral:
			fmt.Printf("(integer %s) ", t.data)
		case ttPutbyte:
			fmt.Print("[putbyte] ")
		case ttSemicolon:
			fmt.Print("[;]\n")
		case ttEOF:
			fmt.Print("[EOF]\n")
		}
	}
}

func parse(ts []token) (node, error) {
	if ttCounter != 8 {
		panic("parse: not all tokens implemented")
	}
	if ntCounter != 7 {
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
		case ts[0].t == ttIntegerLiteral:
			val, err := strconv.ParseInt(ts[0].data, 10, 32)
			if err != nil {
				return node{t: ntEmpty}, err
			}

			il := node{
				t:           ntIntLiteral,
				intLitValue: uint32(val),
			}

			return il, nil

		// character literal
		case ts[0].t == ttCharLiteral:
			if len(ts[0].data) != 1 {
				panic("length of char literal != 1")
			}
			cl := node{
				t:            ntCharLiteral,
				charLitValue: uint8(ts[0].data[0]),
			}
			return cl, nil

		default:
			return node{t: ntEmpty}, fmt.Errorf("unexpected token type while trying to parse 1 token")
		}
	case len(ts) == 2:
		switch {
		// putbyte
		case ts[0].t == ttPutbyte:
			op := node{
				t: ntOpPutbyte,
			}

			arg, err := parse(ts[1:2])
			if err != nil {
				return node{t: ntEmpty}, fmt.Errorf("failed to parse arg of 'putbyte'")
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
					return node{t: ntEmpty}, err
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
		for _, child := range n.rootChildNodes {
			printAST(child)
		}

	case ntOpAssign:
		fmt.Print("[=]: ")
		for _, child := range n.opArgs {
			printAST(child)
		}
		fmt.Println("")

	case ntOpPutbyte:
		fmt.Print("[putbyte]: ")
		for _, child := range n.opArgs {
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

func checkVarExistence(vars []variable, name string) int {
	for i, v := range vars {
		if v.name == name {
			return i
		}
	}
	return -1
}

func codegen(root node) (string, error) {
	if ntCounter != 7 {
		panic("codegen: not all nodes implemented")
	}

	var code string = `.section .text
.globl _start

putbyte:
	pushq	%rbp
	movq	%rsp, %rbp

	movb	%dil, -1(%rbp)

	movq	$1, %rax
	movq	$1, %rdi
	leaq	-1(%rbp), %rsi
	movq	$1, %rdx
	syscall

	leave
	ret

_start:
	mov %rsp, %rbp

`
	var globalVars []variable
	var globalOffset uint32 = 0

	for _, node := range root.rootChildNodes {
		switch node.t {
		case ntOpAssign:
			args := node.opArgs

			// we assume the opposite (it is checked in the parser)
			if len(args) != 2 {
				panic("codegen: incorrect number of arguments for '='")
			}
			if args[0].t != ntVariable {
				panic("codegen: left argument of '=' is not a variable")
			}

			varInd := checkVarExistence(globalVars, args[0].varName)

			if varInd == -1 {
				globalOffset += 4
				v := variable{
					offset: globalOffset,
					name:   args[0].varName,
				}
				globalVars = append(globalVars, v)
				varInd = len(globalVars) - 1
			}

			switch args[1].t {
			case ntIntLiteral:
				leftVar := globalVars[varInd]
				code += fmt.Sprintf("\tmovl\t$%d, -%d(%%rbp)\n", args[1].intLitValue, leftVar.offset)

			case ntCharLiteral:
				leftVar := globalVars[varInd]
				code += fmt.Sprintf("\tmovl\t$%d, -%d(%%rbp)\n", args[1].charLitValue, leftVar.offset)

			case ntVariable:
				rightVarInd := checkVarExistence(globalVars, args[1].varName)
				if rightVarInd == -1 {
					return "", fmt.Errorf("codegen: trying to assign to a non-existing variable")
				}

				leftVar := globalVars[varInd]
				rightVar := globalVars[rightVarInd]

				code += fmt.Sprintf("\tmovl\t-%d(%%rbp), %%eax\n", rightVar.offset)
				code += fmt.Sprintf("\tmovl\t%%eax, -%d(%%rbp)\n", leftVar.offset)

			default:
				return "", fmt.Errorf("codegen: incorrect right argument of '='")
			}

		case ntOpPutbyte:
			args := node.opArgs

			// we assume the opposite (it is checked in the parser)
			if len(args) != 1 {
				panic("codegen: incorrect number of arguments for 'putbyte'")
			}

			// To make sure:
			// 1. Stack frames don't collide
			// 2. The stack is 16 byte aligned
			rspOffset := globalOffset + (16-globalOffset%16)%16

			code += fmt.Sprintf("\tsubq\t$%d, %%rsp\n", rspOffset)

			switch args[0].t {
			case ntVariable:
				varInd := checkVarExistence(globalVars, args[0].varName)
				if varInd == -1 {
					return "", fmt.Errorf("codegen: trying to print a non-existing variable")
				}

				code += fmt.Sprintf("\tmovl\t-%d(%%rbp), %%edi\n", globalVars[varInd].offset)
				code += "\tcall\tputbyte\n"

			case ntIntLiteral:
				code += fmt.Sprintf("\tmovl\t$%d, %%edi\n", args[0].intLitValue)
				code += "\tcall\tputbyte\n"

			case ntCharLiteral:
				code += fmt.Sprintf("\tmovb\t$%d, %%dil\n", args[0].charLitValue)
				code += "\tcall\tputbyte\n"

			default:
				return "", fmt.Errorf("codegen: incorrect argument of 'putbyte'")
			}

			code += fmt.Sprintf("\taddq\t$%d, %%rsp\n", rspOffset)
		}
	}

	code += `
	movq	$60, %rax
	movq	$0, %rdi
	syscall
`

	return code, nil
}
