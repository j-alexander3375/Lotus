package main

import (
	"fmt"
	"math"
	"strings"
)

// interpreter.go - Tree-walking interpreter for the Lotus language.
// Evaluates AST nodes directly without compiling to native code.

// ============================================================================
// Runtime Values
// ============================================================================

type ValueKind int

const (
	KindInt ValueKind = iota
	KindFloat
	KindBool
	KindString
	KindChar
	KindNull
	KindArray
	KindStruct
	KindFunction
	KindVoid
)

// Value is the universal runtime value type.
type Value struct {
	Kind    ValueKind
	Int     int
	Float   float64
	Bool    bool
	Str     string          // string value or struct type name
	Char    rune
	Arr     []*Value
	Fields  map[string]*Value
	Fn      *InterpFn
	builtin *builtinFn
}

// InterpFn holds a user-defined or lambda function with its captured closure.
type InterpFn struct {
	Params []FunctionParam
	Body   []ASTNode
	Expr   ASTNode // for single-expression lambdas
	Env    *Environment
	Name   string
}

// builtinFn wraps a native Go function as a Lotus callable.
type builtinFn struct {
	name string
	fn   func([]*Value) (*Value, error)
}

// TypeName returns a human-readable type string.
func (v *Value) TypeName() string {
	switch v.Kind {
	case KindInt:
		return "int"
	case KindFloat:
		return "float"
	case KindBool:
		return "bool"
	case KindString:
		return "str"
	case KindChar:
		return "char"
	case KindNull:
		return "null"
	case KindArray:
		return "array"
	case KindStruct:
		if v.Str != "" {
			return v.Str
		}
		return "struct"
	case KindFunction:
		if v.Fn != nil && v.Fn.Name != "" {
			return "fn(" + v.Fn.Name + ")"
		}
		if v.builtin != nil {
			return "fn(" + v.builtin.name + ")"
		}
		return "fn"
	case KindVoid:
		return "()"
	default:
		return "?"
	}
}

