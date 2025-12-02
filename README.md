![banner](https://raw.githubusercontent.com/2105789/Lango/main/banner.jpg)

# Lango Programming Language

A modern, dynamically-typed programming language implemented in Go that combines simplicity with powerful features for building APIs, web clients, and data processing applications.

## ⚡ Features

### Core Language
- **Dynamic Typing** - Flexible type system
- **First-Class Functions** - Functions as values, closures, higher-order functions
- **Object-Oriented** - Classes with methods, properties, and inheritance via `this`
- **Object Literals** - Create objects with `{key: value}` syntax
- **Control Flow** - if/else, while, for loops
- **Arrays/Lists** - Dynamic arrays with indexing and nested structures
- **Escape Sequences** - Support for `\n`, `\t`, `\r`, `\\`, `\"`, `\'` in strings
- **Lexical Scoping** - Block-level scope with proper closure support

### 🌐 Networking & HTTP
- **HTTP Client** - Full support for GET, POST, PUT, DELETE requests
- **Structured Responses** - Access status codes, headers, and body separately
- **Custom Headers** - Send custom headers with requests
- **JSON Integration** - Seamless JSON parsing and serialization

### 📦 Built-in Libraries

#### HTTP Functions
```lango
var response = httpGet("https://api.example.com/users/1");
var data = jsonParse(response.body);
print data.name;
```

#### JSON Operations
```lango
var json = jsonStringify(object);
var prettyJson = jsonStringifyPretty(object);
var parsed = jsonParse(jsonString);
```

#### File I/O
```lango
fileWrite("data.txt", "Hello, Lango!");
var content = fileRead("data.txt");
fileAppend("data.txt", " More text.");
var exists = fileExists("data.txt");
fileDelete("data.txt");
fileCopy("source.txt", "dest.txt");
var files = fileList(".");
```

#### String Functions
```lango
strLen(str)              // Length
strUpper(str)            // Uppercase
strLower(str)            // Lowercase
strSplit(str, delim)     // Split to array
strJoin(array, delim)    // Join from array
strContains(str, substr) // Contains check
strReplace(str, old, new)// Replace all
strTrim(str)             // Trim whitespace
```

#### Array Functions
```lango
arrayLen(arr)            // Length
arrayPush(arr, item)     // Add item
arrayPop(arr)            // Get last item
arraySlice(arr, start, end) // Slice
arrayConcat(arr1, arr2)  // Concatenate
```

#### Math Functions
```lango
mathFloor(num)           // Round down
mathCeil(num)            // Round up
mathRound(num)           // Round nearest
mathAbs(num)             // Absolute value
mathMax(a, b)            // Maximum
mathMin(a, b)            // Minimum
mathPow(base, exp)       // Power
mathSqrt(num)            // Square root
mathRandom()             // Random 0-1
```

#### Type Conversion
```lango
toNumber(value)          // Convert to number
toString(value)          // Convert to string
toBoolean(value)         // Convert to boolean
```

#### Utilities
```lango
time()                   // Unix timestamp
sleep(milliseconds)      // Sleep
```

## 🚀 Quick Start

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/Lango.git
cd Lango
```

2. **Build the interpreter**
```bash
go build
```

### Running Lango

**REPL Mode:**
```bash
./Lango        # Linux/macOS
.\Lango.exe    # Windows
```

**Script File:**
```bash
./Lango script.lango        # Linux/macOS
.\Lango.exe script.lango    # Windows
```

## 📝 Examples

### Example 1: Basic Syntax

```lango
// Variables and types
var name = "Lango";
var version = 1.0;
var isAwesome = true;

// Functions
fun greet(name) {
  return "Hello, " + name + "!";
}

print greet("World");

// Classes
class Counter {
  init(start) {
    this.value = start;
  }
  
  increment() {
    this.value = this.value + 1;
  }
}

var counter = Counter();
counter.init(0);
counter.increment();
print counter.value; // Output: 1
```

### Example 2: Building an API Client

```lango
// Fetch user data from API
var response = httpGet("https://jsonplaceholder.typicode.com/users/1");

if (response.status == 200) {
  var user = jsonParse(response.body);
  
  print "Name: " + user.name;
  print "Email: " + user.email;
  print "City: " + user.address.city;
  
  // Save to file
  fileWrite("user.json", jsonStringifyPretty(user));
  print "Saved to user.json";
}
```

### Example 3: Data Processing

```lango
// Process CSV data
var csv = "John,25,NYC,Jane,30,LA,Bob,35,SF";
var parts = strSplit(csv, ",");

var users = [];
for (var i = 0; i < arrayLen(parts); i = i + 3) {
  var name = parts[i];
  var age = toNumber(parts[i + 1]);
  var city = parts[i + 2];
  
  print name + " is " + toString(age) + " years old from " + city;
}
```

### Example 4: File Operations

```lango
// Read, process, and save data
var data = fileRead("input.txt");
var lines = strSplit(data, "\n");

var output = "Processed Data:\n";
for (var i = 0; i < arrayLen(lines); i = i + 1) {
  var line = strTrim(lines[i]);
  if (strLen(line) > 0) {
    output = output + strUpper(line) + "\n";
  }
}

fileWrite("output.txt", output);
```

## 📚 Language Features

### Variables
```lango
var x = 10;
var name = "Lango";
var flag = true;
var nothing = nil;
```

### Control Flow
```lango
if (x > 5) {
  print "Greater";
} else {
  print "Smaller";
}

while (x > 0) {
  print x;
  x = x - 1;
}

for (var i = 0; i < 10; i = i + 1) {
  print i;
}
```

### Functions & Closures
```lango
fun makeCounter() {
  var count = 0;
  fun increment() {
    count = count + 1;
    return count;
  }
  return increment;
}

var counter = makeCounter();
print counter(); // 1
print counter(); // 2
```

### Classes & Objects
```lango
class Person {
  setName(name) {
    this.name = name;
  }
  
  greet() {
    print "Hi, I'm " + this.name;
  }
}

var person = Person();
person.setName("Alice");
person.greet();
```

### Arrays
```lango
var numbers = [1, 2, 3, 4, 5];
print numbers[0];        // 1
numbers[1] = 10;         // Modify
var nested = [[1, 2], [3, 4]];
print nested[0][1];      // 2
```

### Object Literals
```lango
// Create objects with literal syntax
var person = {name: "Alice", age: 30, city: "NYC"};
print person.name;       // Alice
person.age = 31;         // Modify property

// Nested objects
var config = {server: {host: "localhost", port: 8080}};
print config.server.host; // localhost

// Mixed values
var data = {str: "hello", num: 42, arr: [1, 2, 3]};
```

### Escape Sequences
```lango
print "Line 1\nLine 2";     // Newline
print "Tab\there";          // Tab
print "Quote: \"hi\"";      // Escaped quotes  
print "Back\\slash";        // Backslash
var multiline = "First\n\tIndented\nLast";
```

## 🎯 Use Cases

- **REST API Clients** - Build clients to consume REST APIs
- **Web Scraping** - Fetch and process web data
- **Data Processing** - Transform and analyze data
- **File Automation** - Automate file operations
- **Configuration Management** - Generate and manage config files
- **System Scripting** - Quick automation scripts

## 📖 Documentation

- [FEATURES.md](FEATURES.md) - Comprehensive feature documentation
- [examples/](examples/) - Example scripts demonstrating all features
  - `api_example.lango` - HTTP requests and JSON handling
  - `file_example.lango` - File I/O operations
  - `stdlib_example.lango` - Standard library functions
  - `complete_api_demo.lango` - Complete API client demo

## 🏗️ Architecture

Lango is a tree-walking interpreter written in Go:

- **`scanner.go`** - Lexical analysis (tokenization)
- **`parser.go`** - Syntax analysis (builds AST)
- **`interpreter.go`** - Tree-walking interpreter
- **`network.go`** - HTTP client implementation
- **`json.go`** - JSON parsing and serialization
- **`fileio.go`** - File system operations
- **`stdlib.go`** - Standard library functions
- **`environment.go`** - Variable scopes and environments
- **`expr.go`, `stmt.go`** - AST node definitions

### Performance

All networking, file I/O, and standard library functions are implemented as native Go functions, providing:
- Fast execution (near-native performance for I/O)
- Efficient memory usage
- Reliable error handling
- Thread-safe operations

## 🔧 Building from Source

**Requirements:**
- Go 1.23.3 or later

**Build:**
```bash
go mod init Lango  # First time only
go build
```

**Run tests:**
```bash
./Lango ./build/test_script.lango
```

## 🤝 Contributing

Contributions are welcome! Areas for improvement:

- [x] Add escape sequence support to scanner (✅ Completed)
- [x] Implement object literal syntax `{key: value}` (✅ Completed)
- [ ] Add async/await for non-blocking operations
- [ ] WebSocket support
- [ ] Database drivers (SQL, MongoDB)
- [ ] Regular expression support
- [ ] HTTP server capabilities
- [ ] Module/package system

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

Built following principles from "Crafting Interpreters" by Robert Nystrom.

## 🌟 Why Lango?

- **Simple Syntax** - Easy to learn and read
- **Powerful Features** - Networking, JSON, File I/O built-in
- **Go Performance** - Fast execution with native Go functions
- **Practical** - Designed for real-world scripting tasks
- **Extensible** - Easy to add new built-in functions

---

**Made with ❤️ using Go**
