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

type variable struct {
	name   string
	offset int32
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

	code, err := codegen(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code)
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

	case ntOpPutByte:
		fmt.Print("[@]: ")
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

/*
	ntOpAssign
	ntOpPutByte
	ntVariable
	ntIntLiteral
*/

func checkVarExistence(vars []variable, name string) int {
	for i, v := range vars {
		if v.name == name {
			return i
		}
	}
	return -1
}

func codegen(root node) (string, error) {
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
	var globalOffset int32 = 0

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
				globalOffset -= 4
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
				code += fmt.Sprintf("\tmovl\t$%d, %d(%%rbp)\n", args[1].intLitValue, leftVar.offset)

			case ntVariable:
				rightVarInd := checkVarExistence(globalVars, args[1].varName)
				if rightVarInd == -1 {
					return "", fmt.Errorf("codegen: trying to assign to a non-existing variable")
				}

				leftVar := globalVars[varInd]
				rightVar := globalVars[rightVarInd]

				code += fmt.Sprintf("\tmovl\t%d(%%rbp), %%eax\n", rightVar.offset)
				code += fmt.Sprintf("\tmovl\t%%eax, %d(%%rbp)\n", leftVar.offset)

			default:
				return "", fmt.Errorf("codegen: incorrect right argument of '='")
			}

		case ntOpPutByte:
			args := node.opArgs

			// we assume the opposite (it is checked in the parser)
			if len(args) != 1 {
				panic("codegen: incorrect number of arguments for '@'")
			}

			var allignOffset int32 = 0
			if (-globalOffset)%16 != 0 {
				allignOffset = -(16 - (-globalOffset)%16)
				code += fmt.Sprintf("\taddq\t$%d, %%rsp\n", allignOffset)
			}

			switch args[0].t {
			case ntVariable:
				varInd := checkVarExistence(globalVars, args[0].varName)
				if varInd == -1 {
					return "", fmt.Errorf("codegen: trying to print a non-existing variable")
				}

				code += fmt.Sprintf("\tmovl\t%d(%%rbp), %%edi\n", globalVars[varInd].offset)
				code += "\tcall\tputbyte\n"

			case ntIntLiteral:
				code += fmt.Sprintf("\tmovl\t$%d, %%edi\n", args[0].intLitValue)
				code += "\tcall\tputbyte\n"

			default:
				return "", fmt.Errorf("codegen: incorrect argument of '@'")
			}

			if allignOffset != 0 {
				code += fmt.Sprintf("\tsubq\t$%d, %%rsp\n", allignOffset)
			}
		}
	}

	code += `
	movq	$60, %rax
	movq	$0, %rdi
	syscall`

	return code, nil
}