// Display returns a quoted representation for REPL output.
func (v *Value) Display() string {
	switch v.Kind {
	case KindInt:
		return fmt.Sprintf("%d", v.Int)
	case KindFloat:
		s := fmt.Sprintf("%g", v.Float)
		if !strings.ContainsAny(s, ".e") {
			s += ".0"
		}
		return s
	case KindBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case KindString:
		return fmt.Sprintf("%q", v.Str)
	case KindChar:
		return fmt.Sprintf("'%c'", v.Char)
	case KindNull:
		return "null"
	case KindArray:
		parts := make([]string, len(v.Arr))
		for i, elem := range v.Arr {
			parts[i] = elem.Display()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case KindStruct:
		parts := []string{}
		for k, fv := range v.Fields {
			parts = append(parts, k+": "+fv.Display())
		}
		return "{ " + strings.Join(parts, ", ") + " }"
	case KindFunction:
		if v.builtin != nil {
			return "<builtin " + v.builtin.name + ">"
		}
		if v.Fn != nil && v.Fn.Name != "" {
			return "<fn " + v.Fn.Name + ">"
		}
		return "<lambda>"
	case KindVoid:
		return "()"
	}
	return "?"
}

// ============================================================================
// Control-flow signals (carried as errors)
// ============================================================================

type returnSignal struct{ val *Value }
type breakSignal struct{}
type continueSignal struct{}

func (r returnSignal) Error() string   { return "return" }
func (b breakSignal) Error() string    { return "break" }
func (c continueSignal) Error() string { return "continue" }

// ============================================================================
// Environment (lexical scope chain)
// ============================================================================

type Environment struct {
	vars   map[string]*Value
	consts map[string]struct{}
	outer  *Environment
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{
		vars:   make(map[string]*Value),
		consts: make(map[string]struct{}),
		outer:  outer,
	}
}

func (e *Environment) Get(name string) (*Value, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}

// Set assigns to an existing binding anywhere in the scope chain.
// Falls back to defining in the current scope if not found.
func (e *Environment) Set(name string, val *Value) error {
	scope := e
	for scope != nil {
		if _, ok := scope.vars[name]; ok {
			if _, isConst := scope.consts[name]; isConst {
				return fmt.Errorf("cannot assign to constant %q", name)
			}
			scope.vars[name] = val
			return nil
		}
		scope = scope.outer
	}
	e.vars[name] = val
	return nil
}

// Define introduces a new binding in the current scope.
func (e *Environment) Define(name string, val *Value, isConst bool) {
	e.vars[name] = val
	if isConst {
		e.consts[name] = struct{}{}
	}
}

// Snapshot returns a flat view of all visible bindings (inner shadows outer).
func (e *Environment) Snapshot() map[string]*Value {
	result := make(map[string]*Value)
	var collect func(env *Environment)
	collect = func(env *Environment) {
		if env.outer != nil {
			collect(env.outer)
		}
		for k, v := range env.vars {
			result[k] = v
		}
	}
	collect(e)
	return result
}

// ============================================================================
// Interpreter
// ============================================================================

type Interpreter struct {
	globals *Environment
}

func NewInterpreter() *Interpreter {
	interp := &Interpreter{globals: NewEnvironment(nil)}
	interp.registerBuiltins()
	return interp
}

// EvalInGlobals evaluates a node in the global scope and returns the result.
func (interp *Interpreter) EvalInGlobals(node ASTNode) (*Value, error) {
	return interp.evalNode(node, interp.globals)
}

// EvalStatements evaluates a slice of statements and returns the last non-void value.
func (interp *Interpreter) EvalStatements(stmts []ASTNode, env *Environment) (*Value, error) {
	last := &Value{Kind: KindVoid}
	for _, stmt := range stmts {
		v, err := interp.evalNode(stmt, env)
		if err != nil {
			return nil, err
		}
		if v != nil && v.Kind != KindVoid {
			last = v
		}
	}
	return last, nil
}

// evalNode is the main dispatch function.
func (interp *Interpreter) evalNode(node ASTNode, env *Environment) (*Value, error) {
	if node == nil {
		return &Value{Kind: KindVoid}, nil
	}

	switch n := node.(type) {

	// ── Literals ──────────────────────────────────────────────────────────────
	case *IntLiteral:
		return &Value{Kind: KindInt, Int: n.Value}, nil

	case *FloatLiteral:
		return &Value{Kind: KindFloat, Float: float64(n.Value) / 1000.0}, nil

	case *BoolLiteral:
		return &Value{Kind: KindBool, Bool: n.Value}, nil

	case *StringLiteral:
		return &Value{Kind: KindString, Str: n.Value}, nil

	case *CharLiteral:
		runes := []rune(n.Value)
		if len(runes) > 0 {
			return &Value{Kind: KindChar, Char: runes[0]}, nil
		}
		return &Value{Kind: KindChar}, nil

	case *NullLiteral:
		return &Value{Kind: KindNull}, nil

	case *InterpolatedString:
		return interp.evalInterpolated(n, env)

	// ── Variables / constants ─────────────────────────────────────────────────
	case *Identifier:
		v, ok := env.Get(n.Name)
		if !ok {
			return nil, fmt.Errorf("undefined: %q", n.Name)
		}
		return v, nil

	case *VariableDeclaration:
		var val *Value
		var err error
		if n.Value != nil {
			val, err = interp.evalNode(n.Value, env)
			if err != nil {
				return nil, err
			}
		} else {
			val = &Value{Kind: KindNull}
		}
		env.Define(n.Name, val, false)
		return &Value{Kind: KindVoid}, nil

	case *ConstantDeclaration:
		val, err := interp.evalNode(n.Value, env)
		if err != nil {
			return nil, err
		}
		env.Define(n.Name, val, true)
		return &Value{Kind: KindVoid}, nil

	case *Assignment:
		val, err := interp.evalNode(n.Value, env)
		if err != nil {
			return nil, err
		}
		return interp.assign(n.Target, val, env)

	case *CompoundAssignment:
		return interp.evalCompoundAssign(n, env)

	// ── Arithmetic / bitwise ──────────────────────────────────────────────────
	case *BinaryOp:
		return interp.evalBinaryOp(n, env)

	case *UnaryOp:
		return interp.evalUnaryOp(n, env)

	case *IncrementOp:
		return interp.evalIncrementOp(n, env)

	case *BitwiseOp:
		return interp.evalBitwiseOp(n, env)

	// ── Comparisons / logic ───────────────────────────────────────────────────
	case *Comparison:
		return interp.evalComparison(n, env)

	case *LogicalOp:
		return interp.evalLogicalOp(n, env)

	case *TernaryOp:
		cond, err := interp.evalNode(n.Condition, env)
		if err != nil {
			return nil, err
		}
		if isTruthy(cond) {
			return interp.evalNode(n.TrueExpr, env)
		}
		return interp.evalNode(n.FalseExpr, env)

	// ── Control flow ──────────────────────────────────────────────────────────
	case *IfStatement:
		return interp.evalIf(n, env)

	case *WhileLoop:
		return interp.evalWhile(n, env)

	case *ForLoop:
		return interp.evalFor(n, env)

	case *BreakStatement:
		return nil, breakSignal{}

	case *ContinueStatement:
		return nil, continueSignal{}

	case *ReturnStatement:
		if n.Value == nil {
			return nil, returnSignal{val: &Value{Kind: KindVoid}}
		}
		val, err := interp.evalNode(n.Value, env)
		if err != nil {
			return nil, err
		}
		return nil, returnSignal{val: val}

	// ── Functions ─────────────────────────────────────────────────────────────
	case *FunctionDefinition:
		fn := &InterpFn{Params: n.Parameters, Body: n.Body, Env: env, Name: n.Name}
		val := &Value{Kind: KindFunction, Fn: fn}
		env.Define(n.Name, val, false)
		return &Value{Kind: KindVoid}, nil

	case *LambdaExpression:
		fn := &InterpFn{Params: n.Parameters, Body: n.Body, Expr: n.Expression, Env: env}
		return &Value{Kind: KindFunction, Fn: fn}, nil

	case *FunctionReference:
		v, ok := env.Get(n.FunctionName)
		if !ok {
			return nil, fmt.Errorf("undefined function %q", n.FunctionName)
		}
		return v, nil

	case *FunctionCall:
		return interp.evalCall(n, env)

	case *PipeExpression:
		return interp.evalPipe(n, env)

	case *DecoratedFunction:
		// Apply decorators then define the underlying function
		return interp.evalNode(n.Function, env)

	// ── Arrays ────────────────────────────────────────────────────────────────
	case *ArrayLiteral:
		arr := make([]*Value, len(n.Elements))
		for i, elem := range n.Elements {
			v, err := interp.evalNode(elem, env)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		return &Value{Kind: KindArray, Arr: arr}, nil

	case *ArrayAccess:
		return interp.evalArrayAccess(n, env)

	case *ArrayDeclaration:
		return interp.evalArrayDecl(n, env)

	// ── Structs ───────────────────────────────────────────────────────────────
	case *StructDefinition:
		return interp.defineStruct(n, env)

	case *StructLiteral:
		return interp.evalStructLiteral(n, env)

	case *FieldAccess:
		return interp.evalFieldAccess(n, env)

	// ── Match ─────────────────────────────────────────────────────────────────
	case *MatchExpression:
		return interp.evalMatch(n, env)

	// ── Internal helper node ──────────────────────────────────────────────────
	case *valueNode:
		return n.v, nil

	// ── Ignored nodes ─────────────────────────────────────────────────────────
	case *ImportStatement:
		return &Value{Kind: KindVoid}, nil

	default:
		return &Value{Kind: KindVoid}, nil
	}
}

// ============================================================================
// Builtins
// ============================================================================

func (interp *Interpreter) registerBuiltins() {
	define := func(name string, fn func([]*Value) (*Value, error)) {
		bfn := &builtinFn{name: name, fn: fn}
		interp.globals.Define(name, &Value{Kind: KindFunction, builtin: bfn}, false)
	}

	define("print", func(args []*Value) (*Value, error) {
		parts := make([]string, len(args))
		for i, a := range args {
			parts[i] = valueStr(a)
		}
		fmt.Print(strings.Join(parts, " "))
		return &Value{Kind: KindVoid}, nil
	})

	define("println", func(args []*Value) (*Value, error) {
		parts := make([]string, len(args))
		for i, a := range args {
			parts[i] = valueStr(a)
		}
		fmt.Println(strings.Join(parts, " "))
		return &Value{Kind: KindVoid}, nil
	})

	define("len", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("len() expects 1 argument")
		}
		switch args[0].Kind {
		case KindString:
			return &Value{Kind: KindInt, Int: len([]rune(args[0].Str))}, nil
		case KindArray:
			return &Value{Kind: KindInt, Int: len(args[0].Arr)}, nil
		}
		return nil, fmt.Errorf("len() unsupported for %s", args[0].TypeName())
	})

	define("abs", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("abs() expects 1 argument")
		}
		if args[0].Kind == KindInt {
			v := args[0].Int
			if v < 0 {
				v = -v
			}
			return &Value{Kind: KindInt, Int: v}, nil
		}
		return &Value{Kind: KindFloat, Float: math.Abs(toFloat(args[0]))}, nil
	})

	define("sqrt", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("sqrt() expects 1 argument")
		}
		return &Value{Kind: KindFloat, Float: math.Sqrt(toFloat(args[0]))}, nil
	})

	define("pow", func(args []*Value) (*Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("pow() expects 2 arguments")
		}
		return &Value{Kind: KindFloat, Float: math.Pow(toFloat(args[0]), toFloat(args[1]))}, nil
	})

	define("floor", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("floor() expects 1 argument")
		}
		return &Value{Kind: KindInt, Int: int(math.Floor(toFloat(args[0])))}, nil
	})

	define("ceil", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("ceil() expects 1 argument")
		}
		return &Value{Kind: KindInt, Int: int(math.Ceil(toFloat(args[0])))}, nil
	})

	define("round", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("round() expects 1 argument")
		}
		return &Value{Kind: KindInt, Int: int(math.Round(toFloat(args[0])))}, nil
	})

	define("min", func(args []*Value) (*Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("min() expects 2 arguments")
		}
		if toFloat(args[0]) <= toFloat(args[1]) {
			return args[0], nil
		}
		return args[1], nil
	})

	define("max", func(args []*Value) (*Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("max() expects 2 arguments")
		}
		if toFloat(args[0]) >= toFloat(args[1]) {
			return args[0], nil
		}
		return args[1], nil
	})

	define("sin", mathFn1(math.Sin))
	define("cos", mathFn1(math.Cos))
	define("tan", mathFn1(math.Tan))
	define("asin", mathFn1(math.Asin))
	define("acos", mathFn1(math.Acos))
	define("atan2", func(args []*Value) (*Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("atan2() expects 2 arguments")
		}
		return &Value{Kind: KindFloat, Float: math.Atan2(toFloat(args[0]), toFloat(args[1]))}, nil
	})

	define("toInt", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toInt() expects 1 argument")
		}
		switch args[0].Kind {
		case KindInt:
			return args[0], nil
		case KindFloat:
			return &Value{Kind: KindInt, Int: int(args[0].Float)}, nil
		case KindBool:
			if args[0].Bool {
				return &Value{Kind: KindInt, Int: 1}, nil
			}
			return &Value{Kind: KindInt, Int: 0}, nil
		case KindString:
			var i int
			if _, err := fmt.Sscanf(args[0].Str, "%d", &i); err != nil {
				return nil, fmt.Errorf("cannot convert %q to int", args[0].Str)
			}
			return &Value{Kind: KindInt, Int: i}, nil
		}
		return nil, fmt.Errorf("toInt() unsupported for %s", args[0].TypeName())
	})

	define("toFloat", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toFloat() expects 1 argument")
		}
		return &Value{Kind: KindFloat, Float: toFloat(args[0])}, nil
	})

	define("toBool", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toBool() expects 1 argument")
		}
		return &Value{Kind: KindBool, Bool: isTruthy(args[0])}, nil
	})

	define("toString", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toString() expects 1 argument")
		}
		return &Value{Kind: KindString, Str: valueStr(args[0])}, nil
	})

	define("concat", func(args []*Value) (*Value, error) {
		var sb strings.Builder
		for _, a := range args {
			sb.WriteString(valueStr(a))
		}
		return &Value{Kind: KindString, Str: sb.String()}, nil
	})

	define("contains", func(args []*Value) (*Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("contains() expects 2 arguments")
		}
		if args[0].Kind == KindString && args[1].Kind == KindString {
			return &Value{Kind: KindBool, Bool: strings.Contains(args[0].Str, args[1].Str)}, nil
		}
		if args[0].Kind == KindArray {
			for _, elem := range args[0].Arr {
				if valuesEqual(elem, args[1]) {
					return &Value{Kind: KindBool, Bool: true}, nil
				}
			}
			return &Value{Kind: KindBool, Bool: false}, nil
		}
		return nil, fmt.Errorf("contains() unsupported for %s", args[0].TypeName())
	})

	define("toUpper", func(args []*Value) (*Value, error) {
		if len(args) != 1 || args[0].Kind != KindString {
			return nil, fmt.Errorf("toUpper() expects a string")
		}
		return &Value{Kind: KindString, Str: strings.ToUpper(args[0].Str)}, nil
	})

	define("toLower", func(args []*Value) (*Value, error) {
		if len(args) != 1 || args[0].Kind != KindString {
			return nil, fmt.Errorf("toLower() expects a string")
		}
		return &Value{Kind: KindString, Str: strings.ToLower(args[0].Str)}, nil
	})

	define("trim", func(args []*Value) (*Value, error) {
		if len(args) != 1 || args[0].Kind != KindString {
			return nil, fmt.Errorf("trim() expects a string")
		}
		return &Value{Kind: KindString, Str: strings.TrimSpace(args[0].Str)}, nil
	})

	define("startsWith", func(args []*Value) (*Value, error) {
		if len(args) != 2 || args[0].Kind != KindString || args[1].Kind != KindString {
			return nil, fmt.Errorf("startsWith() expects two strings")
		}
		return &Value{Kind: KindBool, Bool: strings.HasPrefix(args[0].Str, args[1].Str)}, nil
	})

	define("endsWith", func(args []*Value) (*Value, error) {
		if len(args) != 2 || args[0].Kind != KindString || args[1].Kind != KindString {
			return nil, fmt.Errorf("endsWith() expects two strings")
		}
		return &Value{Kind: KindBool, Bool: strings.HasSuffix(args[0].Str, args[1].Str)}, nil
	})

	define("substring", func(args []*Value) (*Value, error) {
		if len(args) < 2 || args[0].Kind != KindString {
			return nil, fmt.Errorf("substring() expects (str, start[, end])")
		}
		runes := []rune(args[0].Str)
		start := args[1].Int
		end := len(runes)
		if len(args) >= 3 {
			end = args[2].Int
		}
		if start < 0 || end > len(runes) || start > end {
			return nil, fmt.Errorf("substring() index out of range")
		}
		return &Value{Kind: KindString, Str: string(runes[start:end])}, nil
	})

	define("push", func(args []*Value) (*Value, error) {
		if len(args) < 2 || args[0].Kind != KindArray {
			return nil, fmt.Errorf("push() expects (array, value)")
		}
		args[0].Arr = append(args[0].Arr, args[1])
		return &Value{Kind: KindVoid}, nil
	})

	define("pop", func(args []*Value) (*Value, error) {
		if len(args) != 1 || args[0].Kind != KindArray {
			return nil, fmt.Errorf("pop() expects an array")
		}
		n := len(args[0].Arr)
		if n == 0 {
			return nil, fmt.Errorf("pop() on empty array")
		}
		val := args[0].Arr[n-1]
		args[0].Arr = args[0].Arr[:n-1]
		return val, nil
	})

	define("assert", func(args []*Value) (*Value, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("assert() requires at least 1 argument")
		}
		if !isTruthy(args[0]) {
			msg := "assertion failed"
			if len(args) > 1 {
				msg = valueStr(args[1])
			}
			return nil, fmt.Errorf("%s", msg)
		}
		return &Value{Kind: KindVoid}, nil
	})

	define("typeof", func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("typeof() expects 1 argument")
		}
		return &Value{Kind: KindString, Str: args[0].TypeName()}, nil
	})

	define("exit", func(args []*Value) (*Value, error) {
		code := 0
		if len(args) > 0 && args[0].Kind == KindInt {
			code = args[0].Int
		}
		return nil, fmt.Errorf("exit(%d)", code)
	})
}

