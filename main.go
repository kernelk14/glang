// TODO(#2): Add Loops for Glang.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Version number
const VERSION = "0.0.3-beta"

// Stage: it's either `devel`, or `release`
const STAGE = "devel"

// Special Functions
var red = color.New(color.FgRed).SprintFunc()

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	redC   = red("ERROR:")
)

func err(err_str string) {
	print("%s %s", redC, err_str)
}

type GlangOps struct {
	Data  int
	Ident string
	Oper  string
}

// Here are the list of operations, after you add an operation here, go to function simulate()
// and you will find an `if` statement, you will increment the number in the `if` statement.
const (
	OP_PUSH  = iota
	OP_WRITE = iota
	OP_PLUS  = iota
	OP_MINUS = iota
	OP_MULTI = iota
	OP_DIV   = iota
	OP_POS   = iota
	OP_GOTO  = iota
	OP_COUNT = iota
)

// Stack {{{
type Stack []int

func (s *Stack) is_empty() bool {
	return len(*s) == 0
}

func (s *Stack) push(data int) {
	*s = append(*s, data)
}

func (s *Stack) pop() (int, bool) {
	if s.is_empty() {
		return 0, false
	} else {
		idx := len(*s) - 1
		elem := (*s)[idx]
		*s = (*s)[:idx]
		return elem, true
	}
}

// }}}

// The easiest way to print
var print = fmt.Printf

func evaluate(program string) {
	program_split := strings.Split(program, " ")
	var stack Stack
	var goto_stack Stack
	for i := 0; i < len(program_split); {
		code := program_split[i]
		// TODO: Comment out this Debug Messages if ever.
		if strings.HasPrefix(code, "\"") == true {
			err("Strings are not implemented yet.\n")
			break
		}
		switch code {
		case "+":
			print("%s  : A Plus instruction\n", code)
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(a + b)
			// print("Glang Debug [1st integer is in the stack]: ", aa)
			// print("\nGlang Debug [2nd integer is in the stack]: ", ab)
			// print("\n")
		case "-":
			print("%s  : A Minus instruction\n", code)
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(b - a)
			// print("Glang Debug [1st integer is in the stack]: ", aa)
			// print("\nGlang Debug [2nd integer is in the stack]: ", ab)
			// print("\n")
		case "*":
			print("%s  : A Multiply Instruction\n", code)
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(a * b)
			// print("Glang Debug [1st integer is in the stack]: ", aa)
			// print("\nGlang Debug [2nd integer is in the stack]: ", ab)
			// print("\n")
		case "/":
			print("%s  : A Divide Instruction\n", code)
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(a / b)
		case "write":
			print("%s  : A Write instruction\n", code)
			a, _ := stack.pop()
			/* if aa {
			    print("Glang Debug [Result]: %d\n", a)
			} */
			print("Glang Debug [Result]: %d\n", a)
		case "write\n": // I hardcoded this instruction with newline, because i dont have much knowledge in slicing newlines in Go.
			print("%s  : A Write instruction\n", code)
			a, _ := stack.pop()
			/* if aa {
			    print("Glang Debug [Result]: %d\n", a)
			} */
			print("Glang Debug [Result]: %d\n", a)
		case "pos":
			goto_stack.push(int(program[i+1]))
		case "goto":
			g, _ := goto_stack.pop()
			i = g
		case "end":
			// break
			os.Exit(3)
		default:
			print("%s  : Integers to be pushed\n", code)
			c_psh, _ := strconv.Atoi(code)
			stack.push(c_psh)
			// print("Glang Debug [Stack Atoi() __err__ trace]: ", err)
			// print("\n")
		}
		i += 1
		print("----------------------------------------------\n")
		// print("%s\n", code)
	}
}

func check_err(e error) {
	if e != nil {
		panic(e)
		// print("Ur Mom\n")
	}
}

// strings.ReplaceAll(string(program_file), "\n", " ")
func TrimNewLines(s string) string {
	re := regexp.MustCompile(` +\r?\n +`)
	return re.ReplaceAllString(s, " ")
}

func usage() {
	print("Glang Programming Language\n")
	print("Usage: glang <filename>\n")
}

func simulate(file string) {
	program_file, err := os.ReadFile(os.Args[1])
	check_err(err)
	// var program = string(program_file)
	program := strings.TrimSuffix(string(program_file), "\n")

	// OP_COUNT should be incremented.
	if OP_COUNT != 8 {
		print("%s: Operations not handled properly.\n", red("ERROR"))
	} else {
		evaluate(TrimNewLines(string(program)))
	}
}

func main() {
	if STAGE == "devel" {
		print("%s: Glang is in the development stage, be careful.\n\n", yellow("WARNING"))
	}
	if len(os.Args) < 2 {
		usage()

		// print("%s: Did not supply enough arguments. Maybe you forgot the filename?\n", red("ERROR"))
		err("Did not supply enough arguments. Maybe you forgot the filename?\n")
		os.Exit(1)
		// panic("Did not supply enough arguments.")
	} else {
		fileExt := filepath.Ext(os.Args[1])
		if fileExt != ".glg" {
			err("Supplying non-Glang file. Make sure it's the right file.\n")
		} else {
			simulate(string(os.Args[1]))
		}
	}
}
