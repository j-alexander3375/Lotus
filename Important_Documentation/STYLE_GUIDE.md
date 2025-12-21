# Lotus Language Style Guide

## Philosophy

Lotus is a systems programming language that bridges **Rust's memory safety philosophy**, **C++'s performance and control**, and **Go's simplicity and clarity**. The style guide promotes **modularity, safety, and readability** while remaining distinctive from its influences.

---

## Core Principles

1. **Explicit Over Implicit** (Rust influence)
   - Make intentions clear in code
   - Avoid magic or hidden behavior
   - Use explicit keywords and types

2. **Safety By Default** (Rust influence)
   - Memory safety through static analysis
   - Type safety enforced at compile time
   - Clear ownership and lifetime semantics

3. **Performance Matters** (C++ influence)
   - No unnecessary abstractions
   - Direct hardware access when needed
   - Zero-cost abstractions

4. **Clarity Over Cleverness** (Go influence)
   - Code is read more often than written
   - Simple, straightforward syntax
   - Minimize cognitive load

5. **Modularity First**
   - Clear module boundaries
   - Explicit dependencies via imports
   - Easy to understand relationships between modules

---

## File Organization

### Naming Conventions

**Files:**
- Use `snake_case.lts` for Lotus source files
- Shorter names for simple modules: `io.lts`, `math.lts`
- Descriptive names for implementations: `string_utils.lts`, `memory_allocator.lts`
- Module directories use `snake_case`: `std::collections`, `std::io`

**Modules:**
```lotus
// Good: Clear, descriptive
use "string::formatting";
use "memory::allocator";
use "collections::vector";

// Avoid: Too generic
use "util";
use "helpers";
```

### Module Structure

**Single Responsibility Principle**

```lotus
// math.lts - GOOD: Single purpose
use "io";

fn abs(x: int) -> int {
    if x < 0 {
        ret -x;
    }
    ret x;
}

fn max(a: int, b: int) -> int {
    if a > b { ret a; }
    ret b;
}
```

```lotus
// utils.lts - BAD: Mixed responsibilities
use "io";
use "mem";

fn abs(x: int) -> int { /* ... */ }
fn print_vector(v: int[]) { /* ... */ }
fn malloc_aligned(size: int, align: int) { /* ... */ }
fn format_string(s: string) -> string { /* ... */ }
```

---

## Naming Style

### Quick Reference

| Element | Style | Example |
|---------|-------|---------|
| Functions | `snake_case` | `fn calculate_fibonacci(n: int)` |
| Constants | `UPPER_SNAKE_CASE` | `const int MAX_VECTOR_SIZE = 1000;` |
| Structs | `snake_case` | `struct vector { /* ... */ }` |
| Enums | `snake_case` | `enum result { ok, error }` |
| Enum Variants | `snake_case` | `result.ok` |
| Variables | `snake_case` | `int result = 0;` |
| Module Names | `lowercase` | `use "io"; use "math";` |
| Type Names | `lowercase` | `int`, `string`, `bool` |
| Boolean Vars | `snake_case` + prefix | `is_valid`, `has_error`, `can_allocate` |

**Lotus-only cues**
- Prefer `ret expr;` over implicit returns.
- Bindings are type-first by default: `int count = 0;` instead of `let count: int`.
- Imports stay string-based (`use "io::printf";`) to keep module boundaries explicit.

### Detailed Identifiers

**Functions: `snake_case`**
```lotus
fn calculate_fibonacci(n: int) -> int {
    int result = 0;
    int prev = 0;
    int curr = 1;

    while result < n {
        prev = curr;
        curr = result;
        result = prev + curr;
    }

    ret result;
}

fn main() {
    int fib_value = calculate_fibonacci(10);
    println("Result: %d", fib_value);
}
```

**Variables: `snake_case`**
```lotus
fn calculate_fibonacci(n: int) -> int {
    int result = 0;
    int prev = 0;
    int curr = 1;

    while result < n {
        prev = curr;
        curr = result;
        result = prev + curr;
    }

    ret result;
}
```

**Structs: `snake_case`**
```lotus
struct vector {
    data: *int,
    capacity: int,
    length: int,
}

class linked_list_node {
    value: int,
    next: *linked_list_node,
}
```

**Enums: `snake_case`**
```lotus
enum result {
    ok,
    error,
}
```