// ============================================================================
// Eval helpers
// ============================================================================

func (interp *Interpreter) evalInterpolated(n *InterpolatedString, env *Environment) (*Value, error) {
	var sb strings.Builder
	for _, part := range n.Parts {
		if !part.IsExpr {
			sb.WriteString(part.Text)
		} else {
			v, err := interp.evalNode(part.Expr, env)
			if err != nil {
				return nil, err
			}
			sb.WriteString(valueStr(v))
		}
	}
	return &Value{Kind: KindString, Str: sb.String()}, nil
}

func (interp *Interpreter) evalBinaryOp(n *BinaryOp, env *Environment) (*Value, error) {
	left, err := interp.evalNode(n.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := interp.evalNode(n.Right, env)
	if err != nil {
		return nil, err
	}

	// String concatenation with +
	if n.Operator == TokenPlus && (left.Kind == KindString || right.Kind == KindString) {
		return &Value{Kind: KindString, Str: valueStr(left) + valueStr(right)}, nil
	}

	if isNumeric(left) && isNumeric(right) {
		if left.Kind == KindFloat || right.Kind == KindFloat {
			a, b := toFloat(left), toFloat(right)
			switch n.Operator {
			case TokenPlus:
				return &Value{Kind: KindFloat, Float: a + b}, nil
			case TokenMinus:
				return &Value{Kind: KindFloat, Float: a - b}, nil
			case TokenStar:
				return &Value{Kind: KindFloat, Float: a * b}, nil
			case TokenSlash:
				if b == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return &Value{Kind: KindFloat, Float: a / b}, nil
			case TokenPercent:
				return &Value{Kind: KindFloat, Float: math.Mod(a, b)}, nil
			}
		} else {
			a, b := left.Int, right.Int
			switch n.Operator {
			case TokenPlus:
				return &Value{Kind: KindInt, Int: a + b}, nil
			case TokenMinus:
				return &Value{Kind: KindInt, Int: a - b}, nil
			case TokenStar:
				return &Value{Kind: KindInt, Int: a * b}, nil
			case TokenSlash:
				if b == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return &Value{Kind: KindInt, Int: a / b}, nil
			case TokenPercent:
				if b == 0 {
					return nil, fmt.Errorf("modulo by zero")
				}
				return &Value{Kind: KindInt, Int: a % b}, nil
			}
		}
	}

	return nil, fmt.Errorf("unsupported binary op on %s and %s", left.TypeName(), right.TypeName())
}

func (interp *Interpreter) evalUnaryOp(n *UnaryOp, env *Environment) (*Value, error) {
	val, err := interp.evalNode(n.Operand, env)
	if err != nil {
		return nil, err
	}
	switch n.Operator {
	case TokenMinus:
		if val.Kind == KindInt {
			return &Value{Kind: KindInt, Int: -val.Int}, nil
		}
		if val.Kind == KindFloat {
			return &Value{Kind: KindFloat, Float: -val.Float}, nil
		}
	case TokenExclaim:
		return &Value{Kind: KindBool, Bool: !isTruthy(val)}, nil
	case TokenTilde:
		if val.Kind == KindInt {
			return &Value{Kind: KindInt, Int: ^val.Int}, nil
		}
	}
	return nil, fmt.Errorf("unsupported unary operator on %s", val.TypeName())
}

func (interp *Interpreter) evalIncrementOp(n *IncrementOp, env *Environment) (*Value, error) {
	ident, ok := n.Operand.(*Identifier)
	if !ok {
		return nil, fmt.Errorf("increment/decrement requires an identifier")
	}
	val, exists := env.Get(ident.Name)
	if !exists {
		return nil, fmt.Errorf("undefined: %q", ident.Name)
	}
	var newVal *Value
	switch val.Kind {
	case KindInt:
		delta := 1
		if n.Operator == TokenMinusMinus {
			delta = -1
		}
		newVal = &Value{Kind: KindInt, Int: val.Int + delta}
	case KindFloat:
		delta := 1.0
		if n.Operator == TokenMinusMinus {
			delta = -1.0
		}
		newVal = &Value{Kind: KindFloat, Float: val.Float + delta}
	default:
		return nil, fmt.Errorf("increment/decrement unsupported for %s", val.TypeName())
	}
	env.Set(ident.Name, newVal)
	if n.IsPrefix {
		return newVal, nil
	}
	return val, nil // postfix returns old value
}

func (interp *Interpreter) evalBitwiseOp(n *BitwiseOp, env *Environment) (*Value, error) {
	left, err := interp.evalNode(n.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := interp.evalNode(n.Right, env)
	if err != nil {
		return nil, err
	}
	if left.Kind != KindInt || right.Kind != KindInt {
		return nil, fmt.Errorf("bitwise ops require int operands")
	}
	a, b := left.Int, right.Int
	switch n.Operator {
	case TokenAmpersand:
		return &Value{Kind: KindInt, Int: a & b}, nil
	case TokenPipe:
		return &Value{Kind: KindInt, Int: a | b}, nil
	case TokenCaret:
		return &Value{Kind: KindInt, Int: a ^ b}, nil
	case TokenLShift:
		return &Value{Kind: KindInt, Int: a << uint(b)}, nil
	case TokenRShift:
		return &Value{Kind: KindInt, Int: a >> uint(b)}, nil
	}
	return nil, fmt.Errorf("unknown bitwise operator")
}

func (interp *Interpreter) evalComparison(n *Comparison, env *Environment) (*Value, error) {
	left, err := interp.evalNode(n.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := interp.evalNode(n.Right, env)
	if err != nil {
		return nil, err
	}
	var result bool
	switch n.Operator {
	case TokenEqual:
		result = valuesEqual(left, right)
	case TokenNotEqual:
		result = !valuesEqual(left, right)
	case TokenLess:
		result = compareValues(left, right, "<")
	case TokenGreater:
		result = compareValues(left, right, ">")
	case TokenLessEq:
		result = compareValues(left, right, "<=")
	case TokenGreaterEq:
		result = compareValues(left, right, ">=")
	default:
		return nil, fmt.Errorf("unknown comparison operator")
	}
	return &Value{Kind: KindBool, Bool: result}, nil
}

func (interp *Interpreter) evalLogicalOp(n *LogicalOp, env *Environment) (*Value, error) {
	left, err := interp.evalNode(n.Left, env)
	if err != nil {
		return nil, err
	}
	switch n.Operator {
	case TokenAnd:
		if !isTruthy(left) {
			return &Value{Kind: KindBool, Bool: false}, nil
		}
		right, err := interp.evalNode(n.Right, env)
		if err != nil {
			return nil, err
		}
		return &Value{Kind: KindBool, Bool: isTruthy(right)}, nil
	case TokenOr:
		if isTruthy(left) {
			return &Value{Kind: KindBool, Bool: true}, nil
		}
		right, err := interp.evalNode(n.Right, env)
		if err != nil {
			return nil, err
		}
		return &Value{Kind: KindBool, Bool: isTruthy(right)}, nil
	}
	return nil, fmt.Errorf("unknown logical operator")
}

func (interp *Interpreter) evalCompoundAssign(n *CompoundAssignment, env *Environment) (*Value, error) {
	ident, ok := n.Target.(*Identifier)
	if !ok {
		return nil, fmt.Errorf("compound assignment requires a simple identifier")
	}
	current, exists := env.Get(ident.Name)
	if !exists {
		return nil, fmt.Errorf("undefined: %q", ident.Name)
	}
	rhs, err := interp.evalNode(n.Value, env)
	if err != nil {
		return nil, err
	}
	// Map compound operator to binary operator token
	var binOp TokenType
	switch n.Operator {
	case TokenPlusEq:
		binOp = TokenPlus
	case TokenMinusEq:
		binOp = TokenMinus
	case TokenStarEq:
		binOp = TokenStar
	case TokenSlashEq:
		binOp = TokenSlash
	case TokenPercentEq:
		binOp = TokenPercent
	default:
		return nil, fmt.Errorf("unknown compound assignment operator")
	}
	result, err := interp.evalBinaryOp(&BinaryOp{Operator: binOp, Left: &valueNode{current}, Right: &valueNode{rhs}}, env)
	if err != nil {
		return nil, err
	}
	env.Set(ident.Name, result)
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) evalIf(n *IfStatement, env *Environment) (*Value, error) {
	cond, err := interp.evalNode(n.Condition, env)
	if err != nil {
		return nil, err
	}
	if isTruthy(cond) {
		return interp.EvalStatements(n.ThenBody, NewEnvironment(env))
	}
	if len(n.ElseBody) > 0 {
		return interp.EvalStatements(n.ElseBody, NewEnvironment(env))
	}
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) evalWhile(n *WhileLoop, env *Environment) (*Value, error) {
	for {
		cond, err := interp.evalNode(n.Condition, env)
		if err != nil {
			return nil, err
		}
		if !isTruthy(cond) {
			break
		}
		_, err = interp.EvalStatements(n.Body, NewEnvironment(env))
		if err != nil {
			if _, ok := err.(breakSignal); ok {
				break
			}
			if _, ok := err.(continueSignal); ok {
				continue
			}
			return nil, err
		}
	}
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) evalFor(n *ForLoop, env *Environment) (*Value, error) {
	inner := NewEnvironment(env)
	if n.Init != nil {
		if _, err := interp.evalNode(n.Init, inner); err != nil {
			return nil, err
		}
	}
	for {
		if n.Condition != nil {
			cond, err := interp.evalNode(n.Condition, inner)
			if err != nil {
				return nil, err
			}
			if !isTruthy(cond) {
				break
			}
		}
		_, err := interp.EvalStatements(n.Body, NewEnvironment(inner))
		if err != nil {
			if _, ok := err.(breakSignal); ok {
				break
			}
			if _, ok := err.(continueSignal); ok {
				// fall through to update
			} else {
				return nil, err
			}
		}
		if n.Update != nil {
			if _, err := interp.evalNode(n.Update, inner); err != nil {
				return nil, err
			}
		}
	}
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) evalCall(n *FunctionCall, env *Environment) (*Value, error) {
	callee, ok := env.Get(n.Name)
	if !ok {
		return nil, fmt.Errorf("undefined: %q", n.Name)
	}
	args := make([]*Value, len(n.Args))
	for i, argNode := range n.Args {
		v, err := interp.evalNode(argNode, env)
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	return interp.callValue(callee, args, n.Name)
}

func (interp *Interpreter) callValue(callee *Value, args []*Value, callSite string) (*Value, error) {
	if callee.Kind != KindFunction {
		return nil, fmt.Errorf("%q is not callable", callSite)
	}
	if callee.builtin != nil {
		return callee.builtin.fn(args)
	}
	if callee.Fn == nil {
		return nil, fmt.Errorf("invalid function value for %q", callSite)
	}
	fn := callee.Fn
	fnEnv := NewEnvironment(fn.Env)
	for i, param := range fn.Params {
		var v *Value
		if i < len(args) {
			v = args[i]
		} else {
			v = &Value{Kind: KindNull}
		}
		fnEnv.Define(param.Name, v, false)
	}
	if fn.Expr != nil {
		return interp.evalNode(fn.Expr, fnEnv)
	}
	_, err := interp.EvalStatements(fn.Body, fnEnv)
	if err != nil {
		if ret, ok := err.(returnSignal); ok {
			return ret.val, nil
		}
		return nil, err
	}
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) evalPipe(n *PipeExpression, env *Environment) (*Value, error) {
	left, err := interp.evalNode(n.Left, env)
	if err != nil {
		return nil, err
	}
	callee, ok := env.Get(n.Function)
	if !ok {
		return nil, fmt.Errorf("undefined: %q", n.Function)
	}
	return interp.callValue(callee, []*Value{left}, n.Function)
}

func (interp *Interpreter) evalArrayAccess(n *ArrayAccess, env *Environment) (*Value, error) {
	arr, err := interp.evalNode(n.Array, env)
	if err != nil {
		return nil, err
	}
	idx, err := interp.evalNode(n.Index, env)
	if err != nil {
		return nil, err
	}
	if arr.Kind == KindString {
		runes := []rune(arr.Str)
		i := idx.Int
		if i < 0 || i >= len(runes) {
			return nil, fmt.Errorf("string index %d out of range [0, %d)", i, len(runes))
		}
		return &Value{Kind: KindChar, Char: runes[i]}, nil
	}
	if arr.Kind != KindArray {
		return nil, fmt.Errorf("cannot index into %s", arr.TypeName())
	}
	i := idx.Int
	if i < 0 || i >= len(arr.Arr) {
		return nil, fmt.Errorf("array index %d out of range [0, %d)", i, len(arr.Arr))
	}
	return arr.Arr[i], nil
}

func (interp *Interpreter) evalArrayDecl(n *ArrayDeclaration, env *Environment) (*Value, error) {
	var arr []*Value
	if n.Size != nil {
		sizeVal, err := interp.evalNode(n.Size, env)
		if err != nil {
			return nil, err
		}
		arr = make([]*Value, sizeVal.Int)
		for i := range arr {
			arr[i] = &Value{Kind: KindNull}
		}
	}
	env.Define(n.Name, &Value{Kind: KindArray, Arr: arr}, false)
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) defineStruct(n *StructDefinition, env *Environment) (*Value, error) {
	fields := n.Fields
	structName := n.Name
	StructRegistry[structName] = n
	bfn := &builtinFn{
		name: structName,
		fn: func(args []*Value) (*Value, error) {
			fieldMap := make(map[string]*Value)
			for i, f := range fields {
				if i < len(args) {
					fieldMap[f.Name] = args[i]
				} else {
					fieldMap[f.Name] = &Value{Kind: KindNull}
				}
			}
			return &Value{Kind: KindStruct, Str: structName, Fields: fieldMap}, nil
		},
	}
	env.Define(structName, &Value{Kind: KindFunction, builtin: bfn}, false)
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) evalStructLiteral(n *StructLiteral, env *Environment) (*Value, error) {
	fields := make(map[string]*Value)
	for name, expr := range n.Fields {
		v, err := interp.evalNode(expr, env)
		if err != nil {
			return nil, err
		}
		fields[name] = v
	}
	return &Value{Kind: KindStruct, Str: n.StructName, Fields: fields}, nil
}

func (interp *Interpreter) evalFieldAccess(n *FieldAccess, env *Environment) (*Value, error) {
	obj, err := interp.evalNode(n.Object, env)
	if err != nil {
		return nil, err
	}
	if obj.Kind != KindStruct {
		return nil, fmt.Errorf("field access on %s (expected struct)", obj.TypeName())
	}
	v, ok := obj.Fields[n.FieldName]
	if !ok {
		return nil, fmt.Errorf("no field %q in %s", n.FieldName, obj.Str)
	}
	return v, nil
}

func (interp *Interpreter) evalMatch(n *MatchExpression, env *Environment) (*Value, error) {
	subject, err := interp.evalNode(n.Value, env)
	if err != nil {
		return nil, err
	}
	for _, mc := range n.Cases {
		if mc.IsDefault {
			return interp.EvalStatements(mc.Body, NewEnvironment(env))
		}
		matched, bindEnv, err := interp.matchPattern(mc.Pattern, subject, env)
		if err != nil {
			return nil, err
		}
		if matched {
			// Check guard
			if mc.Guard != nil {
				g, err := interp.evalNode(mc.Guard, bindEnv)
				if err != nil {
					return nil, err
				}
				if !isTruthy(g) {
					continue
				}
			}
			return interp.EvalStatements(mc.Body, bindEnv)
		}
	}
	return &Value{Kind: KindVoid}, nil
}

func (interp *Interpreter) matchPattern(pattern ASTNode, subject *Value, env *Environment) (bool, *Environment, error) {
	inner := NewEnvironment(env)
	switch p := pattern.(type) {
	case *WildcardPattern:
		return true, inner, nil
	case *LiteralPattern:
		pval, err := interp.evalNode(p.Value, env)
		if err != nil {
			return false, nil, err
		}
		return valuesEqual(pval, subject), inner, nil
	case *BindingPattern:
		inner.Define(p.Name, subject, false)
		return true, inner, nil
	case *RangePattern:
		lo, err := interp.evalNode(p.Start, env)
		if err != nil {
			return false, nil, err
		}
		hi, err := interp.evalNode(p.End, env)
		if err != nil {
			return false, nil, err
		}
		sv := toFloat(subject)
		return sv >= toFloat(lo) && sv <= toFloat(hi), inner, nil
	}
	return false, inner, nil
}

func (interp *Interpreter) assign(target ASTNode, val *Value, env *Environment) (*Value, error) {
	switch t := target.(type) {
	case *Identifier:
		if err := env.Set(t.Name, val); err != nil {
			return nil, err
		}
		return &Value{Kind: KindVoid}, nil
	case *ArrayAccess:
		arr, err := interp.evalNode(t.Array, env)
		if err != nil {
			return nil, err
		}
		idx, err := interp.evalNode(t.Index, env)
		if err != nil {
			return nil, err
		}
		if arr.Kind != KindArray {
			return nil, fmt.Errorf("cannot index-assign into %s", arr.TypeName())
		}
		i := idx.Int
		if i < 0 || i >= len(arr.Arr) {
			return nil, fmt.Errorf("array index %d out of range", i)
		}
		arr.Arr[i] = val
		return &Value{Kind: KindVoid}, nil
	case *FieldAccess:
		obj, err := interp.evalNode(t.Object, env)
		if err != nil {
			return nil, err
		}
		if obj.Kind != KindStruct {
			return nil, fmt.Errorf("field assignment on %s", obj.TypeName())
		}
		obj.Fields[t.FieldName] = val
		return &Value{Kind: KindVoid}, nil
	}
	return nil, fmt.Errorf("unsupported assignment target")
}

// ============================================================================
// valueNode: lets a *Value appear as an ASTNode (for evalBinaryOp reuse)
// ============================================================================

type valueNode struct{ v *Value }

func (vn *valueNode) astNode()      {}
func (vn *valueNode) Loc() Location { return Location{} }

// evalNode handles *valueNode inline in the default branch, so we add a case
// by returning its value directly.  The switch above dispatches to default for
// unknown types; add a small helper.
func init() {
	// Ensure valueNode satisfies ASTNode at compile time.
	var _ ASTNode = (*valueNode)(nil)
}

// ============================================================================
// Utility functions
// ============================================================================

func isTruthy(v *Value) bool {
	switch v.Kind {
	case KindBool:
		return v.Bool
	case KindInt:
		return v.Int != 0
	case KindFloat:
		return v.Float != 0
	case KindString:
		return v.Str != ""
	case KindNull, KindVoid:
		return false
	default:
		return true
	}
}

func isNumeric(v *Value) bool {
	return v.Kind == KindInt || v.Kind == KindFloat
}

func toFloat(v *Value) float64 {
	switch v.Kind {
	case KindInt:
		return float64(v.Int)
	case KindFloat:
		return v.Float
	case KindBool:
		if v.Bool {
			return 1
		}
		return 0
	}
	return 0
}

func valueStr(v *Value) string {
	switch v.Kind {
	case KindString:
		return v.Str
	case KindInt:
		return fmt.Sprintf("%d", v.Int)
	case KindFloat:
		return fmt.Sprintf("%g", v.Float)
	case KindBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case KindChar:
		return string(v.Char)
	case KindNull:
		return "null"
	case KindVoid:
		return ""
	case KindArray:
		parts := make([]string, len(v.Arr))
		for i, elem := range v.Arr {
			parts[i] = valueStr(elem)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	}
	return v.Display()
}

func valuesEqual(a, b *Value) bool {
	if a.Kind != b.Kind {
		if isNumeric(a) && isNumeric(b) {
			return toFloat(a) == toFloat(b)
		}
		return false
	}
	switch a.Kind {
	case KindInt:
		return a.Int == b.Int
	case KindFloat:
		return a.Float == b.Float
	case KindBool:
		return a.Bool == b.Bool
	case KindString:
		return a.Str == b.Str
	case KindChar:
		return a.Char == b.Char
	case KindNull, KindVoid:
		return true
	}
	return false
}

func compareValues(a, b *Value, op string) bool {
	if a.Kind == KindString && b.Kind == KindString {
		switch op {
		case "<":
			return a.Str < b.Str
		case ">":
			return a.Str > b.Str
		case "<=":
			return a.Str <= b.Str
		case ">=":
			return a.Str >= b.Str
		}
	}
	av, bv := toFloat(a), toFloat(b)
	switch op {
	case "<":
		return av < bv
	case ">":
		return av > bv
	case "<=":
		return av <= bv
	case ">=":
		return av >= bv
	}
	return false
}

func mathFn1(fn func(float64) float64) func([]*Value) (*Value, error) {
	return func(args []*Value) (*Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("expects 1 argument")
		}
		return &Value{Kind: KindFloat, Float: fn(toFloat(args[0]))}, nil
	}
}
