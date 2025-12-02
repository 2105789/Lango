# Lango Language - Enhanced Features

This document describes the new features added to Lango in the latest update.

## Overview

Lango now includes comprehensive networking capabilities, JSON support, file I/O operations, and an extensive standard library, making it suitable for building APIs, web clients, and data processing applications.

## New Features

### 1. HTTP Networking

Lango now supports making HTTP requests with full support for GET, POST, PUT, and DELETE methods.

#### HTTP Functions

- **`httpGet(url, headers?)`** - Perform HTTP GET request
- **`httpPost(url, body, headers?)`** - Perform HTTP POST request  
- **`httpPut(url, body, headers?)`** - Perform HTTP PUT request
- **`httpDelete(url, headers?)`** - Perform HTTP DELETE request

#### Response Object

All HTTP functions return a response object with the following properties:
- `status` - HTTP status code (number)
- `statusText` - Status text (e.g., "200 OK")
- `body` - Response body as string
- `headers` - Response headers as map

#### Example

```lango
var response = httpGet("https://api.example.com/data");
if (response.status == 200) {
  print "Success: " + response.body;
}
```

### 2. JSON Support

Full JSON parsing and serialization capabilities.

#### JSON Functions

- **`jsonParse(jsonString)`** - Parse JSON string into Lango objects/arrays
- **`jsonStringify(value)`** - Convert Lango value to JSON string
- **`jsonStringifyPretty(value)`** - Convert to pretty-printed JSON string

#### Example

```lango
var data = jsonParse('{"name":"John","age":30}');
print data.name;  // Output: John
print data.age;   // Output: 30

var json = jsonStringify(data);
print json;  // Output: {"name":"John","age":30}
```

### 3. File I/O Operations

Comprehensive file system operations.

#### File Functions

- **`fileRead(path)`** - Read file content as string
- **`fileWrite(path, content)`** - Write string to file (overwrites)
- **`fileAppend(path, content)`** - Append content to file
- **`fileExists(path)`** - Check if file exists (returns boolean)
- **`fileDelete(path)`** - Delete a file
- **`fileList(directory)`** - List files in directory (returns array)
- **`fileCopy(srcPath, dstPath)`** - Copy a file

#### Example

```lango
fileWrite("data.txt", "Hello, Lango!");
var content = fileRead("data.txt");
print content;  // Output: Hello, Lango!

fileAppend("data.txt", " More text.");
fileDelete("data.txt");
```

### 4. String Functions

Rich string manipulation capabilities.

- **`strLen(str)`** - Get string length
- **`strUpper(str)`** - Convert to uppercase
- **`strLower(str)`** - Convert to lowercase
- **`strSplit(str, delimiter)`** - Split string into array
- **`strJoin(array, delimiter)`** - Join array into string
- **`strContains(str, substr)`** - Check if contains substring
- **`strReplace(str, old, new)`** - Replace all occurrences
- **`strTrim(str)`** - Remove leading/trailing whitespace

#### Example

```lango
var text = "  Hello, World!  ";
print strTrim(text);           // Output: Hello, World!
print strUpper(text);          // Output:   HELLO, WORLD!  
print strLen(text);            // Output: 17

var parts = strSplit("a,b,c", ",");
print arrayLen(parts);         // Output: 3
```

### 5. Array Functions

Enhanced array operations.

- **`arrayLen(array)`** - Get array length
- **`arrayPush(array, item)`** - Add item to end (returns new array)
- **`arrayPop(array)`** - Get last item
- **`arraySlice(array, start, end)`** - Get array slice
- **`arrayConcat(array1, array2)`** - Concatenate arrays

#### Example

```lango
var arr = [1, 2, 3];
print arrayLen(arr);           // Output: 3

var extended = arrayPush(arr, 4);
print extended;                // Output: [1, 2, 3, 4]

var sliced = arraySlice(arr, 0, 2);
print sliced;                  // Output: [1, 2]
```

### 6. Math Functions

Comprehensive mathematical operations.

- **`mathFloor(num)`** - Round down
- **`mathCeil(num)`** - Round up
- **`mathRound(num)`** - Round to nearest integer
- **`mathAbs(num)`** - Absolute value
- **`mathMax(a, b)`** - Maximum of two numbers
- **`mathMin(a, b)`** - Minimum of two numbers
- **`mathPow(base, exponent)`** - Power function
- **`mathSqrt(num)`** - Square root
- **`mathRandom()`** - Random number between 0 and 1

#### Example

```lango
print mathFloor(4.7);          // Output: 4
print mathCeil(4.3);           // Output: 5
print mathRound(4.5);          // Output: 5
print mathPow(2, 8);           // Output: 256
print mathSqrt(16);            // Output: 4
```

### 7. Type Conversion

Convert between different types.

- **`toNumber(value)`** - Convert to number
- **`toString(value)`** - Convert to string
- **`toBoolean(value)`** - Convert to boolean

#### Example

```lango
var num = toNumber("42");
print num + 8;                 // Output: 50

var str = toString(123);
print "Number: " + str;        // Output: Number: 123
```

### 8. Utility Functions

Additional helpful utilities.

- **`time()`** - Get current Unix timestamp
- **`sleep(milliseconds)`** - Sleep for specified milliseconds

#### Example

```lango
var timestamp = time();
print timestamp;

print "Waiting...";
sleep(1000);  // Wait 1 second
print "Done!";
```

## Examples

The `examples/` directory contains comprehensive demonstrations:

- **`api_example.lango`** - HTTP requests and JSON handling
- **`file_example.lango`** - File I/O operations
- **`stdlib_example.lango`** - Standard library functions

Run examples with:
```bash
.\Lango.exe .\examples\api_example.lango
.\Lango.exe .\examples\file_example.lango
.\Lango.exe .\examples\stdlib_example.lango
```

## Building an API Client

Here's a complete example of building an API client:

```lango
// Fetch user data
var response = httpGet("https://jsonplaceholder.typicode.com/users/1");

if (response.status == 200) {
  var user = jsonParse(response.body);
  
  print "User: " + user.name;
  print "Email: " + user.email;
  print "City: " + user.address.city;
  
  // Save to file
  var formatted = jsonStringifyPretty(user);
  fileWrite("user.json", formatted);
  print "Saved user data to user.json";
}
```

## Technical Implementation

### Module Structure

The enhanced Lango interpreter includes these new Go modules:

- **`network.go`** - HTTP client implementation
- **`json.go`** - JSON parsing and serialization
- **`fileio.go`** - File system operations
- **`stdlib.go`** - Standard library functions

### Built-in Function Registration

All functions are registered as native Go functions in the global environment during interpreter initialization, providing optimal performance while maintaining Lango's simplicity.

### JSON Object Access

Parsed JSON objects are represented as Go `map[string]interface{}` and support property access using dot notation, just like Lango class instances.

## Performance

All networking, file I/O, and standard library functions are implemented in Go, providing:
- Native Go performance
- Efficient memory usage
- Reliable error handling
- Thread-safe operations

## Limitations

- JSON strings in Lango code currently require avoiding escape sequences
- No object literal syntax (use JSON strings instead)
- Synchronous HTTP requests only (no async/await)

## Future Enhancements

Potential next-phase features:
- Async/await support for network operations
- WebSocket support
- Database drivers (SQL, NoSQL)
- Regular expression support
- More advanced string operations
- Module/package system

## Conclusion

With these enhancements, Lango is now a capable language for:
- Building REST API clients
- Data processing and transformation
- File-based operations
- Web scraping and automation
- Configuration management
- System scripting

The combination of Go's performance with Lango's simple syntax makes it ideal for rapid development of network-enabled applications.
