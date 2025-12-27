package main

import (
	"fmt"
	"strings"
)

// docs.go - Offline documentation for the Lotus compiler
// Provides embedded documentation viewable via the -docs flag.

// Documentation sections
const (
	DocSectionOverview  = "overview"
	DocSectionSyntax    = "syntax"
	DocSectionTypes     = "types"
	DocSectionStdlib    = "stdlib"
	DocSectionExamples  = "examples"
	DocSectionMemory    = "memory"
	DocSectionFunctions = "functions"
)

// AvailableSections lists all documentation sections
var AvailableSections = []string{
	DocSectionOverview,
	DocSectionSyntax,
	DocSectionTypes,
	DocSectionStdlib,
	DocSectionExamples,
	DocSectionMemory,
	DocSectionFunctions,
}

// PrintDocs displays documentation to stdout
func PrintDocs(section string) {
	if section == "" {
		printFullDocs()
		return
	}

	section = strings.ToLower(strings.TrimSpace(section))
	switch section {
	case DocSectionOverview:
		printOverview()
	case DocSectionSyntax:
		printSyntax()
	case DocSectionTypes:
		printTypes()
	case DocSectionStdlib:
		printStdlib()
	case DocSectionExamples:
		printExamples()
	case DocSectionMemory:
		printMemory()
	case DocSectionFunctions:
		printFunctions()
	default:
		fmt.Printf("Unknown documentation section: %s\n\n", section)
		fmt.Println("Available sections:")
		for _, s := range AvailableSections {
			fmt.Printf("  - %s\n", s)
		}
		fmt.Println("\nUsage: lotus -docs-section <section>")
	}
}

func printHeader(title string) {
	line := strings.Repeat("=", len(title)+4)
	fmt.Println(line)
	fmt.Printf("  %s\n", title)
	fmt.Println(line)
	fmt.Println()
}

func printSubHeader(title string) {
	fmt.Printf("── %s ──\n\n", title)
}

// printRaw prints a raw string without format specifier warnings.
// This is used to print code examples that contain format directives like %d.
func printRaw(s string) {
	fmt.Print(s)
}

func printFullDocs() {
	printHeader("LOTUS LANGUAGE DOCUMENTATION")

	fmt.Println("Lotus is a systems programming language combining Rust's safety,")
	fmt.Println("C++'s performance, and Go's simplicity.")
	fmt.Println()
	fmt.Println("Documentation Sections:")
	fmt.Println("  overview   - Language overview and philosophy")
	fmt.Println("  syntax     - Language syntax reference")
	fmt.Println("  types      - Type system and declarations")
	fmt.Println("  stdlib     - Standard library reference")
	fmt.Println("  examples   - Code examples")
	fmt.Println("  memory     - Memory management")
	fmt.Println("  functions  - Function declarations")
	fmt.Println()
	fmt.Println("View a specific section:")
	fmt.Println("  lotus -docs-section <section>")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  lotus -docs-section stdlib")
	fmt.Println()

	// Print quick reference
	printSubHeader("Quick Reference")
	fmt.Println("  File extension:    .lts")
	fmt.Println("  Entry point:       fn int main() { ... }")
	fmt.Println("  Comments:          // single line")
	fmt.Println("  Import:            use \"module\";")
	fmt.Println("  Variable:          type name = value;")
	fmt.Println("  Function:          fn type name(params) { ... }")
	fmt.Println()
	fmt.Printf("Lotus Compiler Version %s\n", Version)
}

func printOverview() {
	printHeader("LOTUS LANGUAGE OVERVIEW")

	fmt.Print(`Lotus is a systems programming language designed to combine:

  • Rust's memory safety philosophy
  • C++'s performance and control  
  • Go's simplicity and clarity

CORE PRINCIPLES

  1. Explicit Over Implicit
     - Make intentions clear in code
     - Avoid magic or hidden behavior
     - Use explicit keywords and types

  2. Safety By Default
     - Memory safety through static analysis
     - Type safety enforced at compile time
     - Clear ownership and lifetime semantics

  3. Performance Matters
     - No unnecessary abstractions
     - Direct hardware access when needed
     - Zero-cost abstractions

  4. Clarity Over Cleverness
     - Code is read more often than written
     - Simple, straightforward syntax
     - Minimize cognitive load

FILE ORGANIZATION

  • Source files use .lts extension
  • Use snake_case for file names: string_utils.lts
  • Entry point is fn main() { ... }
`)
}

