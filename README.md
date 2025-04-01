![banner](https://raw.githubusercontent.com/2105789/Lango/main/banner.jpg)
# Lango Language Interpreter

A simple interpreter for the Lango programming language, implemented in Go.

This project implements a tree-walking interpreter featuring:

*   Dynamic Typing
*   Basic Arithmetic (+, -, *, /, %)
*   Logical Operators (and, or, !)
*   Comparison Operators (>, >=, <, <=, ==, !=)
*   Variables (declaration with `var`, assignment)
*   Control Flow (if/else, while, for)
*   Blocks and Lexical Scoping (`{ ... }`)
*   Functions (first-class, closures, `return`)
*   Classes (declaration, instantiation, methods, properties, `this`)
*   Arrays/Lists (literals `[]`, indexing `[]`, assignment `[]=`) 

## Features

### Variables and Types

```lango
var message = "Hello";
var count = 10;
var flag = true;
var nothing = nil;

print message + " World!"; // Output: Hello World!
print count * 2;          // Output: 20
print !flag;              // Output: false
```

### Control Flow

```lango
var x = 5;
if (x > 3) {
  print "x is greater than 3";
} else {
  print "x is not greater than 3";
}

var i = 0;
while (i < 3) {
  print i; // 0, 1, 2
  i = i + 1;
}

for (var j = 5; j > 2; j = j - 1) {
  print j; // 5, 4, 3
}
```

### Functions

Functions are first-class citizens and support closures.

```lango
fun sayHi(name) {
  print "Hi, " + name + "!";
}

sayHi("Lango"); // Output: Hi, Lango!

fun makeAdder(n) {
  fun adder(x) {
    return x + n;
  }
  return adder;
}

var addTwo = makeAdder(2);
print addTwo(5); // Output: 7
```

### Classes

Basic class support with methods, properties, and `this`.

```lango
class Greeter {
  setName(name) {
    this.name = name;
  }
  
  greet() {
    print "Hello, my name is " + this.name + ".";
  }
}

var g = Greeter();
g.setName("LangoClass");
g.greet(); // Output: Hello, my name is LangoClass.
print g.name; // Output: LangoClass
```

### Arrays / Lists

Simple dynamic arrays (lists).

```lango
var list = [1, "two", true];
print list;      // Output: [1, two, true]
print list[0];   // Output: 1

list[1] = 2;     // Modify element
print list;      // Output: [1, 2, true]

var nested = [[1, 2], [3, 4]];
print nested[0][1]; // Output: 2
```

## Building and Running

1.  **Navigate** to the project directory.
2.  **Initialize Go Module** (if not already done):
    ```bash
    go mod init Lango # Or your preferred module path
    ```
3.  **Build** the interpreter:
    ```bash
    go build
    ```
4.  **Run**:
    *   **REPL:**
        ```bash
        ./Lango        # Linux/macOS
        .\Lango.exe   # Windows
        ```
    *   **Script File:**
        ```bash
        ./Lango your_script.lango        # Linux/macOS
        .\Lango.exe your_script.lango   # Windows
        ```

## Development

Built following principles often found in interpreters like the one described in "Crafting Interpreters". Code includes:

*   `scanner.go`: Lexical analysis (tokenization).
*   `parser.go`: Syntax analysis (builds Abstract Syntax Tree - AST).
*   `expr.go`, `stmt.go`: AST node definitions.
*   `interpreter.go`: Tree-walking interpreter logic.
*   `environment.go`: Handles variable scopes.
*   `token.go`, `tokentype.go`: Token definitions.
*   `main.go`: Main entry point, REPL, file running.
*   `astprinter.go`: Utility to print the AST (optional).
