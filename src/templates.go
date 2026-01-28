package main

import "fmt"

// templates.go - AST nodes and types for C++-style templates

// ============================================================================
// Template AST Nodes
// ============================================================================

// TemplateParameter represents a parameter in a template declaration
type TemplateParameter struct {
	IsType  bool        // true for typename/type parameters, false for value parameters
	Name    string      // Parameter name (e.g., "T", "N")
	Type    TokenType   // For value parameters (e.g., int N)
	Default interface{} // Default value/type if specified
}

// TemplateDeclaration represents a template<...> declaration header
type TemplateDeclaration struct {
	BaseNode
	Parameters  []TemplateParameter // Template parameters
	Declaration ASTNode             // The templated function/struct/class
}

func (t *TemplateDeclaration) astNode() {}

// TemplateInstantiation represents a template being instantiated with concrete arguments
// Example: vector<int>, max<float>
type TemplateInstantiation struct {
	BaseNode
	TemplateName string    // Name of the template
	TypeArgs     []string  // Type arguments (e.g., ["int"], ["string", "int"])
	ValueArgs    []ASTNode // Value arguments for non-type parameters
}

func (t *TemplateInstantiation) astNode() {}

// TemplateSpecialization represents a template specialization
// template<> fn void foo<int>() { ... }
type TemplateSpecialization struct {
	BaseNode
	TemplateName string    // Name of template being specialized
	TypeArgs     []string  // Specialized type arguments
	ValueArgs    []ASTNode // Specialized value arguments
	Declaration  ASTNode   // The specialized definition
	IsPartial    bool      // true for partial specialization
}

func (t *TemplateSpecialization) astNode() {}

// ============================================================================
// Template Type System
// ============================================================================

// TemplateDefinition stores a template definition before instantiation
type TemplateDefinition struct {
	Name            string
	Parameters      []TemplateParameter
	Body            ASTNode            // FunctionDefinition, StructDefinition, or ClassDefinition
	Instantiations  map[string]ASTNode // Cache: "int,10" -> instantiated version
	Specializations []*TemplateSpecialization
}

// TemplateContext tracks template parameters during instantiation
type TemplateContext struct {
	TypeBindings  map[string]TokenType   // "T" -> TokenTypeInt
	ValueBindings map[string]interface{} // "N" -> 10
}

// ============================================================================
// Namespace AST Nodes
// ============================================================================

// NamespaceDeclaration represents a namespace { ... } block
type NamespaceDeclaration struct {
	BaseNode
	Name     string    // Namespace name
	Body     []ASTNode // Declarations inside namespace
	IsInline bool      // true for inline namespace
}

func (n *NamespaceDeclaration) astNode() {}

// QualifiedName represents a namespace-qualified identifier
// Example: std::vector, math::max
type QualifiedName struct {
	BaseNode
	Namespaces []string // ["std", "collections"]
	Name       string   // "vector"
}

func (q *QualifiedName) astNode() {}

// UsingDeclaration represents a using statement
// using std::vector;
// using namespace std;
type UsingDeclaration struct {
	BaseNode
	IsNamespace bool     // true for "using namespace", false for specific item
	Namespaces  []string // Namespace path (e.g., ["std", "collections"])
	Name        string   // Item name (empty for "using namespace")
	Alias       string   // Optional alias (e.g., "using std::vector as vec")
}

func (u *UsingDeclaration) astNode() {}

// ============================================================================
// Namespace Type System
// ============================================================================

// Namespace represents a namespace scope
type Namespace struct {
	Name      string
	Parent    *Namespace
	Children  map[string]*Namespace
	Functions map[string][]*FunctionDefinition // Support overloading
	Templates map[string]*TemplateDefinition
	Structs   map[string]*StructDefinition
	Classes   map[string]*ClassDefinition
	Enums     map[string]*EnumDefinition
	Variables map[string]*VariableDeclaration
	Constants map[string]*ConstantDeclaration
}

// NewNamespace creates a new namespace
func NewNamespace(name string, parent *Namespace) *Namespace {
	return &Namespace{
		Name:      name,
		Parent:    parent,
		Children:  make(map[string]*Namespace),
		Functions: make(map[string][]*FunctionDefinition),
		Templates: make(map[string]*TemplateDefinition),
		Structs:   make(map[string]*StructDefinition),
		Classes:   make(map[string]*ClassDefinition),
		Enums:     make(map[string]*EnumDefinition),
		Variables: make(map[string]*VariableDeclaration),
		Constants: make(map[string]*ConstantDeclaration),
	}
}

