package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestIntegrationCompileAndRun tests end-to-end compilation and execution
func TestIntegrationCompileAndRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the compiler first
	tmpDir := t.TempDir()
	compilerPath := filepath.Join(tmpDir, "lotus")

	// Build in the current directory
	cmd := exec.Command("go", "build", "-o", compilerPath, ".")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build compiler: %v\nOutput: %s", err, output)
	}

	tests := []struct {
		name           string
		source         string
		expectedOutput string
		shouldFail     bool
	}{
		{
			name: "simple_return",
			source: `fn int main() {
				ret 42;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "variable_declaration",
			source: `fn int main() {
				int x = 5;
				ret x;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "float32_constant",
			source: `const float32 PI = 3.14;
			fn int main() {
				ret 0;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "float64_constant",
			source: `const float64 E = 2.71828;
			fn int main() {
				ret 0;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "arithmetic_operations",
			source: `fn int main() {
				int a = 10;
				int b = 5;
				int sum = a + b;
				int diff = a - b;
				int prod = a * b;
				int quot = a / b;
				ret quot;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "if_statement",
			source: `fn int main() {
				int x = 10;
				if (x > 5) {
					ret 1;
				}
				ret 0;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "while_loop",
			source: `fn int main() {
				int i = 0;
				while (i < 3) {
					i = i + 1;
				}
				ret i;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "for_loop",
			source: `fn int main() {
				int sum = 0;
				for (int i = 0; i < 5; i = i + 1) {
					sum = sum + i;
				}
				ret sum;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "function_call",
			source: `fn int add(int a, int b) {
				ret a + b;
			}
			fn int main() {
				int result = add(3, 4);
				ret result;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "struct_definition",
			source: `struct Point {
				int x;
				int y;
			}
			fn int main() {
				ret 0;
			}`,
			expectedOutput: "",
			shouldFail:     false, // Struct parsing is now implemented
		},
		{
			name: "enum_definition",
			source: `enum Color {
				RED,
				GREEN,
				BLUE
			}
			fn int main() {
				ret 0;
			}`,
			expectedOutput: "",
			shouldFail:     true, // Enum parsing not yet implemented
		},
		{
			name: "bitwise_operations",
			source: `fn int main() {
				int a = 5;
				int b = 3;
				int and_result = a & b;
				int or_result = a | b;
				int xor_result = a ^ b;
				ret xor_result;
			}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "func_negate",
			source: `use "func";
fn int is_positive(int x) {
    if x > 0 { ret 1; }
    ret 0;
}
fn int main() {
    int a = func::negate(fn is_positive, 5);
    int b = func::negate(fn is_positive, 0);
    if a != 0 { ret 1; }
    if b == 0 { ret 2; }
    ret 0;
}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "func_both",
			source: `use "func";
fn int is_positive(int x) {
    if x > 0 { ret 1; }
    ret 0;
}
fn int is_even(int x) {
    if x % 2 == 0 { ret 1; }
    ret 0;
}
fn int main() {
    int a = func::both(fn is_positive, fn is_even, 4);
    int b = func::both(fn is_positive, fn is_even, 3);
    if a == 0 { ret 1; }
    if b != 0 { ret 2; }
    ret 0;
}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "func_either_fn",
			source: `use "func";
fn int is_positive(int x) {
    if x > 0 { ret 1; }
    ret 0;
}
fn int is_even(int x) {
    if x % 2 == 0 { ret 1; }
    ret 0;
}
fn int main() {
    int a = func::either_fn(fn is_positive, fn is_even, 3);
    int b = func::either_fn(fn is_positive, fn is_even, -3);
    if a == 0 { ret 1; }
    if b != 0 { ret 2; }
    ret 0;
}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "func_on",
			source: `use "func";
fn int add(int a, int b) { ret a + b; }
fn int double(int x) { ret x * 2; }
fn int main() {
    int result = func::on(fn add, fn double, 3, 4);
    if result != 14 { ret 1; }
    ret 0;
}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "func_converge",
			source: `use "func";
fn int add(int a, int b) { ret a + b; }
fn int double(int x) { ret x * 2; }
fn int negate_val(int x) { ret 0 - x; }
fn int main() {
    int result = func::converge(fn add, fn double, fn negate_val, 5);
    if result != 5 { ret 1; }
    ret 0;
}`,
			expectedOutput: "",
			shouldFail:     false,
		},
		{
			name: "func_clamp",
			source: `use "func";
fn int main() {
    int a = func::clamp(0, 10, 5);
    int b = func::clamp(0, 10, -3);
    int c = func::clamp(0, 10, 15);
    if a != 5 { ret 1; }
    if b != 0 { ret 2; }
    if c != 10 { ret 3; }
    ret 0;
}`,
			expectedOutput: "",
			shouldFail:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source file
			sourceFile := filepath.Join(tmpDir, tt.name+".lts")
			err := os.WriteFile(sourceFile, []byte(tt.source), 0644)
			if err != nil {
				t.Fatalf("Failed to write source file: %v", err)
			}

			// Compile
			outputBinary := filepath.Join(tmpDir, tt.name)
			cmd := exec.Command(compilerPath, "-o", outputBinary, sourceFile)
			compileOutput, err := cmd.CombinedOutput()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected compilation to fail, but it succeeded")
				}
				return
			}

			if err != nil {
				t.Fatalf("Compilation failed: %v\nOutput: %s", err, compileOutput)
			}

			// Run if we have expected output
			if tt.expectedOutput != "" {
				cmd = exec.Command(outputBinary)
				runOutput, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("Execution failed: %v\nOutput: %s", err, runOutput)
				}

				actualOutput := strings.TrimSpace(string(runOutput))
				if actualOutput != tt.expectedOutput {
					t.Errorf("Expected output '%s', got '%s'", tt.expectedOutput, actualOutput)
				}
			}
		})
	}
}