**Constants: `UPPER_SNAKE_CASE`**
```lotus
const int MAX_VECTOR_SIZE = 1000;
const string DEFAULT_ENCODING = "utf-8";
const bool ENABLE_DEBUG_LOGGING = false;
```

**Module Names: `lowercase` (no underscores in standard library)**
```lotus
use "io";        // Good
use "math";      // Good
use "io_utils";  // Avoid in stdlib (use submodule instead)
```

### Boolean Variables: Be Explicit

```lotus
// Good: Clear what state is being checked
bool is_valid = check_validity(data);
bool has_error = error != null;
bool can_allocate = available_memory > size;
bool is_initialized = ptr != null;

if is_valid {
    process_data(data);
}

// Avoid: Ambiguous
bool valid = check_validity(data);
bool error = error != null;  // Confusing - is error a boolean?
bool result = check(data);
```

---

## Formatting and Layout

### Indentation and Spacing

**Use 4 spaces for indentation** (not tabs)

```lotus
fn process_data(input: string) -> bool {
    use "io";
    use "string";
    
    int length = len(input);

    if length > 0 {
        string processed = format_string(input);
        println("Processed: %s", processed);
        ret true;
    }

    println("Error: Empty input");
    ret false;
}
```

### Line Length

**Maximum 100 characters per line**
- Promotes readability
- Works on most displays
- Encourages clearer code

```lotus
// Good: Clear, readable
fn calculate_complex_value(
    first_parameter: int,
    second_parameter: int,
    third_parameter: int
) -> int {
    first_parameter + second_parameter + third_parameter
}

// Avoid: Line too long
fn calculate_complex_value(first_parameter: int, second_parameter: int, third_parameter: int) -> int { first_parameter + second_parameter + third_parameter }
```

### Brace Placement

**Opening braces on the same line** (influenced by Go)

```lotus
// Good
fn calculate(x: int) -> int {
    if x > 0 {
        x * 2
    } else {
        0
    }
}

// Avoid
fn calculate(x: int) -> int
{
    if x > 0
    {
        x * 2
    }
}
```

### Blank Lines

- One blank line between module declarations
- One blank line between function definitions
- Two blank lines between major sections

```lotus
use "io";
use "math";

const int MAX_SIZE = 100;

// ============================================================================
// String operations
// ============================================================================

fn reverse_string(s: string) -> string {
    // Implementation
}

fn concatenate_strings(s1: string, s2: string) -> string {
    // Implementation
}

// ============================================================================
// Memory utilities
// ============================================================================

fn allocate_aligned(size: int, alignment: int) -> *int {
    // Implementation
}
```

---

## Comments and Documentation

### Inline Comments

```lotus
fn fibonacci(n: int) -> int {
    // Base cases for Fibonacci sequence
    if n <= 1 {
        ret n;
    }
    
    // Use iterative approach for better performance
    int prev = 0;
    int curr = 1;
    
    while n > 1 {
        int next = prev + curr;
        prev = curr;
        curr = next;
        n = n - 1;
    }
    
    ret curr;
}
```

### Function Documentation

```lotus
/// Calculates the absolute value of an integer.
/// 
/// Arguments:
///   x - The input integer
/// 
/// Returns:
///   The absolute value of x
/// 
/// Example:
///   int result = abs(-42);  // Returns 42
fn abs(x: int) -> int {
    if x < 0 { ret -x; }
    ret x;
}
```

### Module-level Comments

```lotus
/// The math module provides common mathematical operations.
///
/// This module includes:
/// - Basic arithmetic (abs, min, max)
/// - Exponential and power functions
/// - Trigonometric functions (future)
/// 
/// All functions use integer arithmetic unless otherwise specified.

use "io";

// Implementation follows...
```

---

## Type System and Declarations

### Type Clarity

**Always be explicit about types**

```lotus
// Good: Clear types everywhere
fn process_list(items: *int, count: int) -> int {
    int total = 0;
    int index = 0;
    
    while index < count {
        total = total + items[index];
        index = index + 1;
    }
    
    ret total;
}

// Avoid: Implicit or unclear types
fn process_list(items, count) {
    total = 0;
    for i in 0..count {
        total += items[i];
    }
    return total;
}
```

### Pointer Usage (Rust-influenced clarity)

