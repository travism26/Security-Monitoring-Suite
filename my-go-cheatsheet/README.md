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
- [Context](#context)
- [Generics](#generics)
- [Reflection](#reflection)
- [Testing](#testing)
- [Common Stdlib Packages](#common-stdlib-packages)
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

// Array of strings to string
strArray := []string{"1", "2", "3"}
str := strings.Join(strArray, ",") // join string array => "1,2,3"
byteSlice3 := []byte(str)          // string to byte slice => [49, 44, 50, 44, 51]
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

### Basic Interface

```go
// Simple interface definition
type Writer interface {
    Write([]byte) (int, error)
}

// Interface with multiple methods
type ReadWriter interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
}
```

### Interface Composition

```go
// Interfaces can be composed of other interfaces
type ReadWriter interface {
    Reader
    Writer
}

// Real-world example
type SystemMonitor interface {
    CPUMonitor
    MemoryMonitor
    DiskMonitor
}
```

### Interface Implementation

```go
// Interface
type Animal interface {
    Speak() string
}

// Implicit implementation
type Dog struct {
    Name string
}

// Dog implements Animal interface implicitly without declaring it unlike java
func (d Dog) Speak() string {
    return fmt.Sprintf("%s says Woof!", d.Name)
}

// Usage
func main() {
    var animal Animal = Dog{Name: "Rex"}
    fmt.Println(animal.Speak())  // Output: Rex says Woof!
}
```

### Interface Composition with Structs

```go
// Interfaces
type Logger interface {
    Log(string)
}

type Processor interface {
    Process() error
}

// Implementation using composition
type Worker struct {
    logger    Logger    // Composition through embedding
    processor Processor
}

// Constructor pattern
func NewWorker(l Logger, p Processor) *Worker {
    return &Worker{
        logger:    l,
        processor: p,
    }
}

// Methods using composed interfaces
func (w *Worker) DoWork() error {
    w.logger.Log("Starting work")
    return w.processor.Process()
}
```

### Empty Interface

```go
// Empty interface can hold any type
var i interface{}
i = 42          // holds an int
i = "hello"     // holds a string
i = struct{}{}  // holds a struct

// Type assertion
str, ok := i.(string)
if ok {
    fmt.Printf("Value is a string: %s\n", str)
}

// Type switch
switch v := i.(type) {
case int:
    fmt.Printf("Integer: %d\n", v)
case string:
    fmt.Printf("String: %s\n", v)
default:
    fmt.Printf("Unknown type\n")
}
```

### Interface Best Practices

```go
// Good: Small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Good: Interface segregation
type FileHandler interface {
    Open() error
    Close() error
}

type FileReader interface {
    FileHandler
    Reader
}

// Bad: Large, monolithic interface
type DoEverything interface {
    DoThis() error
    DoThat() error
    DoSomethingElse() error
    // ... many more methods
}
```

### Common Design Patterns with Interfaces

```go
// Factory Pattern
type Logger interface {
    Log(string)
}

func NewLogger(logType string) Logger {
    switch logType {
    case "file":
        return &FileLogger{}
    case "console":
        return &ConsoleLogger{}
    default:
        return &NullLogger{}
    }
}

// Strategy Pattern
type SortStrategy interface {
    Sort([]int)
}

type Sorter struct {
    strategy SortStrategy
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
    s.strategy = strategy
}

func (s *Sorter) Sort(data []int) {
    s.strategy.Sort(data)
}
```

### Testing with Interfaces

```go
// Interface for testing
type DataStore interface {
    Save(data []byte) error
    Load() ([]byte, error)
}

// Mock implementation for testing
type MockDataStore struct {
    data []byte
    err  error
}

func (m *MockDataStore) Save(data []byte) error {
    if m.err != nil {
        return m.err
    }
    m.data = data
    return nil
}

// Test using mock
func TestDataProcessor(t *testing.T) {
    mock := &MockDataStore{}
    processor := NewDataProcessor(mock)
    // ... test implementation
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

### Sync Package

```go
// WaitGroup for synchronizing goroutines
var wg sync.WaitGroup
wg.Add(1)  // Add a counter
go func() {
    defer wg.Done()  // Decrement counter when done
    // Do work
}()
wg.Wait()  // Wait for all goroutines to finish

// Mutex for protecting shared resources
var mu sync.Mutex
mu.Lock()
// Critical section
mu.Unlock()

// RWMutex for read/write locks
var rwmu sync.RWMutex
rwmu.RLock()  // Multiple readers can acquire lock
// Read operations
rwmu.RUnlock()

rwmu.Lock()   // Only one writer can acquire lock
// Write operations
rwmu.Unlock()

// Once for one-time initialization
var once sync.Once
once.Do(func() {
    // This will only execute once
})
```

## Error Handling

### Basic Error Handling

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

### Error Wrapping (Go 1.13+)

```go
// Wrapping errors
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Unwrapping errors
var customErr *CustomError
if errors.As(err, &customErr) {
    // Handle custom error
}

// Check error type
if errors.Is(err, io.EOF) {
    // Handle EOF
}
```

### Error Best Practices

```go
// Define error types for specific cases
var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid input")
)

// Use error variables for comparison
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}

// Custom error types with additional context
type ValidationError struct {
    Field string
    Error error
}

func (v *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %v", v.Field, v.Error)
}

func (v *ValidationError) Unwrap() error {
    return v.Error
}
```

## Context

### Basic Context Usage

```go
// Creating context
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// Context with value
ctx = context.WithValue(ctx, "key", "value")
value := ctx.Value("key").(string)

// Using context in functions
func DoWork(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(2 * time.Second):
        return nil
    }
}
```

### Context Patterns

```go
// HTTP Server with context
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // Use context for timeouts, cancellation, etc.
}

// Database operations with context
func (db *DB) QueryWithContext(ctx context.Context, query string) (*Result, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Perform query
    }
}

