# Lango Programming Language

Lango is a simple interpreted programming language implemented in Go. It supports basic arithmetic operations, variable declarations, control flow statements, and more.

Inspired by [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom.

## Table of Contents

1. [Features](#features)
2. [Project Structure](#project-structure)
3. [Usage](#usage)
4. [Language Syntax](#language-syntax)
5. [Examples](#examples)

## Features

- Basic arithmetic operations (+, -, *, /, %)
- Variable declarations and assignments
- Control flow statements (if, while, for)
- Print statements
- Support for numbers, strings, and boolean values

## Project Structure

The project consists of several Go files, each responsible for a specific part of the language implementation:

- `main.go`: Entry point of the interpreter
- `scanner.go`: Tokenizes the input source code
- `parser.go`: Parses the tokens into an Abstract Syntax Tree (AST)
- `interpreter.go`: Executes the parsed AST
- `expr.go`: Defines expression types
- `stmt.go`: Defines statement types
- `token.go`: Defines token types and structure
- `environment.go`: Manages variable scoping and storage
- `astprinter.go`: Utility for printing the AST (useful for debugging)

## Usage

To run a Lango script:

```
go run . <script_name>.lango
```

To start an interactive REPL (Read-Eval-Print Loop):

```
go run .
```

## Language Syntax

### Variables

```lango
var x = 10;
var y = "Hello, World!";
var z = true;
```

### Arithmetic Operations

```lango
var a = 5 + 3 * 2;
var b = 10 / 2 - 1;
var c = 7 % 3;
```

### Control Flow

#### If Statement

```lango
if (x > 5) {
    print "x is greater than 5";
} else {
    print "x is not greater than 5";
}
```

#### While Loop

```lango
var i = 0;
while (i < 5) {
    print i;
    i = i + 1;
}
```

#### For Loop

```lango
for (var j = 0; j < 5; j = j + 1) {
    print j;
}
```

### Print Statement

```lango
print "Hello, World!";
print 42;
print true;
```

## Examples

Here are some examples demonstrating the features of Lango:

### Example 1: Basic Arithmetic and Variables

```lango
var x = 10;
var y = 5;
print x + y;
print x * y;
print x / y;
print x % y;
```

Output:
```
15
50
2
0
```

### Example 2: Control Flow

```lango
var i = 0;
while (i < 5) {
    if (i % 2 == 0) {
        print "Even";
    } else {
        print "Odd";
    }
    i = i + 1;
}
```

Output:
```
Even
Odd
Even
Odd
Even
```

### Example 3: For Loop and String Concatenation

```lango
for (var i = 1; i <= 5; i = i + 1) {
    print "Count: " + i;
}
```

Output:
```
Count: 1
Count: 2
Count: 3
Count: 4
Count: 5
```

## Implementation Details

### Scanner (`scanner.go`)

The scanner tokenizes the input source code. It reads the source character by character and produces a list of tokens. Key functions include:

- `ScanTokens()`: Main function that scans the entire source and returns a list of tokens
- `scanToken()`: Scans a single token
- `isDigit()`, `isAlpha()`, `isAlphaNumeric()`: Helper functions for character classification

### Parser (`parser.go`)

The parser takes the list of tokens from the scanner and builds an Abstract Syntax Tree (AST). Key functions include:

- `Parse()`: Main parsing function that returns a list of statements
- `declaration()`, `statement()`, `expression()`: Parse different language constructs
- `match()`, `consume()`: Helper functions for token matching and consumption

### Interpreter (`interpreter.go`)

The interpreter executes the AST produced by the parser. It implements the Visitor pattern to traverse and execute each node of the AST. Key functions include:

- `Interpret()`: Main interpretation function
- `VisitXxxExpr()` and `VisitXxxStmt()`: Visitor methods for different expression and statement types
- `execute()` and `evaluate()`: Helper functions for statement execution and expression evaluation

### Environment (`environment.go`)

The Environment manages variable scoping and storage. It supports nested scopes for block-level variable declarations. Key functions include:

- `Define()`: Defines a new variable in the current scope
- `Get()`: Retrieves a variable's value
- `Assign()`: Assigns a new value to an existing variable

## Conclusion

Lango is a simple yet functional programming language that demonstrates the basics of language design and implementation. It provides a foundation for further experimentation and expansion, such as adding functions, classes, or more advanced control flow constructs.