// TestIntegrationExampleFiles tests compilation of example files
func TestIntegrationExampleFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the compiler first
	tmpDir := t.TempDir()
	compilerPath := filepath.Join(tmpDir, "lotus")

	cmd := exec.Command("go", "build", "-o", compilerPath, ".")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build compiler: %v\nOutput: %s", err, output)
	}

	// Find example files
	examplesDir := filepath.Join("..", "examples")
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		t.Skip("Examples directory not found")
	}

	files, err := filepath.Glob(filepath.Join(examplesDir, "*.lts"))
	if err != nil {
		t.Fatalf("Failed to list example files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No example files found")
	}

	for _, file := range files {
		fileName := filepath.Base(file)
		t.Run(fileName, func(t *testing.T) {
			outputBinary := filepath.Join(tmpDir, strings.TrimSuffix(fileName, ".lts"))
			cmd := exec.Command(compilerPath, "-o", outputBinary, file)
			compileOutput, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Compilation of %s failed: %v\nOutput: %s", fileName, err, compileOutput)
			}
		})
	}
}

// TestIntegrationTestFiles tests compilation of test files
func TestIntegrationTestFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the compiler first
	tmpDir := t.TempDir()
	compilerPath := filepath.Join(tmpDir, "lotus")

	cmd := exec.Command("go", "build", "-o", compilerPath, ".")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build compiler: %v\nOutput: %s", err, output)
	}

	// Find test files
	testsDir := filepath.Join("..", "tests")
	if _, err := os.Stat(testsDir); os.IsNotExist(err) {
		t.Skip("Tests directory not found")
	}

	files, err := filepath.Glob(filepath.Join(testsDir, "*.lts"))
	if err != nil {
		t.Fatalf("Failed to list test files: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No test files found")
	}

	for _, file := range files {
		fileName := filepath.Base(file)
		t.Run(fileName, func(t *testing.T) {
			outputBinary := filepath.Join(tmpDir, strings.TrimSuffix(fileName, ".lts"))
			cmd := exec.Command(compilerPath, "-o", outputBinary, file)
			compileOutput, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Compilation of %s failed: %v\nOutput: %s", fileName, err, compileOutput)
			}
		})
	}
}

// TestIntegrationTypeSystem tests the type system fixes
func TestIntegrationTypeSystem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the compiler first
	tmpDir := t.TempDir()
	compilerPath := filepath.Join(tmpDir, "lotus")

	cmd := exec.Command("go", "build", "-o", compilerPath, ".")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build compiler: %v\nOutput: %s", err, output)
	}

	// Test that float32 and float64 work in constants and variables
	source := `const float32 PI32 = 3.14;
const float64 PI64 = 3.14159265359;

fn int main() {
	float32 x = 1.5;
	float64 y = 2.71828;

	int result = 0;
	ret result;
}`

	sourceFile := filepath.Join(tmpDir, "float_test.lts")
	err = os.WriteFile(sourceFile, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write source file: %v", err)
	}

	outputBinary := filepath.Join(tmpDir, "float_test")
	cmd = exec.Command(compilerPath, "-o", outputBinary, sourceFile)
	compileOutput, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Compilation failed: %v\nOutput: %s", err, compileOutput)
	}

	t.Log("Float type system test passed successfully")
}
