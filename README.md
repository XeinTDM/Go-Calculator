# Go Calculator

Welcome to the Go Calculator! This is a command-line calculator application written in Go, capable of performing a variety of mathematical operations, including basic arithmetic, exponentiation, and trigonometric functions.

### Features
- Supports basic arithmetic operations: addition (`+`), subtraction (`-`), multiplication (`*`), division (`/`), and exponentiation (`^`).
- Handles nested and multiple operations with parentheses.
- Includes trigonometric functions: `sin(x)`, `cos(x)`, `tan(x)`.
- Provides square root function: `sqrt(x)`.
- User-friendly interface with prompt for user input.
- Exits cleanly with the `exit` command.

## Getting Started
### Prerequisites
- Go (1.16 or higher)

### Installation
1. **Clone the repository:**
```bash
git clone https://github.com/XeinTDM/Go-Calculator.git
cd go-calculator
```

2. **Run the executable:**
```bash
go build -o calculator.exe calc.go
```

3. **(Optional) Compress the executable with UPX for smaller size:**
```bash
upx --best calculator.exe
```

### Usage
1. **Run the calculator:**
```bash
./calculator
```

2. **Enter your calculations when prompted.**
- Example calculations:
```bash
Enter calculation: 3 + 5 * (2 - 4)
Result: -7.000000
```
```bash
Enter calculation: sin(3.14 / 2)
Result: 1.000000
```

3. **Exit the calculator:**
Type `exit` and press Enter to quit the program.

## How It Works
The calculator reads input from the user, processes the input to convert it into a format that can be evaluated, and then computes the result. It uses the Shunting Yard algorithm to handle operator precedence and associativity, converting infix expressions to postfix notation for easier evaluation.

### Key Components

- **Tokenizer:** Splits the input string into meaningful tokens (numbers, operators, functions).
- **Shunting Yard Algorithm:** Converts infix expressions to postfix notation.
- **Postfix Evaluator:** Computes the result from the postfix expression.
- **Error Handling:** Manages various errors such as invalid operators, mismatched parentheses, and division by zero.