func printSyntax() {
	printHeader("LOTUS SYNTAX REFERENCE")

	fmt.Print(`COMMENTS
  // This is a single-line comment

IMPORTS
  use "io";                    // Import entire module
  use "io::printf";            // Import specific function
  use "io::*";                 // Wildcard import
  use "io" as output;          // Aliased import

VARIABLE DECLARATIONS
  int x = 42;                  // Typed variable
  string name = "Lotus";       // String variable
  bool flag = true;            // Boolean variable
  int arr[5];                  // Array declaration

CONTROL FLOW

  If Statement:
    if condition {
        // code
    } elif other_condition {
        // code
    } else {
        // code
    }

  While Loop:
    while condition {
        // code
    }

  For Loop:
    for int i = 0; i < 10; i = i + 1 {
        // code
    }

OPERATORS

  Arithmetic:   +  -  *  /  %
  Comparison:   ==  !=  <  >  <=  >=
  Logical:      &&  ||  !
  Bitwise:      &  |  ^  ~  <<  >>
  Assignment:   =

KEYWORDS
  fn, ret, if, elif, else, while, for, 
  use, as, true, false, nil
`)
}

func printTypes() {
	printHeader("LOTUS TYPE SYSTEM")

	fmt.Print(`PRIMITIVE TYPES

  int       - Signed integer (platform-dependent size)
  int8      - 8-bit signed integer
  int16     - 16-bit signed integer
  int32     - 32-bit signed integer
  int64     - 64-bit signed integer
  
  uint      - Unsigned integer
  uint8     - 8-bit unsigned integer (byte)
  uint16    - 16-bit unsigned integer
  uint32    - 32-bit unsigned integer
  uint64    - 64-bit unsigned integer
  
  float32   - 32-bit floating point
  float64   - 64-bit floating point
  
  bool      - Boolean (true/false)
  string    - String type
  char      - Single character

ARRAY TYPES
  int arr[10];                 // Fixed-size array
  int matrix[3][3];            // Multi-dimensional array

POINTER TYPES
  *int ptr;                    // Pointer to int
  *void data;                  // Generic pointer

TYPE DECLARATION SYNTAX
  type variable_name = initial_value;

EXAMPLES
  int count = 0;
  float64 pi = 3.14159;
  bool active = true;
  string message = "Hello";
  int numbers[5] = {1, 2, 3, 4, 5};
`)
}

func printStdlib() {
	printHeader("LOTUS STANDARD LIBRARY")

	printRaw(`IMPORT SYNTAX
  use "module";                // Import all functions
  use "module::function";      // Import specific function

AVAILABLE MODULES

── io (Input/Output) ──

  print(args...)              Output to stdout
  println(args...)            Output with newline
  printf(format, args...)     Formatted output
  fprintf(fd, format, ...)    Output to file descriptor
  sprint(args...)             Format to string
  sprintf(format, args...)    Format to string
  sprintln(args...)           Format to string with newline

  Format specifiers:
    %d  - Integer (decimal)
    %b  - Integer (binary)
    %o  - Integer (octal)
    %x  - Integer (hex lowercase)
    %X  - Integer (hex uppercase)
    %c  - Character (byte)
    %s  - String
    %q  - Quoted string
    %v  - Default format
    %%  - Literal percent

── mem (Memory Management) ──

  malloc(size)                Allocate memory
  free(ptr)                   Deallocate memory
  sizeof(type)                Get size of type in bytes

── math (Mathematical Functions) ──

  abs(n)                      Absolute value
  min(a, b)                   Minimum of two values
  max(a, b)                   Maximum of two values
  sqrt(x)                     Square root
  pow(base, exp)              Power (base^exp)
  floor(x)                    Largest integer ≤ x
  ceil(x)                     Smallest integer ≥ x
  round(x)                    Nearest integer
  gcd(a, b)                   Greatest common divisor
  lcm(a, b)                   Least common multiple

── str (String Functions) ──

  len(s)                      String length
  concat(a, b)                Concatenate strings
  substr(s, start, len)       Extract substring
  index(s, substr)            Find substring index
  upper(s)                    Convert to uppercase
  lower(s)                    Convert to lowercase
  trim(s)                     Remove whitespace
  split(s, delim)             Split string
  replace(s, old, new)        Replace occurrences

── file (File Operations) ──

  open(path, mode)            Open file
  close(fd)                   Close file
  read(fd, buf, size)         Read from file
  write(fd, buf, size)        Write to file
  seek(fd, offset, whence)    Seek in file
`)
}

