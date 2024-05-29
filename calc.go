package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	addOperator      = "+"
	subtractOperator = "-"
	multiplyOperator = "*"
	divideOperator   = "/"
	powerOperator    = "^"
	leftParen        = "("
	rightParen       = ")"

	exitCommand = "exit"
)

var (
	operators     = []string{addOperator, subtractOperator, multiplyOperator, divideOperator, powerOperator}
	precedence    = map[string]int{addOperator: 1, subtractOperator: 1, multiplyOperator: 2, divideOperator: 2, powerOperator: 3}
	associativity = map[string]string{addOperator: "L", subtractOperator: "L", multiplyOperator: "L", divideOperator: "L", powerOperator: "R"}

	errInvalidOperator    = fmt.Errorf("invalid operator. Use +, -, *, /, or ^")
	errDivideByZero       = fmt.Errorf("cannot divide by zero")
	errInsufficientValues = fmt.Errorf("insufficient values for operation")
	errMismatchedParens   = fmt.Errorf("mismatched parentheses")
)

type Calculator struct {
	reader *bufio.Reader
}

func NewCalculator() *Calculator {
	return &Calculator{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (c *Calculator) Run() {
	fmt.Println("Welcome to the Go Calculator!")
	fmt.Println("Please enter your calculation with multiple operations allowed, including parentheses.")
	fmt.Println("Operators: + for addition, - for subtraction, * for multiplication, / for division, ^ for exponentiation")
	fmt.Println("Functions: sin(x), cos(x), tan(x), sqrt(x)")
	fmt.Println("Type 'exit' to quit the program.")

	for {
		fmt.Print("Enter calculation: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if strings.ToLower(input) == exitCommand {
			fmt.Println("Exiting the calculator. Goodbye!")
			break
		}

		result, err := c.evaluateExpression(input)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Please check your input and try again.")
			continue
		}

		fmt.Printf("Result: %f\n", result)
	}
}

func (c *Calculator) evaluateExpression(input string) (float64, error) {
	input = strings.ReplaceAll(input, " ", "")
	tokens, err := c.tokenize(input)
	if err != nil {
		return 0, err
	}

	postfix, err := c.infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	return c.evaluatePostfix(postfix)
}

func (c *Calculator) tokenize(input string) ([]string, error) {
	var tokens []string
	var number strings.Builder
	functionRegex := regexp.MustCompile(`^(sin|cos|tan|sqrt)\(`)

	for i := 0; i < len(input); {
		char := rune(input[i])
		if unicode.IsDigit(char) || char == '.' {
			number.WriteRune(char)
			i++
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if isOperatorOrParen(string(char)) {
				tokens = append(tokens, string(char))
				i++
			} else if match := functionRegex.FindString(input[i:]); match != "" {
				j := i + len(match)
				count := 1
				for count > 0 && j < len(input) {
					if input[j] == '(' {
						count++
					} else if input[j] == ')' {
						count--
					}
					j++
				}
				if count == 0 {
					tokens = append(tokens, input[i:j])
					i = j
				} else {
					return nil, fmt.Errorf("unmatched function parentheses")
				}
			} else {
				return nil, fmt.Errorf("invalid character: %s", string(char))
			}
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens, nil
}

func isOperatorOrParen(token string) bool {
	return token == addOperator || token == subtractOperator || token == multiplyOperator || token == divideOperator || token == powerOperator || token == leftParen || token == rightParen
}

func (c *Calculator) infixToPostfix(tokens []string) ([]string, error) {
	var postfix []string
	var stack []string

	for _, token := range tokens {
		if c.isNumber(token) {
			postfix = append(postfix, token)
		} else if c.isFunction(token) {
			stack = append(stack, token)
		} else if c.isOperator(token) {
			for len(stack) > 0 && (stack[len(stack)-1] != leftParen) && ((associativity[token] == "L" && precedence[stack[len(stack)-1]] >= precedence[token]) || (associativity[token] == "R" && precedence[stack[len(stack)-1]] > precedence[token])) {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == leftParen {
			stack = append(stack, token)
		} else if token == rightParen {
			for len(stack) > 0 && stack[len(stack)-1] != leftParen {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errMismatchedParens
			}
			stack = stack[:len(stack)-1]
			if len(stack) > 0 && c.isFunction(stack[len(stack)-1]) {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		} else {
			return nil, fmt.Errorf("invalid token: %s", token)
		}
	}
	for len(stack) > 0 {
		if stack[len(stack)-1] == leftParen {
			return nil, errMismatchedParens
		}
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return postfix, nil
}

func (c *Calculator) evaluatePostfix(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		if c.isNumber(token) {
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token)
			}
			stack = append(stack, value)
		} else if c.isFunction(token) {
			result, err := c.evaluateFunction(token)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
		} else if c.isOperator(token) {
			if len(stack) < 2 {
				return 0, errInsufficientValues
			}
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			var result float64
			var err error
			switch token {
			case addOperator:
				result = c.add(a, b)
			case subtractOperator:
				result = c.subtract(a, b)
			case multiplyOperator:
				result = c.multiply(a, b)
			case divideOperator:
				result, err = c.divide(a, b)
				if err != nil {
					return 0, err
				}
			case powerOperator:
				result = c.power(a, b)
			default:
				return 0, errInvalidOperator
			}

			stack = append(stack, result)
		} else {
			return 0, fmt.Errorf("invalid token: %s", token)
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("error evaluating expression")
	}

	return stack[0], nil
}

func (c *Calculator) isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func (c *Calculator) isOperator(token string) bool {
	return isOperatorOrParen(token) && token != leftParen && token != rightParen
}

func (c *Calculator) isFunction(token string) bool {
	return strings.HasPrefix(token, "sin(") || strings.HasPrefix(token, "cos(") || strings.HasPrefix(token, "tan(") || strings.HasPrefix(token, "sqrt(")
}

func (c *Calculator) evaluateFunction(token string) (float64, error) {
	parts := strings.SplitN(token, "(", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid function format: %s", token)
	}
	funcName := parts[0]
	argStr := strings.TrimSuffix(parts[1], ")")

	// Evaluate the argument expression
	arg, err := c.evaluateExpression(argStr)
	if err != nil {
		return 0, fmt.Errorf("invalid function argument: %s", argStr)
	}

	switch funcName {
	case "sin":
		return math.Sin(arg), nil
	case "cos":
		return math.Cos(arg), nil
	case "tan":
		return math.Tan(arg), nil
	case "sqrt":
		return math.Sqrt(arg), nil
	default:
		return 0, fmt.Errorf("unsupported function: %s", funcName)
	}
}

func (c *Calculator) add(a, b float64) float64 {
	return a + b
}

func (c *Calculator) subtract(a, b float64) float64 {
	return a - b
}

func (c *Calculator) multiply(a, b float64) float64 {
	return a * b
}

func (c *Calculator) divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errDivideByZero
	}
	return a / b, nil
}

func (c *Calculator) power(a, b float64) float64 {
	return math.Pow(a, b)
}

func main() {
	calculator := NewCalculator()
	calculator.Run()
}