```lotus
// Good: Clear when dealing with pointers
fn allocate_and_initialize(size: int) -> *int {
    int* ptr = malloc(size * sizeof(int));
    
    if ptr == null {
        ret null;  // Signal allocation failure
    }
    
    // Initialize the allocated memory
    int i = 0;
    while i < size {
        ptr[i] = 0;
        i = i + 1;
    }
    
    ret ptr;
}

// Avoid: Unclear pointer semantics
fn allocate_and_initialize(size) {
    ptr = malloc(size);
    for i in 0..size {
        *ptr = 0;
    }
    return ptr;
}
```

### Null Handling

**Always consider null explicitly**

```lotus
// Good: Explicit null checking
fn safe_get_value(items: *int, index: int) -> int {
    if items == null {
        ret -1;  // Error indicator
    }
    
    if index < 0 {
        ret -1;  // Invalid index
    }

    ret items[index];
}

// Future: Use Result type
fn safe_get_value(items: *int, index: int) -> result<int, string> {
    if items == null {
        Err("Null pointer")
    } else if index < 0 {
        Err("Invalid index")
    } else {
        Ok(items[index])
    }
}
```

---

## Control Flow

### If/Else Statements

**Prefer clear conditions over nested ternaries**

```lotus
// Good: Readable and clear
if value > MAX_THRESHOLD {
    HANDLE_LARGE_VALUE(value);
} else if value > NORMAL_THRESHOLD {
    HANDLE_NORMAL_VALUE(value);
} else {
    HANDLE_SMALL_VALUE(value);
}

// Avoid: Hard to read
int result = (value > MAX ? HANDLE_LARGE : (value > NORMAL ? HANDLE_NORMAL : HANDLE_SMALL))(value);
```

### Loops

**Prefer `while` loops for clarity**

```lotus
// Good: Explicit iteration
fn sum_array(items: *int, count: int) -> int {
    int total = 0;
    int i = 0;
    
    while i < count {
        total = total + items[i];
        i = i + 1;
    }
    
    ret total;
}

// Also acceptable: For-loop (when available)
fn sum_array_for(items: *int, count: int) -> int {
    int total = 0;
    
    for i in 0..count {
        total = total + items[i];
    }
    
    ret total;
}
```

### Guard Clauses

**Use early returns to avoid nesting**

```lotus
// Good: Flat control flow
fn validate_input(name: string, age: int) -> bool {
    if name == null || len(name) == 0 {
        println("Error: Name cannot be empty");
        ret false;
    }
    
    if age < 0 || age > 150 {
        println("Error: Age must be 0-150");
        ret false;
    }
    
    // Main logic here
    process_user(name, age);
    ret true;
}

// Avoid: Deeply nested conditions
fn validate_input(name: string, age: int) -> bool {
    if name != null && len(name) > 0 {
        if age >= 0 && age <= 150 {
            process_user(name, age);
            ret true;
        } else {
            println("Error: Age invalid");
            ret false;
        }
    } else {
        println("Error: Name empty");
        ret false;
    }
}
```

---

## Functions and Methods

### Function Signatures

**Be explicit about inputs and outputs**

```lotus
// Good: Clear what goes in and comes out
fn calculate_average(
    values: *int,
    count: int
) -> int {
    if count == 0 {
        ret 0;
    }
    
    int total = 0;
    int i = 0;
    
    while i < count {
        total = total + values[i];
        i = i + 1;
    }
    
    ret total / count;
}

// Avoid: Unclear what this does
fn calc(v, n) {
    int t = 0;
    for i in 0..n {
        t = t + v[i];
    }
    ret t / n;
}
```

### Return Values

**Single, clear return point when possible**

```lotus
// Good: Single return
fn find_max(values: *int, count: int) -> int {
    if count == 0 {
        ret -1;
    }
    
    int max = values[0];
    int i = 1;
    
    while i < count {
        if values[i] > max {
            max = values[i];
        }
        i = i + 1;
    }
    
    ret max;
}

// Acceptable: Multiple returns for error handling
fn allocate_memory(size: int) -> *int {
    if size <= 0 {
        ret null;
    }
    
    int* ptr = malloc(size);
    
    if ptr == null {
        ret null;  // Allocation failed
    }
    
    ret ptr;
}
```

---

## Imports and Module Organization

### Import Ordering