// FullName returns the fully qualified name of this namespace
func (ns *Namespace) FullName() string {
	if ns.Parent == nil || ns.Parent.Name == "" {
		return ns.Name
	}
	return ns.Parent.FullName() + "::" + ns.Name
}

// Lookup searches for a name in this namespace and parent namespaces
func (ns *Namespace) Lookup(name string) interface{} {
	// Check functions
	if funcs, ok := ns.Functions[name]; ok && len(funcs) > 0 {
		return funcs
	}

	// Check templates
	if tmpl, ok := ns.Templates[name]; ok {
		return tmpl
	}

	// Check types
	if s, ok := ns.Structs[name]; ok {
		return s
	}
	if c, ok := ns.Classes[name]; ok {
		return c
	}
	if e, ok := ns.Enums[name]; ok {
		return e
	}

	// Check variables and constants
	if v, ok := ns.Variables[name]; ok {
		return v
	}
	if c, ok := ns.Constants[name]; ok {
		return c
	}

	// Check child namespaces
	if child, ok := ns.Children[name]; ok {
		return child
	}

	// Search parent namespace
	if ns.Parent != nil {
		return ns.Parent.Lookup(name)
	}

	return nil
}

// GetOrCreateChild gets or creates a child namespace
func (ns *Namespace) GetOrCreateChild(name string) *Namespace {
	if child, ok := ns.Children[name]; ok {
		return child
	}
	child := NewNamespace(name, ns)
	ns.Children[name] = child
	return child
}

// ============================================================================
// Name Mangling
// ============================================================================

// MangleTemplateName creates a mangled name for a template instantiation
// Example: vector<int> -> _ZTS6vectorIiE
func MangleTemplateName(name string, typeArgs []string, valueArgs []interface{}) string {
	mangled := "_ZTS" + string(rune(len(name))) + name + "I"

	for _, typeArg := range typeArgs {
		mangled += mangleTypeName(typeArg)
	}

	for _, valueArg := range valueArgs {
		mangled += mangleValue(valueArg)
	}

	mangled += "E"
	return mangled
}

// MangleNamespaceName creates a mangled name for a namespaced function
// Example: math::add -> _ZN4math3addE
func MangleNamespaceName(namespaces []string, name string) string {
	if len(namespaces) == 0 {
		return name
	}

	mangled := "_ZN"
	for _, ns := range namespaces {
		mangled += string(rune(len(ns))) + ns
	}
	mangled += string(rune(len(name))) + name
	mangled += "E"
	return mangled
}

// mangleTypeName converts a type name to its mangled representation
func mangleTypeName(typeName string) string {
	switch typeName {
	case "void":
		return "v"
	case "bool":
		return "b"
	case "char":
		return "c"
	case "int8":
		return "a"
	case "uint8":
		return "h"
	case "int16":
		return "s"
	case "uint16":
		return "t"
	case "int", "int32":
		return "i"
	case "uint", "uint32":
		return "j"
	case "int64":
		return "l"
	case "uint64":
		return "m"
	case "float32":
		return "f"
	case "float64":
		return "d"
	case "string":
		return "Ss" // std::string equivalent
	default:
		// Custom type: use length + name
		return string(rune(len(typeName))) + typeName
	}
}

// mangleValue converts a value to its mangled representation
func mangleValue(value interface{}) string {
	switch v := value.(type) {
	case int:
		if v >= 0 {
			return "Li" + string(rune(v)) + "E"
		} else {
			return "Lin" + string(rune(-v)) + "E"
		}
	case int64:
		if v >= 0 {
			return "Ll" + string(rune(v)) + "E"
		} else {
			return "Lln" + string(rune(-v)) + "E"
		}
	default:
		return "Lv0E" // Unknown value
	}
}

// ============================================================================
// Template Instantiation Engine
// ============================================================================

// TemplateInstantiator handles template instantiation
type TemplateInstantiator struct {
	registry  *TypeRegistry
	templates map[string]*TemplateDefinition
}

