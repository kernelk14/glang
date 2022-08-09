package main
import ( 
    "fmt"
    "strings"
    "strconv"
    "os"
)

type GlangOps struct {
    Data int
    Ident string
    Oper string
}
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
    var program_split = strings.Split(program, " ")
    var stack Stack
    for i := 0; i < len(program_split); i++ {
        var code = program_split[i]
        switch (code) {
            case "+":
                print("%s  : A Plus instruction\n", code)
                var a, aa = stack.pop()
                var b, ab = stack.pop()
                stack.push(a + b)
                print("Glang Debug [1st integer is in the stack]: ", aa)
                print("\nGlang Debug [2nd integer is in the stack]: ", ab)
                print("\n")
                break
            case "write":
                print("%s  : A Write instruction\n", code)
                a, aa := stack.pop()
                if aa == true {
                    print("Glang Debug [Result]: %d\n", a)
                }
                break
            case "write\n":
                print("write  : A Write instruction\n")
                a, aa := stack.pop()
                if aa == true {
                    print("Glang Debug [Result]: %d\n", a)
                }
                break
            default:
                print("%s  : Integers to be pushed\n", code)
                c_psh, err := strconv.Atoi(code)
                stack.push(c_psh)
                print("Glang Debug [Stack Atoi() __err__ trace]: ", err)
                print("\n")
                break
        }
        // print("%s\n", code)
    }
}

func check_err(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    // var program = "34 35 + write"
    
    program, err := os.ReadFile("./hello.glg")
    check_err(err)
    evaluate(string(program)) 
}