```lotus
// 1. Standard library imports
use "io";
use "math";
use "mem";

// 2. Specific imports
use "io::printf";
use "math::max";

// 3. Code organization
struct my_data {
    value: int,
    name: string,
}

fn process_data(data: *my_data) {
    // Implementation
}
```

### Avoiding Circular Dependencies

```lotus
// Good: Clear dependency direction
// module_a.lts
use "module_b";

fn function_a() {
    function_b();
}

// module_b.lts
// Does NOT use module_a
fn function_b() {
    // Independent implementation
}

// Avoid: Circular dependencies
// module_a.lts
use "module_b";

// module_b.lts
use "module_a";  // Creates circular dependency!
```

---

## Memory Management (Systems Programming Focus)

### Ownership and Borrowing (Rust-influenced)

**Make ownership explicit:**

```lotus
// Good: Clear ownership transfer
fn create_and_populate_array(size: int) -> *int {
    int* array = malloc(size * sizeof(int));
    
    if array == null {
        ret null;  // Failed to allocate
    }
    
    // Initialize
    int i = 0;
    while i < size {
        array[i] = i * 2;
        i = i + 1;
    }
    
    ret array;  // Ownership transferred to caller
}

fn main() {
    use "io";
    use "mem";
    
    int* my_array = create_and_populate_array(10);
    
    // Use the array...
    
    free(my_array);  // Caller is responsible for cleanup
}
```

### Resource Cleanup

```lotus
// Pattern: Allocate-Use-Free
fn process_file_data(filename: string) {
    use "mem";
    
    int* buffer = malloc(1024);
    
    if buffer == null {
        println("Error: Failed to allocate buffer");
        ret;
    }
    
    // Use buffer...
    // Process data
    
    free(buffer);  // Always clean up
}

// Pattern: Guard to ensure cleanup
fn safe_operation(size: int) -> bool {
    use "mem";
    
    int* temp = malloc(size);
    
    if temp == null {
        ret false;  // Failed early
    }
    
    // Do work...
    bool result = do_work(temp);
    
    free(temp);
    ret result;
}
```

---

## Error Handling

### Null Checks (Current approach)

```lotus
// Good: Explicit null handling
fn divide(a: int, b: int) -> *int {
    if b == 0 {
        println("Error: Division by zero");
        ret null;
    }
    
    int* result = malloc(sizeof(int));
    
    if result == null {
        println("Error: Memory allocation failed");
        ret null;
    }
    
    result[0] = a / b;
    ret result;
}

fn main() {
    use "io";
    
    int* answer = divide(10, 2);
    
    if answer != null {
        printf("Result: %d\n", answer[0]);
        // free(answer);  // In real code, cleanup
    } else {
        println("Operation failed");
    }
}
```

### Result Types (Future enhancement)

```lotus
// Future style with Result type
enum result<T, E> {
    ok(T),
    err(E),
}

fn divide(a: int, b: int) -> result<int, string> {
    if b == 0 {
        err("Division by zero")
    } else {
        ok(a / b)
    }
}

fn main() {
    use "io";
    
    result<int, string> result = divide(10, 2);
    
    match result {
        ok(value) => printf("Result: %d\n", value),
        err(msg) => printf("Error: %s\n", msg),
    }
}
```

---

## Struct and Class Design

### Struct Definition and Usage (Go-influenced)

```lotus
// Good: Clear, minimal struct
struct point {
    x: int,
    y: int,
}

fn distance(p1: point, p2: point) -> int {
    int dx = p2.x - p1.x;
    int dy = p2.y - p1.y;
    
    // Simplified distance (not real Pythagorean)
    ret dx + dy;
}

fn main() {
    point p1 = { x: 0, y: 0 };
    point p2 = { x: 3, y: 4 };
    
    printf("Distance: %d\n", distance(p1, p2));
}
```

### Class Definition (C++-influenced)

```lotus
// Good: Clear class with encapsulation
class linked_list {
    head: *node,
    size: int,
}

impl linked_list {
    fn new() -> linked_list {
        linked_list {
            head: null,
            size: 0,
        }
    }
    
    fn push(this: *linked_list, value: int) {
        // Implementation
    }
    
    fn pop(this: *linked_list) -> int {
        // Implementation
    }
}
```

### Enums (Rust-influenced)

