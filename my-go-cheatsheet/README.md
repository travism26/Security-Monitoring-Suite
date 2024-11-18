# Go Language Cheatsheet

## Table of Contents

- [Basic Syntax](#basic-syntax)
- [Data Types](#data-types)
- [Variables](#variables)
- [Control Structures](#control-structures)
- [Functions](#functions)
- [Structs and Interfaces](#structs-and-interfaces)
- [Concurrency](#concurrency)
- [Error Handling](#error-handling)
- [Packages and Modules](#packages-and-modules)

## Basic Syntax

### Hello World

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
```

## Data Types

- **Basic Types**

  - `bool`
  - `string`
  - `int`, `int8`, `int16`, `int32`, `int64`
  - `uint`, `uint8`, `uint16`, `uint32`, `uint64`
  - `float32`, `float64`
  - `complex64`, `complex128`
  - `byte` (alias for uint8)
  - `rune` (alias for int32)

- **Composite Types**
  - Arrays
  - Slices
  - Maps
  - Structs
  - Pointers
  - Channels

### Type Conversion

```go
// Basic type conversion
i := 42
f := float64(i)    // int to float64 => 42.00
s := string(i)     // int to string (converts to ASCII/Unicode) => "42"
b := byte(i)       // int to byte => 42

// String conversions
str := "123"
num, err := strconv.Atoi(str)      // string to int => 123
num64, err := strconv.ParseInt(str, 10, 64)  // string to int64 => 123
f64, err := strconv.ParseFloat(str, 64)      // string to float64 => 123.00
byteSlice := []byte(str)                     // string to byte slice => [49, 50, 51]

// Converting back to string
str = strconv.Itoa(num)            // int to string => "123"
str = strconv.FormatInt(num64, 10) // int64 to string => "123"
str = strconv.FormatFloat(f64, 'f', 2, 64)   // float64 to string => "123.00"

// Array/Slice conversion
slice := []int{1, 2, 3}
array := [3]int(slice)             // slice to array => [1, 2, 3]
slice2 := array[:]                 // array to slice => [1, 2, 3]
strArray := []string(slice)        // slice to string array => ["1", "2", "3"]
byteSlice2 := []byte(strArray)     // string array to byte slice => [49, 50, 51]
```

## Variables

### Declaration

```go
var name string = "John"
var age = 25 // Type inference
shorthand := "value" // Short declaration
const PI = 3.14159
```

### Zero Values

- Numbers: `0`
- Booleans: `false`
- Strings: `""`
- Pointers: `nil`

## Control Structures

### If Statement

```go
if x > 0 {
    // code
} else if x < 0 {
    // code
} else {
    // code
}
```

### For Loop

```go
// Standard for loop
for i := 0; i < 10; i++ {
    // code
}

// While-style loop
for condition {
    // code
}

// Range loop
for index, value := range slice {
    // code
}
```

### Switch

```go
switch value {
case 1:
    // code
case 2:
    // code
default:
    // code
}
```

## Functions

### Basic Function

```go
func add(x int, y int) int {
    return x + y
}
```

### Multiple Return Values

```go
func divide(x, y float64) (float64, error) {
    if y == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return x / y, nil
}
```

### Defer

```go
func example() {
    defer fmt.Println("This runs last")
    fmt.Println("This runs first")
}
```

## Structs and Interfaces

### Struct

```go
type Person struct {
    Name string
    Age  int
}

// Method
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, my name is %s", p.Name)
}
```

### Interface

```go
type Animal interface {
    Speak() string
}
```

## Concurrency

### Goroutines

```go
go function() // Start a new goroutine

// Anonymous function
go func() {
    // code
}()
```

### Channels

```go
ch := make(chan int)    // Create channel
ch <- value            // Send value
value := <-ch         // Receive value

// Buffered channel
ch := make(chan int, 100)
```

### Select

```go
select {
case msg1 := <-ch1:
    // Handle msg1
case msg2 := <-ch2:
    // Handle msg2
default:
    // Optional default case
}
```

## Error Handling

```go
if err != nil {
    return err
}

// Custom error
type CustomError struct {
    message string
}

func (e *CustomError) Error() string {
    return e.message
}
```

## Packages and Modules

### Module Initialization

```bash
go mod init module-name
```

### Common Imports

```go
import (
    "fmt"      // Formatted I/O
    "os"       // Operating system functionality
    "strings"  // String manipulation
    "time"     // Time functionality
    "errors"   // Error handling
    "context"  // Context management
)
```

### Testing

```go
// file: example_test.go
func TestFunction(t *testing.T) {
    got := Function()
    want := expectedValue
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

## Best Practices

- Use `gofmt` to format your code
- Handle errors explicitly
- Use meaningful variable names
- Write documentation comments
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `go vet` and `golint` for code analysis

## Common Commands

```bash
go run file.go      # Run a program
go build            # Build the program
go test            # Run tests
go get package     # Download and install packages
go mod tidy        # Clean up module dependencies
```
