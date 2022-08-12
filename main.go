// TODO(#2): Add Loops for Glang.

package main

import (
	"bufio"
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
	yellow  = color.New(color.FgYellow).SprintFunc()
	yellowC = yellow("WARNING:")
	redC    = red("ERROR:")
)

func err(err_str string) {
	print("%s %s", redC, err_str)
}

func warn(warning string) {
	print("%s %s", yellowC, warning)
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

// I've added some tracing for stack
var trace = false

// }}}

// The easiest way to print
var print = fmt.Printf

func evaluate(program string) {
	program_split := strings.Split(program, " ")
	var stack Stack
	var goto_stack Stack
	for i := 0; i < len(program_split); {
		code := program_split[i]
		// TODO(#3): Comment out this Debug Messages if ever.

		// TODO(#4): Add support for string literals
		if strings.HasPrefix(code, "\"") == true {
			err("Strings are not implemented yet.\n")
			break
		}
		switch code {
		case "+":
			if STAGE == "devel" {
				print("%s  : A Plus instruction\n", code)
			}
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(a + b)

		case "-":
			if STAGE == "devel" {
				print("%s  : A Minus instruction\n", code)
			}
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(b - a)

		case "*":
			if STAGE == "devel" {
				print("%s  : A Multiply Instruction\n", code)
			}
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(a * b)

		case "/":
			if STAGE == "devel" {
				print("%s  : A Divide Instruction\n", code)
			}
			a, _ := stack.pop()
			b, _ := stack.pop()
			stack.push(a / b)
		case "write":
			if STAGE == "devel" {
				print("%s  : A Write instruction\n", code)
			}
			a, _ := stack.pop()
			if STAGE == "devel" {
				print("Glang Debug [Result]: %d\n", a)
			} else {
				print("%d\n", a)
			}
		case "write\n": // I hardcoded this instruction with newline, because i dont have much knowledge in slicing newlines in Go.
			if STAGE == "devel" {
				print("%s  : A Write instruction\n", code)
			}
			a, _ := stack.pop()

			if STAGE == "devel" {
				print("Glang Debug [Result]: %d\n", a)
			} else {
				print("%d\n", a)
			}
		case "pos":
			goto_stack.push(int(program[i+1]))
		case "goto":
			g, _ := goto_stack.pop()
			i = g
		case "end":
			// break
			os.Exit(1)
		case "trace":
			warn("Stack tracing started\n")
			print("Stack Contents: %d\n", stack)
		default:
			if STAGE == "devel" {
				print("%s  : Integers to be pushed\n", code)
			}
			c_psh, _ := strconv.Atoi(code)
			stack.push(c_psh)
			if trace == true {
				warn("Stack tracing started\n")
				print("Stack Contents: %d\n", stack)
			}

		}
		i += 1
		if STAGE == "devel" {
			print("----------------------------------------------\n")
		}
		// print(program_split[i])

	}
}

func check_err(e error) {
	if e != nil {
		panic(e)
		// print("Ur Mom\n")
	}
}

func TrimNewLines(s string) string {
	re := regexp.MustCompile(` +\r?\n +`)
	return re.ReplaceAllString(s, " ")
}

func usage() {
	print("Glang Programming Language\n")
	print("Usage: glang <filename>\n")
}

func sliceFileName(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

// TODO: Refine the compilation
func compile_program(file string) {
	program_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer program_file.Close()

	scanner := bufio.NewScanner(program_file)
	for scanner.Scan() {
		// Do something for the compiler
		program := scanner.Text()
		program_split := strings.Split(program, " ")
		file_with_ext := os.Args[1]
		ff := sliceFileName(file_with_ext)
		out, err := os.Create(ff + ".go")
		if err != nil {
			panic(err)
		}
		defer out.Close()
		var stack Stack
		out.WriteString("package main\n")
		out.WriteString("import \"fmt\"\n")
		out.WriteString("func main() {\n")
		for i := 0; i < len(program_split); i++ {
			if STAGE == "devel" {
				print("%s\n", program_split[i])
			}
			code := program_split[i]
			switch code {
			case "+":
				a, _ := stack.pop()
				b, _ := stack.pop()
				stack.push(a + b)
				// i += 1
				// out.WriteString("    var a = %d + %d\n", a, b)
			case "-":
				a, _ := stack.pop()
				b, _ := stack.pop()
				stack.push(b - a)
				// i += 1
			case "write":
				a, _ := stack.pop()
				out.WriteString("    fmt.Printf(\"" + "%")
				out.WriteString("d\\n\", ")
				w_str := strconv.Itoa(a)
				out.WriteString(w_str)
				out.WriteString(")\n")
				// i += 1
			case "write\n":
				a, _ := stack.pop()
				out.WriteString("    fmt.Printf(\"" + "%")
				out.WriteString("d\\n\", ")
				w_str := strconv.Itoa(a)
				out.WriteString(w_str)
				out.WriteString(")\n")
				// i += 1
			default:
				i_push, _ := strconv.Atoi(code)
				stack.push(i_push)
				// i += 1
			}
			// print("Chars: %s\n", program_split[i])
		}
		out.WriteString("}\n")
		// print("Lines: %s\n", program_split)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
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
		// compile_program(TrimNewLines(string(program)))
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
			compile_program(string(os.Args[1]))
			// simulate(string(os.Args[1]))
		}
	}
}