// NewTemplateInstantiator creates a new template instantiator
func NewTemplateInstantiator(registry *TypeRegistry) *TemplateInstantiator {
	return &TemplateInstantiator{
		registry:  registry,
		templates: make(map[string]*TemplateDefinition),
	}
}

// RegisterTemplate registers a template definition
func (ti *TemplateInstantiator) RegisterTemplate(name string, tmpl *TemplateDefinition) {
	ti.templates[name] = tmpl
}

// Instantiate instantiates a template with the given arguments
func (ti *TemplateInstantiator) Instantiate(name string, typeArgs []string, valueArgs []interface{}) (ASTNode, error) {
	tmpl, ok := ti.templates[name]
	if !ok {
		return nil, fmt.Errorf("template '%s' not found", name)
	}

	// Create cache key
	cacheKey := createCacheKey(typeArgs, valueArgs)

	// Check if already instantiated
	if inst, ok := tmpl.Instantiations[cacheKey]; ok {
		return inst, nil
	}

	// Check for matching specialization
	for _, spec := range tmpl.Specializations {
		if matchesSpecialization(spec, typeArgs, valueArgs) {
			return spec.Declaration, nil
		}
	}

	// Create template context with bindings
	ctx := &TemplateContext{
		TypeBindings:  make(map[string]TokenType),
		ValueBindings: make(map[string]interface{}),
	}

	// Bind type parameters
	typeParamIdx := 0
	valueParamIdx := 0
	for _, param := range tmpl.Parameters {
		if param.IsType {
			if typeParamIdx < len(typeArgs) {
				// Convert type name to TokenType
				tokenType := stringToTokenType(typeArgs[typeParamIdx])
				ctx.TypeBindings[param.Name] = tokenType
				typeParamIdx++
			}
		} else {
			if valueParamIdx < len(valueArgs) {
				ctx.ValueBindings[param.Name] = valueArgs[valueParamIdx]
				valueParamIdx++
			}
		}
	}

	// Instantiate the template body with the context
	instantiated := ti.instantiateNode(tmpl.Body, ctx)

	// Cache the instantiation
	tmpl.Instantiations[cacheKey] = instantiated

	return instantiated, nil
}

// instantiateNode recursively instantiates a node with template parameter substitutions
func (ti *TemplateInstantiator) instantiateNode(node ASTNode, ctx *TemplateContext) ASTNode {
	// This would be a deep copy of the AST node with all template parameters replaced
	// For now, return the node as-is (full implementation would recursively copy and substitute)
	return node
}

// Helper functions

func createCacheKey(typeArgs []string, valueArgs []interface{}) string {
	key := ""
	for i, ta := range typeArgs {
		if i > 0 {
			key += ","
		}
		key += ta
	}
	if len(valueArgs) > 0 {
		key += ":"
		for i, va := range valueArgs {
			if i > 0 {
				key += ","
			}
			key += fmt.Sprintf("%v", va)
		}
	}
	return key
}

func matchesSpecialization(spec *TemplateSpecialization, typeArgs []string, valueArgs []interface{}) bool {
	if len(spec.TypeArgs) != len(typeArgs) {
		return false
	}
	for i, ta := range typeArgs {
		if spec.TypeArgs[i] != ta && spec.TypeArgs[i] != "" {
			return false
		}
	}
	// TODO: Check value args
	return true
}

func stringToTokenType(typeName string) TokenType {
	switch typeName {
	case "int":
		return TokenTypeInt
	case "int8":
		return TokenTypeInt8
	case "int16":
		return TokenTypeInt16
	case "int32":
		return TokenTypeInt32
	case "int64":
		return TokenTypeInt64
	case "uint":
		return TokenTypeUint
	case "uint8":
		return TokenTypeUint8
	case "uint16":
		return TokenTypeUint16
	case "uint32":
		return TokenTypeUint32
	case "uint64":
		return TokenTypeUint64
	case "float32":
		return TokenTypeFloat32
	case "float64":
		return TokenTypeFloat64
	case "string":
		return TokenTypeString
	case "char":
		return TokenTypeChar
	case "bool":
		return TokenTypeBool
	case "void":
		return TokenTypeVoid
	default:
		return TokenIdentifier // Custom type
	}
}