func printExamples() {
	printHeader("LOTUS CODE EXAMPLES")

	printRaw(`HELLO WORLD
  use "io";

  fn main() {
      println("Hello, World!");
  }

VARIABLES AND ARITHMETIC
  use "io";

  fn main() {
      a: int = 10;
      b: int = 20;
      sum: int = a + b;
      println("Sum:", sum);
  }

CONTROL FLOW
  use "io";

  fn main() {
      x: int = 5;
      
      if x > 10 {
          println("Greater than 10");
      } elif x > 0 {
          println("Positive");
      } else {
          println("Zero or negative");
      }
  }

LOOPS
  use "io";

  fn main() {
      // While loop
      i: int = 0;
      while i < 5 {
          println(i);
          i = i + 1;
      }

      // For loop
      for j: int = 0; j < 5; j = j + 1 {
          println(j);
      }
  }

FUNCTIONS
  use "io";

  fn int add(int a, int b) {
      ret a + b;
  }

  fn main() {
      result: int = add(3, 4);
      printf("Result: %d\n", result);
  }

MEMORY ALLOCATION
  use "io";
  use "mem";

  fn main() {
      size: int = sizeof(int);
      ptr: *int = malloc(size);
      
      // Use memory...
      
      free(ptr);
      println("Memory freed");
  }

ARRAYS
  use "io";

  fn main() {
      numbers: int[5] = {1, 2, 3, 4, 5};
      
      for i: int = 0; i < 5; i = i + 1 {
          printf("numbers[%d] = %d\n", i, numbers[i]);
      }
  }

FORMATTED OUTPUT
  use "io";

  fn main() {
      name: string = "Lotus";
      version: int = 1;
      
      printf("Language: %s\n", name);
      printf("Version: %d\n", version);
      printf("Hex: %x, Binary: %b\n", 255, 255);
  }
`)
}

func printMemory() {
	printHeader("LOTUS MEMORY MANAGEMENT")

	printRaw(`OVERVIEW

Lotus provides manual memory management with safety-focused design.
The mem module provides core allocation functions.

IMPORTING MEMORY FUNCTIONS
  use "mem";

ALLOCATION

  malloc(size) - Allocate 'size' bytes
    ptr: *int = malloc(sizeof(int));
    
  sizeof(type) - Get size of type in bytes
    size: int = sizeof(int64);  // Returns 8

DEALLOCATION

  free(ptr) - Free previously allocated memory
    free(ptr);
    
  Important: Always free allocated memory to prevent leaks.

EXAMPLE: Dynamic Array
  use "io";
  use "mem";

  fn main() {
      count: int = 10;
      arr: *int = malloc(count * sizeof(int));
      
      // Initialize
      for i: int = 0; i < count; i = i + 1 {
          arr[i] = i * 2;
      }
      
      // Use array
      for i: int = 0; i < count; i = i + 1 {
          printf("arr[%d] = %d\n", i, arr[i]);
      }
      
      // Clean up
      free(arr);
  }

BEST PRACTICES

  1. Always pair malloc with free
  2. Set pointers to nil after freeing
  3. Check allocation success (non-nil return)
  4. Use sizeof() for type-safe allocation
  5. Avoid double-free errors
`)
}

func printFunctions() {
	printHeader("LOTUS FUNCTIONS")

	fmt.Print(`FUNCTION SYNTAX

  fn return_type function_name(param_type param_name, ...) {
      // function body
      ret value;
  }

EXAMPLES

No Return Value (void):
  fn void greet() {
      println("Hello!");
  }

With Return Value:
  fn int square(int x) {
      ret x * x;
  }

Multiple Parameters:
  fn int add(int a, int b) {
      ret a + b;
  }

Calling Functions:
  result: int = square(5);
  sum: int = add(10, 20);

MAIN FUNCTION

Every Lotus program requires a main function as entry point:

  fn main() {
      // program starts here
  }

RETURN STATEMENT

Use 'ret' to return a value:
  fn int double(int x) {
      ret x * 2;
  }

Early return is allowed:
  fn int abs(int x) {
      if x < 0 {
          ret -x;
      }
      ret x;
  }

FUNCTION NAMING CONVENTIONS

  • Use snake_case for function names
  • Choose descriptive, verb-based names
  • Keep names concise but clear

Good:
  fn calculate_area(int width, int height)
  fn parse_input(string data)
  fn is_valid(int value)

Avoid:
  fn calc(int w, int h)
  fn doStuff(string s)
  fn process(int x)
`)
}