// Graceful shutdown
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := server.Shutdown(ctx); err != nil {
    log.Fatal("Server forced to shutdown:", err)
}
```

## Generics (Go 1.18+)

### Basic Generic Types

```go
// Generic function
func Min[T constraints.Ordered](x, y T) T {
    if x < y {
        return x
    }
    return y
}

// Generic data structure
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
    var zero T
    if len(s.items) == 0 {
        return zero, errors.New("empty stack")
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, nil
}
```

### Type Constraints

```go
// Custom constraint
type Number interface {
    ~int | ~float64
}

// Generic function with constraint
func Sum[T Number](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}

// Multiple type parameters
func Map[T, U any](s []T, f func(T) U) []U {
    r := make([]U, len(s))
    for i, v := range s {
        r[i] = f(v)
    }
    return r
}
```

## Reflection

### Basic Reflection

```go
// Get type information
t := reflect.TypeOf(x)
v := reflect.ValueOf(x)

// Get/Set values
if v.Kind() == reflect.Ptr && !v.IsNil() {
    v = v.Elem()
}

if v.CanSet() {
    v.SetString("new value")
}

// Iterate struct fields
t := reflect.TypeOf(struct{
    Name string
    Age  int
}{})

for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Printf("Field: %s, Type: %s\n", field.Name, field.Type)
}
```

### Reflection Use Cases

```go
// Dynamic method calls
func CallMethod(v interface{}, method string, args ...interface{}) []reflect.Value {
    return reflect.ValueOf(v).MethodByName(method).Call(
        MakeValueSlice(args...))
}

// Struct tag parsing
type Person struct {
    Name string `json:"name" validate:"required"`
    Age  int    `json:"age" validate:"gte=0,lte=130"`
}

func ParseTags(v interface{}) map[string]string {
    t := reflect.TypeOf(v)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    tags := make(map[string]string)
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        tags[field.Name] = field.Tag.Get("json")
    }
    return tags
}
```

## Testing

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        x, y     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.x, tt.y)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.x, tt.y, got, tt.expected)
            }
        })
    }
}
```

### Benchmarking

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// Sub-benchmarks
func BenchmarkSort(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            data := make([]int, size)
            for i := 0; i < b.N; i++ {
                b.StopTimer() // Don't measure setup
                rand.Shuffle(len(data), func(i, j int) {
                    data[i], data[j] = data[j], data[i]
                })
                b.StartTimer() // Measure sort
                sort.Ints(data)
            }
        })
    }
}
```

### Test Helpers

```go
// Helper function
func setupTestCase(t *testing.T) func() {
    t.Helper() // Marks this as a helper function

    // Setup
    tmpDir, err := ioutil.TempDir("", "test")
    if err != nil {
        t.Fatal(err)
    }

    // Return cleanup function
    return func() {
        os.RemoveAll(tmpDir)
    }
}

// Using helper
func TestSomething(t *testing.T) {
    cleanup := setupTestCase(t)
    defer cleanup()

    // Test code
}
```

## Common Stdlib Packages

### io/ioutil and os

```go
// Reading files
data, err := ioutil.ReadFile("file.txt")
content := string(data)

// Writing files
err = ioutil.WriteFile("file.txt", []byte("content"), 0644)

// Directory operations
files, err := ioutil.ReadDir(".")
for _, f := range files {
    fmt.Printf("Name: %s, Size: %d\n", f.Name(), f.Size())
}

// File operations
f, err := os.OpenFile("file.txt", os.O_RDWR|os.O_CREATE, 0644)
defer f.Close()
```

### encoding/json

```go
// Marshal
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

p := Person{Name: "John", Age: 30}
data, err := json.Marshal(p)

// Unmarshal
var p2 Person
err = json.Unmarshal(data, &p2)

// Streaming JSON
dec := json.NewDecoder(reader)
for dec.More() {
    var m map[string]interface{}
    err := dec.Decode(&m)
    // Process m
}
```

### net/http

```go
// Simple server
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
})
http.ListenAndServe(":8080", nil)

// HTTP client
resp, err := http.Get("http://example.com")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
```

### time

```go
// Time operations
now := time.Now()
future := now.Add(24 * time.Hour)
duration := future.Sub(now)

// Formatting
formatted := now.Format("2006-01-02 15:04:05")
parsed, err := time.Parse("2006-01-02", "2023-01-01")

// Timers
timer := time.NewTimer(2 * time.Second)
<-timer.C // Wait for timer

// Tickers
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for range ticker.C {
    // Do something every second
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