```lotus
// Good: Clear enumeration with variants
enum status {
    running,
    paused,
    stopped,
    error,
}

fn handle_status(s: status) {
    match s {
        status.running => println("Process is running"),
        status.paused => println("Process is paused"),
        status.stopped => println("Process stopped"),
        status.error => println("Error occurred"),
    }
}
```

---

## Constants and Immutability

### Constant Declaration

```lotus
// Good: Clear constant declarations
const int BUFFER_SIZE = 1024;
const string APP_NAME = "Lotus Compiler";
const bool RELEASE_BUILD = true;

fn main() {
    use "io";
    
    printf("Buffer size: %d\n", BUFFER_SIZE);
    printf("Application: %s\n", APP_NAME);
}
```

### Compile-time Constants vs Runtime

```lotus
// Compile-time constant (known at parse time)
const int MAX_STATIC_SIZE = 1000;

// Runtime constant (computed but immutable)
const int CALCULATED_SIZE = COMPUTE_OPTIMAL_SIZE(100);

fn COMPUTE_OPTIMAL_SIZE(base: int) -> int {
    base * 10
}
```

---

## Performance Considerations

### Inlining Hints

```lotus
// Small functions that should be inlined
fn min(a: int, b: int) -> int {
    if a < b { ret a; }
    ret b;
}

fn max(a: int, b: int) -> int {
    if a > b { ret a; }
    ret b;
}

// Larger function - let compiler decide
fn calculate_complex_statistics(data: *int, count: int) -> int {
    // Multiple operations - may or may not be inlined
    int sum = 0;
    int i = 0;
    
    while i < count {
        sum = sum + data[i];
        i = i + 1;
    }
    
    ret sum / count;
}
```

### Stack vs Heap

```lotus
// Good: Stack allocation for known, small sizes
fn process_small_set() {
    int small_array[10];  // Stack allocated
    int i = 0;
    
    while i < 10 {
        small_array[i] = i;
        i = i + 1;
    }
}

// Heap allocation for dynamic/large sizes
fn process_large_set(size: int) {
    int* large_array = malloc(size * sizeof(int));
    
    if large_array == null {
        ret;
    }
    
    // Use array...
    
    free(large_array);
}
```

---

## Testing and Validation

### Writing Testable Code

```lotus
// Good: Functional, testable
fn validate_email(email: string) -> bool {
    if email == null || len(email) == 0 {
        ret false;
    }
    
    // Simple validation (real implementation more complex)
    bool has_at = contains_char(email, '@');
    bool has_dot = contains_char(email, '.');
    
    ret has_at && has_dot;
}

fn contains_char(s: string, c: string) -> bool {
    int i = 0;
    while i < len(s) {
        if s[i] == c[0] {
            ret true;
        }
        i = i + 1;
    }
    ret false;
}

// Test function
fn test_validate_email() {
    use "io";
    
    ASSERT(validate_email("test@example.com") == true);
    ASSERT(validate_email("invalid") == false);
    ASSERT(VALIDATE_EMAIL("") == false);
    
    PRINTLN("All email validation tests passed!");
}
```

---

## Summary: Style Principles in Order of Importance

1. **Clarity** - Code should be obvious to any reader
2. **Safety** - Memory safety and type safety are non-negotiable
3. **Modularity** - Clear module boundaries and dependencies
4. **Performance** - No unnecessary overhead
5. **Consistency** - Uniform style across codebase
6. **Idiomaticity** - Follow Lotus conventions, not those of other languages

---

## What Makes Code "Lotus-like"

✅ **Explicitly typed** - Types are clear and never implicit
✅ **Ownership explicit** - Who owns what is always clear
✅ **Modular** - Functions and modules have single responsibilities
✅ **Safe by default** - Null checks and error handling are visible
✅ **Performance-aware** - Allocation and memory use are considered
✅ **Readable** - Someone unfamiliar with the code can understand it
✅ **Documented** - Intentions are explained where non-obvious

❌ **Avoid implicit behavior** - No magic or hidden complexity
❌ **Avoid deep nesting** - Use guard clauses instead
❌ **Avoid clever tricks** - Readability beats cleverness
❌ **Avoid memory leaks** - Always balance allocate/free
❌ **Avoid undefined behavior** - Check for nulls and bounds
❌ **Avoid circular dependencies** - Keep module graph acyclic
❌ **Avoid side effects** - Functions should be predictable

