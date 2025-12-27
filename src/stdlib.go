package main

import (
	"fmt"
)

// stdlib.go - Standard Library module system for Lotus
// Defines available stdlib modules and their exported functions/types.

// StdlibModule represents a module in the standard library
type StdlibModule struct {
	Name      string                     // Module name (e.g., "io", "math")
	Functions map[string]*StdlibFunction // Available functions in this module
	Types     map[string]TokenType       // Available types (future)
}

// StdlibFunction represents a function available in the stdlib
type StdlibFunction struct {
	Name     string
	Module   string // Module it belongs to
	NumArgs  int    // -1 for variadic
	ArgTypes []TokenType
	RetType  TokenType
	CodeGen  func(*CodeGenerator, []ASTNode) // Code generation function
}

// StandardLibrary holds all available stdlib modules
var StandardLibrary = map[string]*StdlibModule{
	"io":          createIOModule(),
	"mem":         createMemoryModule(),
	"math":        createMathModule(),
	"str":         createStringModule(),
	"num":         createNumModule(),
	"hash":        createHashModule(),
	"collections": createCollectionsModule(),
	"net":         createNetModule(),
	"http":        createHTTPModule(),
	"file":        createFileModule(),
	"time":        createTimeModule(),
}

// createIOModule creates the I/O standard library module
func createIOModule() *StdlibModule {
	return &StdlibModule{
		Name: "io",
		Functions: map[string]*StdlibFunction{
			"print": {
				Name:    "print",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generatePrintCode,
			},
			"println": {
				Name:    "println",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generatePrintlnCode,
			},
			"printf": {
				Name:    "printf",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOPrintf,
			},
			"fprintf": {
				Name:    "fprintf",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOFprintf,
			},
			"sprint": {
				Name:    "sprint",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOSprint,
			},
			"sprintf": {
				Name:    "sprintf",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOSprintf,
			},
			"sprintln": {
				Name:    "sprintln",
				Module:  "io",
				NumArgs: -1,
				CodeGen: generateIOSprintln,
			},
		},
		Types: map[string]TokenType{},
	}
}

// createMemoryModule creates the memory management stdlib module
func createMemoryModule() *StdlibModule {
	return &StdlibModule{
		Name: "mem",
		Functions: map[string]*StdlibFunction{
			"malloc": {
				Name:    "malloc",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemMalloc,
			},
			"free": {
				Name:    "free",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemFreeNoop,
			},
			"sizeof": {
				Name:    "sizeof",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemSizeof,
			},
			"memcpy": {
				Name:    "memcpy",
				Module:  "mem",
				NumArgs: 3,
				CodeGen: generateMemMemcpy,
			},
			"memset": {
				Name:    "memset",
				Module:  "mem",
				NumArgs: 3,
				CodeGen: generateMemMemset,
			},
			"mmap": {
				Name:    "mmap",
				Module:  "mem",
				NumArgs: 1,
				CodeGen: generateMemMmap,
			},
			"munmap": {
				Name:    "munmap",
				Module:  "mem",
				NumArgs: 2,
				CodeGen: generateMemMunmap,
			},
		},
		Types: map[string]TokenType{},
	}
}

// createMathModule creates the math stdlib module
func createMathModule() *StdlibModule {
	return &StdlibModule{
		Name: "math",
		Functions: map[string]*StdlibFunction{
			"abs": {
				Name:    "abs",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathAbs,
			},
			"min": {
				Name:    "min",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathMin,
			},
			"max": {
				Name:    "max",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathMax,
			},
			"sqrt": {
				Name:    "sqrt",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathSqrt,
			},
			"pow": {
				Name:    "pow",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathPow,
			},
			"floor": {
				Name:    "floor",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathFloor,
			},
			"ceil": {
				Name:    "ceil",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathCeil,
			},
			"round": {
				Name:    "round",
				Module:  "math",
				NumArgs: 1,
				CodeGen: generateMathRound,
			},
			"gcd": {
				Name:    "gcd",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathGcd,
			},
			"lcm": {
				Name:    "lcm",
				Module:  "math",
				NumArgs: 2,
				CodeGen: generateMathLcm,
			},
		},
		Types: map[string]TokenType{},
	}
}

// createStringModule creates the string manipulation stdlib module
func createStringModule() *StdlibModule {
	return &StdlibModule{
		Name: "str",
		Functions: map[string]*StdlibFunction{
			"len": {
				Name:    "len",
				Module:  "str",
				NumArgs: 1,
				CodeGen: generateStringLen,
			},
			"concat": {
				Name:    "concat",
				Module:  "str",
				NumArgs: -1,
				CodeGen: generateStringConcat,
			},
			"compare": {
				Name:    "compare",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringCompare,
			},
			"copy": {
				Name:    "copy",
				Module:  "str",
				NumArgs: 1,
				CodeGen: generateStringCopy,
			},
			"indexOf": {
				Name:    "indexOf",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringIndexOf,
			},
			"contains": {
				Name:    "contains",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringContains,
			},
			"startsWith": {
				Name:    "startsWith",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringStartsWith,
			},
			"endsWith": {
				Name:    "endsWith",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringEndsWith,
			},
			"substring": {
				Name:    "substring",
				Module:  "str",
				NumArgs: 3,
				CodeGen: generateStringSubstring,
			},
			"split": {
				Name:    "split",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringSplit,
			},
			"join": {
				Name:    "join",
				Module:  "str",
				NumArgs: 2,
				CodeGen: generateStringJoin,
			},
			"replace": {
				Name:    "replace",
				Module:  "str",
				NumArgs: 3,
				CodeGen: generateStringReplace,
			},
			"toLower": {
				Name:    "toLower",
				Module:  "str",
				NumArgs: 1,
				CodeGen: generateStringToLower,
			},
			"toUpper": {
				Name:    "toUpper",
				Module:  "str",
				NumArgs: 1,
				CodeGen: generateStringToUpper,
			},
			"trim": {
				Name:    "trim",
				Module:  "str",
				NumArgs: 1,
				CodeGen: generateStringTrim,
			},
		},
		Types: map[string]TokenType{},
	}
}

func createNetModule() *StdlibModule {
	return &StdlibModule{
		Name: "net",
		Functions: map[string]*StdlibFunction{
			"socket":       {Name: "socket", Module: "net", NumArgs: 3, CodeGen: generateNetSocket},
			"connect_ipv4": {Name: "connect_ipv4", Module: "net", NumArgs: 3, CodeGen: generateNetConnectIPv4},
			"send":         {Name: "send", Module: "net", NumArgs: 3, CodeGen: generateNetSend},
			"recv":         {Name: "recv", Module: "net", NumArgs: 3, CodeGen: generateNetRecv},
			"close":        {Name: "close", Module: "net", NumArgs: 1, CodeGen: generateNetClose},
			// UDP support
			"bind_ipv4":   {Name: "bind_ipv4", Module: "net", NumArgs: 3, CodeGen: generateNetBindIPv4},
			"sendto_ipv4": {Name: "sendto_ipv4", Module: "net", NumArgs: 5, CodeGen: generateNetSendtoIPv4},
			"recvfrom":    {Name: "recvfrom", Module: "net", NumArgs: 3, CodeGen: generateNetRecvfrom},
			// IPv6 support
			"connect_ipv6": {Name: "connect_ipv6", Module: "net", NumArgs: 3, CodeGen: generateNetConnectIPv6},
			"bind_ipv6":    {Name: "bind_ipv6", Module: "net", NumArgs: 3, CodeGen: generateNetBindIPv6},
			"sendto_ipv6":  {Name: "sendto_ipv6", Module: "net", NumArgs: 5, CodeGen: generateNetSendtoIPv6},
			// DNS resolution
			"resolve":      {Name: "resolve", Module: "net", NumArgs: 2, CodeGen: generateNetResolve},
			"resolve_ipv6": {Name: "resolve_ipv6", Module: "net", NumArgs: 2, CodeGen: generateNetResolveIPv6},
		},
		Types: map[string]TokenType{},
	}
}

// createHTTPModule creates a minimal HTTP client module built on net helpers
func createHTTPModule() *StdlibModule {
	return &StdlibModule{
		Name: "http",
		Functions: map[string]*StdlibFunction{
			"get":  {Name: "get", Module: "http", NumArgs: 7, CodeGen: generateHTTPGetSimple},
			"post": {Name: "post", Module: "http", NumArgs: 9, CodeGen: generateHTTPPostSimple},
			// Response parsing
			"parse_status":  {Name: "parse_status", Module: "http", NumArgs: 2, CodeGen: generateHTTPParseStatus},
			"get_header":    {Name: "get_header", Module: "http", NumArgs: 4, CodeGen: generateHTTPGetHeader},
			"get_body":      {Name: "get_body", Module: "http", NumArgs: 2, CodeGen: generateHTTPGetBody},
			"parse_headers": {Name: "parse_headers", Module: "http", NumArgs: 3, CodeGen: generateHTTPParseHeaders},
			// Connection pooling
			"pool_new":   {Name: "pool_new", Module: "http", NumArgs: 1, CodeGen: generateHTTPPoolNew},     // pool_new(max_conns) -> pool_ptr
			"pool_get":   {Name: "pool_get", Module: "http", NumArgs: 3, CodeGen: generateHTTPPoolGet},     // pool_get(pool, host_ptr, port) -> fd or -1
			"pool_put":   {Name: "pool_put", Module: "http", NumArgs: 4, CodeGen: generateHTTPPoolPut},     // pool_put(pool, fd, host_ptr, port) -> 0/1
			"pool_close": {Name: "pool_close", Module: "http", NumArgs: 1, CodeGen: generateHTTPPoolClose}, // pool_close(pool) -> void
		},
		Types: map[string]TokenType{},
	}
}

// createCollectionsModule creates the collections stdlib module (data structures)
func createCollectionsModule() *StdlibModule {
	return &StdlibModule{
		Name: "collections",
		Functions: map[string]*StdlibFunction{
			// Dynamic array
			"array_int_new":      {Name: "array_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntNew},
			"array_int_push":     {Name: "array_int_push", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsArrayIntPush},
			"array_int_pop":      {Name: "array_int_pop", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntPop},
			"array_int_len":      {Name: "array_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntLen},
			"array_int_capacity": {Name: "array_int_capacity", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntCapacity},
			"array_int_resize":   {Name: "array_int_resize", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsArrayIntResize},
			"array_int_reserve":  {Name: "array_int_reserve", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsArrayIntReserve},
			"array_int_shrink":   {Name: "array_int_shrink", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntShrink},
			"array_int_get":      {Name: "array_int_get", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsArrayIntGet},
			"array_int_set":      {Name: "array_int_set", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsArrayIntSet},
			"array_int_free":     {Name: "array_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsArrayIntFree},

			// Stack
			"stack_int_new":  {Name: "stack_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsStackIntNew},
			"stack_int_push": {Name: "stack_int_push", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsStackIntPush},
			"stack_int_pop":  {Name: "stack_int_pop", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsStackIntPop},
			"stack_int_len":  {Name: "stack_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsStackIntLen},

			// Queue / Deque
			"queue_int_new":     {Name: "queue_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsQueueIntNew},
			"queue_int_enqueue": {Name: "queue_int_enqueue", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsQueueIntEnqueue},
			"queue_int_dequeue": {Name: "queue_int_dequeue", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsQueueIntDequeue},
			"queue_int_len":     {Name: "queue_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsQueueIntLen},

			"deque_int_new":        {Name: "deque_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntNew},
			"deque_int_push_front": {Name: "deque_int_push_front", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsDequeIntPushFront},
			"deque_int_push_back":  {Name: "deque_int_push_back", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsDequeIntPushBack},
			"deque_int_pop_front":  {Name: "deque_int_pop_front", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntPopFront},
			"deque_int_pop_back":   {Name: "deque_int_pop_back", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntPopBack},
			"deque_int_len":        {Name: "deque_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsDequeIntLen},

			// Heap (min-heap)
			"heap_int_new":  {Name: "heap_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntNew},
			"heap_int_push": {Name: "heap_int_push", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHeapIntPush},
			"heap_int_pop":  {Name: "heap_int_pop", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntPop},
			"heap_int_peek": {Name: "heap_int_peek", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntPeek},
			"heap_int_len":  {Name: "heap_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHeapIntLen},

			// Hash map & set (int keys)
			"hashmap_int_new":    {Name: "hashmap_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntNew},
			"hashmap_int_put":    {Name: "hashmap_int_put", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsHashmapIntPut},
			"hashmap_int_get":    {Name: "hashmap_int_get", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapIntGet},
			"hashmap_int_remove": {Name: "hashmap_int_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapIntRemove},
			"hashmap_int_len":    {Name: "hashmap_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntLen},
			"hashmap_int_clear":  {Name: "hashmap_int_clear", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntClear},
			"hashmap_int_free":   {Name: "hashmap_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapIntFree},

			"hashset_int_new":      {Name: "hashset_int_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntNew},
			"hashset_int_add":      {Name: "hashset_int_add", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetIntAdd},
			"hashset_int_contains": {Name: "hashset_int_contains", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetIntContains},
			"hashset_int_remove":   {Name: "hashset_int_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetIntRemove},
			"hashset_int_len":      {Name: "hashset_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntLen},
			"hashset_int_clear":    {Name: "hashset_int_clear", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntClear},
			"hashset_int_free":     {Name: "hashset_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetIntFree},

			// Hash map & set (string keys)
			"hashmap_str_new":      {Name: "hashmap_str_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapStrNew},
			"hashmap_str_put":      {Name: "hashmap_str_put", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsHashmapStrPut},
			"hashmap_str_get":      {Name: "hashmap_str_get", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapStrGet},
			"hashmap_str_contains": {Name: "hashmap_str_contains", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapStrContains},
			"hashmap_str_remove":   {Name: "hashmap_str_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashmapStrRemove},
			"hashmap_str_len":      {Name: "hashmap_str_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapStrLen},
			"hashmap_str_clear":    {Name: "hashmap_str_clear", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapStrClear},
			"hashmap_str_free":     {Name: "hashmap_str_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashmapStrFree},

			"hashset_str_new":      {Name: "hashset_str_new", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetStrNew},
			"hashset_str_add":      {Name: "hashset_str_add", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetStrAdd},
			"hashset_str_contains": {Name: "hashset_str_contains", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetStrContains},
			"hashset_str_remove":   {Name: "hashset_str_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsHashsetStrRemove},
			"hashset_str_len":      {Name: "hashset_str_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetStrLen},
			"hashset_str_clear":    {Name: "hashset_str_clear", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetStrClear},
			"hashset_str_free":     {Name: "hashset_str_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsHashsetStrFree},

			// Sorted set (BST-based, maintains sorted order)
			"sortedset_int_new":      {Name: "sortedset_int_new", Module: "collections", NumArgs: 0, CodeGen: generateCollectionsSortedsetIntNew},
			"sortedset_int_add":      {Name: "sortedset_int_add", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsSortedsetIntAdd},
			"sortedset_int_contains": {Name: "sortedset_int_contains", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsSortedsetIntContains},
			"sortedset_int_remove":   {Name: "sortedset_int_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsSortedsetIntRemove},
			"sortedset_int_min":      {Name: "sortedset_int_min", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedsetIntMin},
			"sortedset_int_max":      {Name: "sortedset_int_max", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedsetIntMax},
			"sortedset_int_len":      {Name: "sortedset_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedsetIntLen},
			"sortedset_int_free":     {Name: "sortedset_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedsetIntFree},

			// Sorted map (BST-based, maintains sorted order by key)
			"sortedmap_int_new":      {Name: "sortedmap_int_new", Module: "collections", NumArgs: 0, CodeGen: generateCollectionsSortedmapIntNew},
			"sortedmap_int_put":      {Name: "sortedmap_int_put", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsSortedmapIntPut},
			"sortedmap_int_get":      {Name: "sortedmap_int_get", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsSortedmapIntGet},
			"sortedmap_int_contains": {Name: "sortedmap_int_contains", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsSortedmapIntContains},
			"sortedmap_int_remove":   {Name: "sortedmap_int_remove", Module: "collections", NumArgs: 2, CodeGen: generateCollectionsSortedmapIntRemove},
			"sortedmap_int_min_key":  {Name: "sortedmap_int_min_key", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedmapIntMinKey},
			"sortedmap_int_max_key":  {Name: "sortedmap_int_max_key", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedmapIntMaxKey},
			"sortedmap_int_len":      {Name: "sortedmap_int_len", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedmapIntLen},
			"sortedmap_int_free":     {Name: "sortedmap_int_free", Module: "collections", NumArgs: 1, CodeGen: generateCollectionsSortedmapIntFree},

			// Array helper
			"binary_search_int": {Name: "binary_search_int", Module: "collections", NumArgs: 3, CodeGen: generateCollectionsBinarySearchInt},
		},
		Types: map[string]TokenType{},
	}
}

// createNumModule creates the numeric conversions stdlib module
func createNumModule() *StdlibModule {
	return &StdlibModule{
		Name: "num",
		Functions: map[string]*StdlibFunction{
			"toInt8":   {Name: "toInt8", Module: "num", NumArgs: 1, CodeGen: generateNumToInt8},
			"toUint8":  {Name: "toUint8", Module: "num", NumArgs: 1, CodeGen: generateNumToUint8},
			"toInt16":  {Name: "toInt16", Module: "num", NumArgs: 1, CodeGen: generateNumToInt16},
			"toUint16": {Name: "toUint16", Module: "num", NumArgs: 1, CodeGen: generateNumToUint16},
			"toInt32":  {Name: "toInt32", Module: "num", NumArgs: 1, CodeGen: generateNumToInt32},
			"toUint32": {Name: "toUint32", Module: "num", NumArgs: 1, CodeGen: generateNumToUint32},
			"toInt64":  {Name: "toInt64", Module: "num", NumArgs: 1, CodeGen: generateNumToInt64},
			"toUint64": {Name: "toUint64", Module: "num", NumArgs: 1, CodeGen: generateNumToUint64},
			"toBool":   {Name: "toBool", Module: "num", NumArgs: 1, CodeGen: generateNumToBool},
		},
		Types: map[string]TokenType{},
	}
}

// createHashModule creates the hashing stdlib module (cryptographic and non-cryptographic)
func createHashModule() *StdlibModule {
	return &StdlibModule{
		Name: "hash",
		Functions: map[string]*StdlibFunction{
			// Non-cryptographic hashes (fast, simple)
			"crc32":  {Name: "crc32", Module: "hash", NumArgs: 2, CodeGen: generateHashCRC32},    // crc32(data_ptr, len) -> uint32
			"fnv1a":  {Name: "fnv1a", Module: "hash", NumArgs: 2, CodeGen: generateHashFNV1a},    // fnv1a(data_ptr, len) -> uint64
			"djb2":   {Name: "djb2", Module: "hash", NumArgs: 1, CodeGen: generateHashDJB2},      // djb2(string_ptr) -> uint64
			"murmur": {Name: "murmur", Module: "hash", NumArgs: 3, CodeGen: generateHashMurmur3}, // murmur(data_ptr, len, seed) -> uint32

			// Cryptographic hashes
			"sha256": {Name: "sha256", Module: "hash", NumArgs: 3, CodeGen: generateHashSHA256}, // sha256(data_ptr, len, out_buf) -> void
			"md5":    {Name: "md5", Module: "hash", NumArgs: 3, CodeGen: generateHashMD5},       // md5(data_ptr, len, out_buf) -> void
		},
		Types: map[string]TokenType{},
	}
}

// createFileModule creates a file I/O stdlib module (POSIX file operations)
func createFileModule() *StdlibModule {
	return &StdlibModule{
		Name: "file",
		Functions: map[string]*StdlibFunction{
			"open":   {Name: "open", Module: "file", NumArgs: 2, CodeGen: generateFileOpen},     // open(path_ptr, flags) -> fd
			"close":  {Name: "close", Module: "file", NumArgs: 1, CodeGen: generateFileClose},   // close(fd) -> status
			"read":   {Name: "read", Module: "file", NumArgs: 3, CodeGen: generateFileRead},     // read(fd, buf_ptr, size) -> bytes_read
			"write":  {Name: "write", Module: "file", NumArgs: 3, CodeGen: generateFileWrite},   // write(fd, buf_ptr, size) -> bytes_written
			"seek":   {Name: "seek", Module: "file", NumArgs: 3, CodeGen: generateFileSeek},     // seek(fd, offset, whence) -> new_pos
			"stat":   {Name: "stat", Module: "file", NumArgs: 2, CodeGen: generateFileStat},     // stat(path_ptr, stat_buf) -> status
			"exists": {Name: "exists", Module: "file", NumArgs: 1, CodeGen: generateFileExists}, // exists(path_ptr) -> 0/1
		},
		Types: map[string]TokenType{},
	}
}

// createTimeModule creates a time/date utility stdlib module
func createTimeModule() *StdlibModule {
	return &StdlibModule{
		Name: "time",
		Functions: map[string]*StdlibFunction{
			"now":       {Name: "now", Module: "time", NumArgs: 0, CodeGen: generateTimeNow},             // now() -> unix_timestamp
			"sleep":     {Name: "sleep", Module: "time", NumArgs: 1, CodeGen: generateTimeSleep},         // sleep(seconds) -> status
			"millis":    {Name: "millis", Module: "time", NumArgs: 0, CodeGen: generateTimeMillis},       // millis() -> milliseconds
			"nanos":     {Name: "nanos", Module: "time", NumArgs: 0, CodeGen: generateTimeNanos},         // nanos() -> nanoseconds
			"clock":     {Name: "clock", Module: "time", NumArgs: 0, CodeGen: generateTimeClock},         // clock() -> clock_ticks
			"gmtime":    {Name: "gmtime", Module: "time", NumArgs: 2, CodeGen: generateTimeGMTime},       // gmtime(timestamp, tm_buf) -> void
			"localtime": {Name: "localtime", Module: "time", NumArgs: 2, CodeGen: generateTimeLocalTime}, // localtime(timestamp, tm_buf) -> void
		},
		Types: map[string]TokenType{},
	}
}

// stdlibLookup is a function pointer for late-binding module lookups
// This is set after StandardLibrary is initialized to avoid init cycles
var stdlibLookup func(moduleName, funcName string) *StdlibFunction

func init() {
	// Set up the lookup function after StandardLibrary is fully initialized
	stdlibLookup = func(moduleName, funcName string) *StdlibFunction {
		if module, ok := StandardLibrary[moduleName]; ok {
			if fn, ok := module.Functions[funcName]; ok {
				return fn
			}
		}
		return nil
	}
}

// GetModuleFunction retrieves a function from a module
func GetModuleFunction(moduleName, funcName string) *StdlibFunction {
	if stdlibLookup != nil {
		return stdlibLookup(moduleName, funcName)
	}
	return nil
}

// ImportContext tracks what has been imported in the current compilation
type ImportContext struct {
	ImportedModules   map[string]string          // Maps alias to module name
	ImportedFunctions map[string]*StdlibFunction // Maps function name to function
	UseWildcard       bool                       // true if using wildcard import
}

// NewImportContext creates a new import tracking context
func NewImportContext() *ImportContext {
	return &ImportContext{
		ImportedModules:   make(map[string]string),
		ImportedFunctions: make(map[string]*StdlibFunction),
	}
}

// ProcessImport processes an import statement and adds exported items to context
func (ic *ImportContext) ProcessImport(stmt *ImportStatement) error {
	module, exists := StandardLibrary[stmt.Module]
	if !exists {
		return fmt.Errorf("module '%s' not found in standard library", stmt.Module)
	}

	alias := stmt.Alias
	if alias == "" {
		alias = stmt.Module
	}

	ic.ImportedModules[alias] = stmt.Module

	// Process specific imports
	if stmt.IsWildcard {
		// Import all functions from module
		for name, fn := range module.Functions {
			ic.ImportedFunctions[name] = fn
		}
	} else if len(stmt.Items) > 0 {
		// Import specific items
		for _, item := range stmt.Items {
			if fn, ok := module.Functions[item]; ok {
				ic.ImportedFunctions[item] = fn
			}
		}
	} else {
		// No specific items means import all
		for name, fn := range module.Functions {
			ic.ImportedFunctions[name] = fn
		}
	}

	return nil
}

// ============================================================================
// Placeholder code generation functions for stdlib functions
// ============================================================================

func generateMathAbs(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// abs(x): if x < 0 then -x else x
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    sarq $63, %rbx\n")
	cg.textSection.WriteString("    xorq %rbx, %rax\n")
	cg.textSection.WriteString("    subq %rbx, %rax\n")
}

func generateMathMin(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	labelGT := cg.getLabel("min_gt")
	labelEnd := cg.getLabel("min_end")
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.generateExpressionToReg(args[1], "rax")
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString("    cmpq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jle %s\n", labelGT))
	// rcx smaller
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelGT))
	// rax already min
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
}

func generateMathMax(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	labelLT := cg.getLabel("max_lt")
	labelEnd := cg.getLabel("max_end")
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.generateExpressionToReg(args[1], "rax")
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString("    cmpq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", labelLT))
	// rcx larger
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelLT))
	// rax already max
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
}

func generateMathSqrt(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// Integer floor sqrt using hardware sqrt on double; negative -> -1
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    js 2f\n")
	cg.textSection.WriteString("    cvtsi2sd %rax, %xmm0\n")
	cg.textSection.WriteString("    sqrtsd %xmm0, %xmm0\n")
	cg.textSection.WriteString("    cvttsd2si %xmm0, %rax\n")
	cg.textSection.WriteString("    jmp 3f\n")
	cg.textSection.WriteString("2:  movq $-1, %rax\n")
	cg.textSection.WriteString("3:\n")
}

func generateMathPow(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Fast exponentiation for integers; negative exponent -> 0
	cg.generateExpressionToReg(args[0], "rbx") // base
	cg.generateExpressionToReg(args[1], "rcx") // exp
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    js 3f\n")
	cg.textSection.WriteString("1:  testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 2f\n")
	cg.textSection.WriteString("    testb $1, %cl\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("    imulq %rbx, %rax\n")
	cg.textSection.WriteString("4:  shrq $1, %rcx\n")
	cg.textSection.WriteString("    imulq %rbx, %rbx\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  jmp 5f\n")
	cg.textSection.WriteString("3:  xorq %rax, %rax\n")
	cg.textSection.WriteString("5:\n")
}

func generateMathFloor(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// For integers, floor is just the value itself
	cg.generateExpressionToReg(args[0], "rax")
}

func generateMathCeil(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// For integers, ceil is just the value itself
	cg.generateExpressionToReg(args[0], "rax")
}

func generateMathRound(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// For integers, round is just the value itself
	cg.generateExpressionToReg(args[0], "rax")
}

func generateMathGcd(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Greatest Common Divisor using Euclidean algorithm
	labelLoop := cg.getLabel("gcd_loop")
	labelEnd := cg.getLabel("gcd_end")

	// Load first argument into rax
	cg.generateExpressionToReg(args[0], "rax")
	// Load second argument into rcx
	cg.generateExpressionToReg(args[1], "rcx")

	// gcd(a,b):
	//   while b != 0:
	//     temp = b
	//     b = a % b
	//     a = temp
	//   return a
	// rax = a, rcx = b

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", labelEnd))

	// temp = b (save in r8)
	cg.textSection.WriteString("    movq %rcx, %r8\n")
	// b = a % b
	cg.textSection.WriteString("    xorq %rdx, %rdx\n") // clear rdx for division
	cg.textSection.WriteString("    divq %rcx\n")       // rax = a / b, rdx = a % b
	cg.textSection.WriteString("    movq %rdx, %rcx\n") // b = remainder
	// a = temp
	cg.textSection.WriteString("    movq %r8, %rax\n") // a = old b

	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
	// rax contains the GCD
}

func generateMathLcm(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// LCM(a,b) = (a * b) / GCD(a,b)
	// To avoid overflow, compute as: a / gcd(a,b) * b

	labelLoop := cg.getLabel("lcm_gcd_loop")
	labelEnd := cg.getLabel("lcm_gcd_end")

	// Load arguments
	cg.generateExpressionToReg(args[0], "r12") // save a in r12
	cg.generateExpressionToReg(args[1], "r13") // save b in r13

	// Compute GCD first
	cg.textSection.WriteString("    movq %r12, %rax\n") // a
	cg.textSection.WriteString("    movq %r13, %rcx\n") // b

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", labelEnd))
	cg.textSection.WriteString("    movq %rcx, %r8\n")  // temp = b
	cg.textSection.WriteString("    xorq %rdx, %rdx\n") // clear for division
	cg.textSection.WriteString("    divq %rcx\n")       // rax = a/b, rdx = a%b
	cg.textSection.WriteString("    movq %rdx, %rcx\n") // b = remainder
	cg.textSection.WriteString("    movq %r8, %rax\n")  // a = old b
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", labelLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", labelEnd))
	// rax = GCD

	// LCM = (a / GCD) * b
	cg.textSection.WriteString("    movq %rax, %rcx\n")  // rcx = GCD
	cg.textSection.WriteString("    movq %r12, %rax\n")  // rax = original a
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")  // clear for division
	cg.textSection.WriteString("    divq %rcx\n")        // rax = a / GCD
	cg.textSection.WriteString("    imulq %r13, %rax\n") // rax = (a / GCD) * b
}

func generateStringLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	switch v := args[0].(type) {
	case *StringLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", len(v.Value)))
	case *Identifier:
		if l, ok := cg.stringLengths[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", l))
			return
		}
		cg.textSection.WriteString("    movq $0, %rax\n")
	default:
		cg.textSection.WriteString("    movq $0, %rax\n")
	}
}

func generateStringConcat(cg *CodeGenerator, args []ASTNode) {
	// concat(a, b): allocate new buffer [len(a)+len(b)+1], copy a then b, NUL-terminate, return ptr in rax
	if len(args) < 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Resolve first string pointer -> r10, len -> r8
	resolvePtrLen := func(expr ASTNode, ptrReg, lenReg string, loopTag string) {
		switch v := expr.(type) {
		case *StringLiteral:
			lbl, _ := emitStringLiteral(cg, v.Value)
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%%s\n", lbl, ptrReg))
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%%s\n", len(v.Value), lenReg))
		case *Identifier:
			if varInfo, ok := cg.variables[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%%s\n", varInfo.Offset, ptrReg))
				if l, ok := cg.stringLengths[v.Name]; ok {
					cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%%s\n", l, lenReg))
				} else {
					// compute length at runtime
					lLoop := cg.getLabel(loopTag + "_loop")
					lEnd := cg.getLabel(loopTag + "_end")
					cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rbx\n", ptrReg))
					cg.textSection.WriteString(fmt.Sprintf("    movq $0, %%%s\n", lenReg))
					cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
					cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
					cg.textSection.WriteString("    testq %rax, %rax\n")
					cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
					cg.textSection.WriteString(fmt.Sprintf("    inc %%%s\n", lenReg))
					cg.textSection.WriteString("    inc %rbx\n")
					cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
					cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
				}
			} else {
				cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", ptrReg, ptrReg))
				cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", lenReg, lenReg))
			}
		default:
			cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", ptrReg, ptrReg))
			cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", lenReg, lenReg))
		}
	}
	resolvePtrLen(args[0], "r10", "r8", "s1len")
	resolvePtrLen(args[1], "r11", "r9", "s2len")
	// Save string pointers and lengths before mmap clobbers them
	// Use %r13-%r15 and %r12 since %r11 gets clobbered by syscall
	cg.textSection.WriteString("    movq %r10, %r13\n") // save s1 ptr
	cg.textSection.WriteString("    movq %r8, %r12\n")  // save s1 len
	cg.textSection.WriteString("    movq %r11, %r14\n") // save s2 ptr
	cg.textSection.WriteString("    movq %r9, %r15\n")  // save s2 len
	// total len in rdx = r8 + r9
	cg.textSection.WriteString("    movq %r8, %rdx\n")
	cg.textSection.WriteString("    addq %r9, %rdx\n")
	// size = total+1 in rsi
	cg.textSection.WriteString("    movq %rdx, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	// mmap(size)
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// Check mmap error
	lErr := cg.getLabel("cat_mmap_err")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))
	// dest base in rax, save to rbx
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	// copy first string
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq %r13, %rsi\n") // restore s1 ptr from r13
	cg.textSection.WriteString("    movq %r12, %rcx\n") // restore s1 len from r12
	l1 := cg.getLabel("cat_cp1")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", l1))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", l1))
	cg.textSection.WriteString("1:\n")
	// copy second string
	cg.textSection.WriteString("    movq %r14, %rsi\n") // restore s2 ptr from r14
	cg.textSection.WriteString("    movq %r15, %rcx\n") // restore s2 len from r15
	l2 := cg.getLabel("cat_cp2")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", l2))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 2f\n")
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", l2))
	cg.textSection.WriteString("2:\n")
	// NUL terminate
	cg.textSection.WriteString("    movb $0, (%rdi)\n")
	// return base ptr
	cg.textSection.WriteString("    movq %rbx, %rax\n")
	// Jump past error handler
	lEnd := cg.getLabel("cat_end")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	// On mmap error, return NULL
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

func generateStringCompare(cg *CodeGenerator, args []ASTNode) {
	// compare(a,b): returns 0 if equal, <0 if a<b, >0 if a>b (lexicographic)
	if len(args) != 2 {
		return
	}
	// Resolve pointers to the two strings into rbx and rcx
	// Supports string literals and identifiers
	resolveStrPtr := func(expr ASTNode, targetReg string) {
		switch v := expr.(type) {
		case *StringLiteral:
			lbl, _ := emitStringLiteral(cg, v.Value)
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%%s\n", lbl, targetReg))
		case *Identifier:
			if varInfo, ok := cg.variables[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%%s\n", varInfo.Offset, targetReg))
			} else {
				cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", targetReg, targetReg))
			}
		default:
			cg.textSection.WriteString(fmt.Sprintf("    xorq %%%s, %%%s\n", targetReg, targetReg))
		}
	}
	resolveStrPtr(args[0], "rbx")
	resolveStrPtr(args[1], "rcx")
	lLoop := cg.getLabel("strcmp_loop")
	lEnd := cg.getLabel("strcmp_end")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
	cg.textSection.WriteString("    movzbq (%rcx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lEnd))
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    inc %rbx\n")
	cg.textSection.WriteString("    inc %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    subq %rdx, %rax\n")
	// result in rax
}

func generateStringCopy(cg *CodeGenerator, args []ASTNode) {
	// copy(src): allocate new buffer [len(src)+1], copy, NUL-terminate, return ptr
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Resolve src pointer -> r10, len -> r8
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", l))
			} else {
				lLoop := cg.getLabel("cpy_len_loop")
				lEnd := cg.getLabel("cpy_len_end")
				cg.textSection.WriteString("    movq %r10, %rbx\n")
				cg.textSection.WriteString("    movq $0, %r8\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    inc %r8\n")
				cg.textSection.WriteString("    inc %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n    xorq %r8, %r8\n")
		}
	default:
		cg.textSection.WriteString("    xorq %r10, %r10\n    xorq %r8, %r8\n")
	}
	// Save source ptr and length before mmap clobbers registers
	// Use %r13, %r14 instead of %r11, %r12 because syscall clobbers %r11
	cg.textSection.WriteString("    movq %r10, %r13\n") // save src ptr to r13
	cg.textSection.WriteString("    movq %r8, %r14\n")  // save src len to r14
	// size = len+1 in rsi
	cg.textSection.WriteString("    movq %r8, %rsi\n    addq $1, %rsi\n")
	// mmap(size)
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// Check mmap error: if rax is negative, it's an error (Linux returns -errno)
	lErr := cg.getLabel("copy_mmap_err")
	lEnd := cg.getLabel("copy_end")
	cg.textSection.WriteString(fmt.Sprintf("    testq %%rax, %%rax\n"))
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))
	// dest base in rax, save to rbx
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	// copy: rdi=dest, rsi=src, rcx=len
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq %r13, %rsi\n") // restore src from r13
	cg.textSection.WriteString("    movq %r14, %rcx\n") // restore len from r14
	l1 := cg.getLabel("cpy_cp")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", l1))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", l1))
	cg.textSection.WriteString("3:\n")
	cg.textSection.WriteString("    movb $0, (%rdi)\n")
	cg.textSection.WriteString("    movq %rbx, %rax\n")
	// Jump past error handler
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	// On mmap error, return NULL (rax = 0)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// indexOf(s, chOrSubstr): if second arg is string literal of length 1, treat as char
func generateStringIndexOf(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Resolve haystack in rbx
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rbx\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %rbx, %rbx\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rbx, %rbx\n")
	}
	// Resolve needle in al (single byte) or pointer in rcx
	isChar := false
	if lit, ok := args[1].(*StringLiteral); ok && len(lit.Value) == 1 {
		isChar = true
		cg.textSection.WriteString(fmt.Sprintf("    movb $%d, %%al\n", byte(lit.Value[0])))
	} else {
		switch v := args[1].(type) {
		case *StringLiteral:
			lbl, _ := emitStringLiteral(cg, v.Value)
			cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rcx\n", lbl))
		case *Identifier:
			if varInfo, ok := cg.variables[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rcx\n", varInfo.Offset))
			} else {
				cg.textSection.WriteString("    xorq %rcx, %rcx\n")
			}
		default:
			cg.textSection.WriteString("    xorq %rcx, %rcx\n")
		}
	}
	if isChar {
		lLoop := cg.getLabel("idx_loop")
		lNotFound := cg.getLabel("idx_nf")
		lEnd := cg.getLabel("idx_end")
		lDone := cg.getLabel("idx_done")
		cg.textSection.WriteString("    movq $0, %r10\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbq (%rbx), %rdx\n")
		cg.textSection.WriteString("    testq %rdx, %rdx\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lNotFound))
		cg.textSection.WriteString("    cmpb %al, %dl\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    inc %r10\n")
		cg.textSection.WriteString("    inc %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNotFound))
		cg.textSection.WriteString("    movq $-1, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lDone))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
		cg.textSection.WriteString("    movq %r10, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lDone))
		return
	}
	// substring search (naive)
	lOuter := cg.getLabel("sub_outer")
	lInner := cg.getLabel("sub_inner")
	lNext := cg.getLabel("sub_next")
	lFound := cg.getLabel("sub_found")
	lNF := cg.getLabel("sub_nf")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lOuter))
	cg.textSection.WriteString("    movzbq (%rbx), %r8\n")
	cg.textSection.WriteString("    testq %r8, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lNF))
	cg.textSection.WriteString("    movq %rbx, %r10\n")
	cg.textSection.WriteString("    movq %rcx, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lInner))
	cg.textSection.WriteString("    movzbq (%r11), %r8\n")
	cg.textSection.WriteString("    testq %r8, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lFound))
	cg.textSection.WriteString("    movzbq (%r10), %r9\n")
	cg.textSection.WriteString("    cmpq %r8, %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lNext))
	cg.textSection.WriteString("    inc %r10\n")
	cg.textSection.WriteString("    inc %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lInner))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNext))
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    inc %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lOuter))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lFound))
	// rax holds starting index
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_end\n", lOuter))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNF))
	cg.textSection.WriteString("    movq $-1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_end:\n", lOuter))
}

func generateStringContains(cg *CodeGenerator, args []ASTNode) {
	// contains(s, sub) -> indexOf >= 0
	generateStringIndexOf(cg, args)
	// rax has index or -1
	lTrue := cg.getLabel("contains_true")
	lEnd := cg.getLabel("contains_end")
	cg.textSection.WriteString("    cmpq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lTrue))
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrue))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// ============================================================================
// Net module (Linux syscalls)
// ============================================================================

// socket(domain, type, protocol)
func generateNetSocket(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rdx")
	cg.textSection.WriteString("    movq $41, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// connect_ipv4(fd, ip_u32_host, port_host)
func generateNetConnectIPv4(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	cg.generateExpressionToReg(args[1], "rsi") // ip host order
	cg.generateExpressionToReg(args[2], "rdx") // port host order
	cg.textSection.WriteString("    subq $32, %rsp\n")
	cg.textSection.WriteString("    movw $2, (%rsp)\n") // AF_INET
	cg.textSection.WriteString("    movq %rdx, %rcx\n")
	cg.textSection.WriteString("    rorw $8, %cx\n")
	cg.textSection.WriteString("    movw %cx, 2(%rsp)\n") // port network order
	cg.textSection.WriteString("    movq %rsi, %r8\n")
	cg.textSection.WriteString("    movl %r8d, %ecx\n")
	cg.textSection.WriteString("    bswap %ecx\n")
	cg.textSection.WriteString("    movl %ecx, 4(%rsp)\n") // ip network order
	cg.textSection.WriteString("    movq $0, 8(%rsp)\n")
	cg.textSection.WriteString("    movq $0, 16(%rsp)\n")
	cg.textSection.WriteString("    movq %rsp, %rsi\n")
	cg.textSection.WriteString("    movq $16, %rdx\n")
	cg.textSection.WriteString("    movq $42, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $32, %rsp\n")
}

// send(fd, buf, len) -> bytes written
func generateNetSend(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rdx")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// recv(fd, buf, len) -> bytes read
func generateNetRecv(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rdx")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// close(fd)
func generateNetClose(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq $3, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// ============================================================================
// UDP Networking Support
// ============================================================================

// bind_ipv4(fd, ip_u32_host, port_host) -> 0 on success, negative on error
func generateNetBindIPv4(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	cg.generateExpressionToReg(args[1], "rsi") // ip host order
	cg.generateExpressionToReg(args[2], "rdx") // port host order

	// Build sockaddr_in on stack
	cg.textSection.WriteString("    subq $32, %rsp\n")
	cg.textSection.WriteString("    movw $2, (%rsp)\n") // AF_INET = 2
	// Convert port to network byte order (big-endian)
	cg.textSection.WriteString("    movq %rdx, %rcx\n")
	cg.textSection.WriteString("    rorw $8, %cx\n")
	cg.textSection.WriteString("    movw %cx, 2(%rsp)\n") // port network order
	// Convert IP to network byte order
	cg.textSection.WriteString("    movq %rsi, %r8\n")
	cg.textSection.WriteString("    movl %r8d, %ecx\n")
	cg.textSection.WriteString("    bswap %ecx\n")
	cg.textSection.WriteString("    movl %ecx, 4(%rsp)\n") // ip network order
	// Zero padding
	cg.textSection.WriteString("    movq $0, 8(%rsp)\n")
	cg.textSection.WriteString("    movq $0, 16(%rsp)\n")
	// bind syscall: fd in rdi, addr in rsi, addrlen in rdx
	cg.textSection.WriteString("    movq %rsp, %rsi\n")
	cg.textSection.WriteString("    movq $16, %rdx\n")
	cg.textSection.WriteString("    movq $49, %rax\n") // syscall 49 = bind
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $32, %rsp\n")
}

// sendto_ipv4(fd, buf_ptr, buf_len, dest_ip_u32, dest_port) -> bytes sent
func generateNetSendtoIPv4(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 5 {
		return
	}
	// Save fd in r12
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    pushq %rdi\n") // save fd

	// Build sockaddr_in for destination on stack
	cg.textSection.WriteString("    subq $32, %rsp\n")
	cg.textSection.WriteString("    movw $2, (%rsp)\n") // AF_INET = 2

	// Get dest_port and convert to network order
	cg.generateExpressionToReg(args[4], "rcx") // dest_port
	cg.textSection.WriteString("    rorw $8, %cx\n")
	cg.textSection.WriteString("    movw %cx, 2(%rsp)\n")

	// Get dest_ip and convert to network order
	cg.generateExpressionToReg(args[3], "r8") // dest_ip
	cg.textSection.WriteString("    movl %r8d, %ecx\n")
	cg.textSection.WriteString("    bswap %ecx\n")
	cg.textSection.WriteString("    movl %ecx, 4(%rsp)\n")

	// Zero padding
	cg.textSection.WriteString("    movq $0, 8(%rsp)\n")
	cg.textSection.WriteString("    movq $0, 16(%rsp)\n")

	// Now set up sendto syscall
	// rdi = fd, rsi = buf, rdx = len, r10 = flags (0), r8 = addr, r9 = addrlen
	cg.textSection.WriteString("    movq 32(%rsp), %rdi\n") // restore fd from stack
	cg.generateExpressionToReg(args[1], "rsi")              // buf_ptr
	cg.generateExpressionToReg(args[2], "rdx")              // buf_len
	cg.textSection.WriteString("    xorq %r10, %r10\n")     // flags = 0
	cg.textSection.WriteString("    movq %rsp, %r8\n")      // addr ptr
	cg.textSection.WriteString("    movq $16, %r9\n")       // addrlen
	cg.textSection.WriteString("    movq $44, %rax\n")      // syscall 44 = sendto
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $40, %rsp\n") // clean up sockaddr + saved fd
}

// recvfrom(fd, buf_ptr, buf_len) -> bytes received
// Note: This simplified version doesn't return sender info
func generateNetRecvfrom(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")          // fd
	cg.generateExpressionToReg(args[1], "rsi")          // buf_ptr
	cg.generateExpressionToReg(args[2], "rdx")          // buf_len
	cg.textSection.WriteString("    xorq %r10, %r10\n") // flags = 0
	cg.textSection.WriteString("    xorq %r8, %r8\n")   // addr = NULL (don't care about sender)
	cg.textSection.WriteString("    xorq %r9, %r9\n")   // addrlen = NULL
	cg.textSection.WriteString("    movq $45, %rax\n")  // syscall 45 = recvfrom
	cg.textSection.WriteString("    syscall\n")
}

// ============================================================================
// IPv6 Support
// ============================================================================
// sockaddr_in6 structure (28 bytes):
//   offset 0:  sin6_family (2 bytes) = AF_INET6 = 10
//   offset 2:  sin6_port (2 bytes, network byte order)
//   offset 4:  sin6_flowinfo (4 bytes)
//   offset 8:  sin6_addr (16 bytes)
//   offset 24: sin6_scope_id (4 bytes)

const AF_INET6 = 10
const sockaddrIn6Size = 28

// generateNetConnectIPv6 connects to an IPv6 address
// Args: ipv6_addr_ptr (16 bytes), port, socket_type (1=TCP, 2=UDP)
// Returns: socket fd or negative error
func generateNetConnectIPv6(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    movq $-1, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "r12") // ipv6 addr ptr
	cg.generateExpressionToReg(args[1], "r13") // port
	cg.generateExpressionToReg(args[2], "r14") // socket type

	// Create socket: socket(AF_INET6, type, 0)
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdi\n", AF_INET6))
	cg.textSection.WriteString("    movq %r14, %rsi\n") // SOCK_STREAM=1, SOCK_DGRAM=2
	cg.textSection.WriteString("    xorq %rdx, %rdx\n") // protocol = 0
	cg.textSection.WriteString("    movq $41, %rax\n")  // socket syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    js 99f\n")          // return error if negative
	cg.textSection.WriteString("    movq %rax, %r15\n") // save socket fd

	// Allocate sockaddr_in6 on stack
	cg.textSection.WriteString(fmt.Sprintf("    subq $%d, %%rsp\n", sockaddrIn6Size))
	// Zero it out
	cg.textSection.WriteString("    movq %rsp, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rcx\n", sockaddrIn6Size/8))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosq\n")

	// Fill sockaddr_in6
	cg.textSection.WriteString(fmt.Sprintf("    movw $%d, (%%rsp)\n", AF_INET6)) // sin6_family
	// Port in network byte order (big endian)
	cg.textSection.WriteString("    movq %r13, %rax\n")
	cg.textSection.WriteString("    xchgb %al, %ah\n")    // swap bytes
	cg.textSection.WriteString("    movw %ax, 2(%rsp)\n") // sin6_port
	// Copy IPv6 address (16 bytes)
	cg.textSection.WriteString("    movq (%r12), %rax\n")
	cg.textSection.WriteString("    movq %rax, 8(%rsp)\n")
	cg.textSection.WriteString("    movq 8(%r12), %rax\n")
	cg.textSection.WriteString("    movq %rax, 16(%rsp)\n")

	// Connect: connect(fd, addr, addrlen)
	cg.textSection.WriteString("    movq %r15, %rdi\n")
	cg.textSection.WriteString("    movq %rsp, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", sockaddrIn6Size))
	cg.textSection.WriteString("    movq $42, %rax\n") // connect syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsp\n", sockaddrIn6Size))
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    js 98f\n")
	cg.textSection.WriteString("    movq %r15, %rax\n") // return socket fd
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString("98:\n")
	// Close socket on connect error
	cg.textSection.WriteString("    movq %r15, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rax\n") // close
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq $-1, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateNetBindIPv6 binds a socket to an IPv6 address
// Args: socket_fd, ipv6_addr_ptr (16 bytes or NULL for any), port
// Returns: 0 on success, negative on error
func generateNetBindIPv6(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    movq $-1, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "r12") // socket fd
	cg.generateExpressionToReg(args[1], "r13") // ipv6 addr ptr (or 0 for any)
	cg.generateExpressionToReg(args[2], "r14") // port

	lblCopyAddr := cg.getLabel("bind6_copy")
	lblBind := cg.getLabel("bind6_do")

	// Allocate sockaddr_in6 on stack
	cg.textSection.WriteString(fmt.Sprintf("    subq $%d, %%rsp\n", sockaddrIn6Size))
	// Zero it out
	cg.textSection.WriteString("    movq %rsp, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rcx\n", sockaddrIn6Size/8))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosq\n")

	cg.textSection.WriteString(fmt.Sprintf("    movw $%d, (%%rsp)\n", AF_INET6))
	cg.textSection.WriteString("    movq %r14, %rax\n")
	cg.textSection.WriteString("    xchgb %al, %ah\n")
	cg.textSection.WriteString("    movw %ax, 2(%rsp)\n")

	// If addr ptr is not NULL, copy it
	cg.textSection.WriteString("    testq %r13, %r13\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblBind))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopyAddr))
	cg.textSection.WriteString("    movq (%r13), %rax\n")
	cg.textSection.WriteString("    movq %rax, 8(%rsp)\n")
	cg.textSection.WriteString("    movq 8(%r13), %rax\n")
	cg.textSection.WriteString("    movq %rax, 16(%rsp)\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblBind))
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq %rsp, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", sockaddrIn6Size))
	cg.textSection.WriteString("    movq $49, %rax\n") // bind syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsp\n", sockaddrIn6Size))
}

// generateNetSendtoIPv6 sends data to an IPv6 address (for UDP)
// Args: socket_fd, buf_ptr, buf_len, ipv6_addr_ptr, port
// Returns: bytes sent or negative on error
func generateNetSendtoIPv6(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 5 {
		cg.textSection.WriteString("    movq $-1, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "r12") // socket fd
	cg.generateExpressionToReg(args[1], "r13") // buf ptr
	cg.generateExpressionToReg(args[2], "r14") // buf len
	cg.generateExpressionToReg(args[3], "r15") // ipv6 addr ptr
	cg.generateExpressionToReg(args[4], "rbx") // port

	// Allocate sockaddr_in6 on stack
	cg.textSection.WriteString(fmt.Sprintf("    subq $%d, %%rsp\n", sockaddrIn6Size))
	cg.textSection.WriteString("    movq %rsp, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rcx\n", sockaddrIn6Size/8))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosq\n")

	cg.textSection.WriteString(fmt.Sprintf("    movw $%d, (%%rsp)\n", AF_INET6))
	cg.textSection.WriteString("    movq %rbx, %rax\n")
	cg.textSection.WriteString("    xchgb %al, %ah\n")
	cg.textSection.WriteString("    movw %ax, 2(%rsp)\n")
	cg.textSection.WriteString("    movq (%r15), %rax\n")
	cg.textSection.WriteString("    movq %rax, 8(%rsp)\n")
	cg.textSection.WriteString("    movq 8(%r15), %rax\n")
	cg.textSection.WriteString("    movq %rax, 16(%rsp)\n")

	// sendto(fd, buf, len, flags, addr, addrlen)
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq %r13, %rsi\n")
	cg.textSection.WriteString("    movq %r14, %rdx\n")
	cg.textSection.WriteString("    xorq %r10, %r10\n") // flags = 0
	cg.textSection.WriteString("    movq %rsp, %r8\n")  // addr
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", sockaddrIn6Size))
	cg.textSection.WriteString("    movq $44, %rax\n") // sendto syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsp\n", sockaddrIn6Size))
}

// ============================================================================
// DNS Resolution (simplified)
// ============================================================================
// These functions parse /etc/hosts for simple resolution
// For real DNS queries, a full resolver would be needed

// generateNetResolve resolves a hostname to IPv4 address via /etc/hosts
// Args: hostname_ptr, out_ipv4_ptr (4 bytes)
// Returns: 1 on success, 0 on failure
func generateNetResolve(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "r12") // hostname ptr
	cg.generateExpressionToReg(args[1], "r13") // out ipv4 ptr

	lblRead := cg.getLabel("dns_read")
	lblLine := cg.getLabel("dns_line")
	lblCompare := cg.getLabel("dns_cmp")
	lblFound := cg.getLabel("dns_found")
	lblNext := cg.getLabel("dns_next")
	lblNotFound := cg.getLabel("dns_nf")
	lblClose := cg.getLabel("dns_close")

	// Open /etc/hosts
	hostsPath, _ := emitStringLiteral(cg, "/etc/hosts")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rdi\n", hostsPath))
	cg.textSection.WriteString("    xorq %rsi, %rsi\n") // O_RDONLY
	cg.textSection.WriteString("    movq $2, %rax\n")   // open
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lblNotFound))
	cg.textSection.WriteString("    movq %rax, %r14\n") // save fd

	// Allocate buffer on stack (512 bytes)
	cg.textSection.WriteString("    subq $512, %rsp\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRead))
	cg.textSection.WriteString("    movq %r14, %rdi\n")
	cg.textSection.WriteString("    movq %rsp, %rsi\n")
	cg.textSection.WriteString("    movq $512, %rdx\n")
	cg.textSection.WriteString("    movq $0, %rax\n") // read
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jle %s\n", lblClose))
	cg.textSection.WriteString("    movq %rax, %r15\n") // bytes read

	// Parse buffer (simplified: look for hostname match)
	cg.textSection.WriteString("    movq %rsp, %rcx\n") // current pos
	cg.textSection.WriteString("    addq %r15, %rsp\n") // temp end marker

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLine))
	// Skip leading whitespace and comments
	cg.textSection.WriteString("    cmpq %rsp, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblClose))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $35, %al\n") // '#' comment
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblNext))
	cg.textSection.WriteString("    cmpb $10, %al\n") // newline
	cg.textSection.WriteString(fmt.Sprintf("    je %s_skip\n", lblLine))
	cg.textSection.WriteString("    cmpb $32, %al\n") // space
	cg.textSection.WriteString(fmt.Sprintf("    je %s_skip\n", lblLine))
	cg.textSection.WriteString("    cmpb $9, %al\n") // tab
	cg.textSection.WriteString(fmt.Sprintf("    je %s_skip\n", lblLine))

	// This is an IP address, save position
	cg.textSection.WriteString("    movq %rcx, %rbx\n") // IP start
	// Skip to whitespace (end of IP)
	cg.textSection.WriteString(fmt.Sprintf("%s_skipip:\n", lblLine))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    cmpq %rsp, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblClose))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s_ws\n", lblLine))
	cg.textSection.WriteString("    cmpb $9, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s_skipip\n", lblLine))
	cg.textSection.WriteString(fmt.Sprintf("%s_ws:\n", lblLine))
	// Skip whitespace to hostname
	cg.textSection.WriteString(fmt.Sprintf("%s_skiphostws:\n", lblLine))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    cmpq %rsp, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblClose))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s_skiphostws\n", lblLine))
	cg.textSection.WriteString("    cmpb $9, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s_skiphostws\n", lblLine))

	// Compare hostname
	cg.textSection.WriteString("    movq %r12, %rdi\n") // target hostname
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCompare))
	cg.textSection.WriteString("    movzbq (%rdi), %rax\n")
	cg.textSection.WriteString("    movzbq (%rcx), %rdx\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s_check\n", lblCompare))
	cg.textSection.WriteString("    cmpb %al, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCompare))
	cg.textSection.WriteString(fmt.Sprintf("%s_check:\n", lblCompare))
	// Check that host entry also ends (whitespace or newline)
	cg.textSection.WriteString("    cmpb $32, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString("    cmpb $9, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString("    cmpb $10, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString("    cmpb $13, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString("    testq %rdx, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblFound))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	// Skip to end of line
	cg.textSection.WriteString(fmt.Sprintf("%s_eol:\n", lblNext))
	cg.textSection.WriteString("    cmpq %rsp, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblClose))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s_eol\n", lblNext))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLine))

	cg.textSection.WriteString(fmt.Sprintf("%s_skip:\n", lblLine))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLine))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	// Parse IP from rbx and store to r13
	// Simple IPv4 parser: a.b.c.d
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    xorq %r8, %r8\n") // octet index
	cg.textSection.WriteString("    xorq %r9, %r9\n") // current octet value
	cg.textSection.WriteString(fmt.Sprintf("%s_ip:\n", lblFound))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $46, %al\n") // '.'
	cg.textSection.WriteString(fmt.Sprintf("    je %s_dot\n", lblFound))
	cg.textSection.WriteString("    cmpb $48, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jl %s_endip\n", lblFound))
	cg.textSection.WriteString("    cmpb $57, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jg %s_endip\n", lblFound))
	cg.textSection.WriteString("    subb $48, %al\n")
	cg.textSection.WriteString("    imulq $10, %r9\n")
	cg.textSection.WriteString("    addq %rax, %r9\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_ip\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("%s_dot:\n", lblFound))
	cg.textSection.WriteString("    movb %r9b, (%r13,%r8)\n")
	cg.textSection.WriteString("    incq %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_ip\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("%s_endip:\n", lblFound))
	cg.textSection.WriteString("    movb %r9b, (%r13,%r8)\n")
	cg.textSection.WriteString("    movq $1, %r8\n") // success flag
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_done\n", lblClose))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblClose))
	cg.textSection.WriteString("    subq %r15, %rsp\n") // restore stack
	cg.textSection.WriteString(fmt.Sprintf("%s_done:\n", lblClose))
	cg.textSection.WriteString("    addq $512, %rsp\n")
	cg.textSection.WriteString("    movq %r14, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rax\n") // close
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq %r8, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateNetResolveIPv6 is a stub for IPv6 resolution
// For simplicity, returns 0 (not found) - full implementation would need proper DNS
func generateNetResolveIPv6(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// IPv6 resolution is complex - return not found for now
	// A real implementation would parse /etc/hosts for IPv6 or do DNS AAAA queries
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// ============================================================================
// HTTP module - minimal GET over an existing connected socket
// ============================================================================

// get(fd, host_ptr, host_len, path_ptr, path_len, buf_ptr, buf_len) -> bytes read
func generateHTTPGetSimple(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 7 {
		return
	}
	// fd in r12 for reuse
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq %rdi, %r12\n")

	// literals
	lblGet, _ := emitStringLiteral(cg, "GET ")
	lblMid, _ := emitStringLiteral(cg, " HTTP/1.0\r\nHost: ")
	lblEnd, _ := emitStringLiteral(cg, "\r\nConnection: close\r\n\r\n")

	writeLiteral := func(label string, length int) {
		cg.textSection.WriteString(fmt.Sprintf("    movq %%r12, %%rdi\n"))
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", label))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
		cg.textSection.WriteString("    movq $1, %rax\n")
		cg.textSection.WriteString("    syscall\n")
	}

	writeLiteral(lblGet, 4)

	// write path
	cg.generateExpressionToReg(args[3], "rsi") // path ptr
	cg.generateExpressionToReg(args[4], "rdx") // path len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	writeLiteral(lblMid, 17)

	// write host
	cg.generateExpressionToReg(args[1], "rsi") // host ptr
	cg.generateExpressionToReg(args[2], "rdx") // host len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	writeLiteral(lblEnd, 24)

	// read response into caller buffer
	cg.generateExpressionToReg(args[5], "rsi") // buf ptr
	cg.generateExpressionToReg(args[6], "rdx") // buf len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// post(fd, host_ptr, host_len, path_ptr, path_len, body_ptr, body_len, buf_ptr, buf_len) -> bytes read
func generateHTTPPostSimple(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 9 {
		return
	}
	// fd in r12 for reuse
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq %rdi, %r12\n")
	// body_len in r13 for Content-Length header
	cg.generateExpressionToReg(args[6], "rdi")
	cg.textSection.WriteString("    movq %rdi, %r13\n")

	// literals
	lblPost, _ := emitStringLiteral(cg, "POST ")
	lblMid, _ := emitStringLiteral(cg, " HTTP/1.0\r\nHost: ")
	lblContentType, _ := emitStringLiteral(cg, "\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: ")
	lblEnd, _ := emitStringLiteral(cg, "\r\nConnection: close\r\n\r\n")

	writeLiteral := func(label string, length int) {
		cg.textSection.WriteString("    movq %r12, %rdi\n")
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", label))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rdx\n", length))
		cg.textSection.WriteString("    movq $1, %rax\n")
		cg.textSection.WriteString("    syscall\n")
	}

	writeLiteral(lblPost, 5)

	// write path
	cg.generateExpressionToReg(args[3], "rsi") // path ptr
	cg.generateExpressionToReg(args[4], "rdx") // path len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	writeLiteral(lblMid, 17)

	// write host
	cg.generateExpressionToReg(args[1], "rsi") // host ptr
	cg.generateExpressionToReg(args[2], "rdx") // host len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	writeLiteral(lblContentType, 62)

	// Write Content-Length as decimal string
	// Convert r13 (body_len) to decimal and write
	// Use stack buffer for conversion
	cg.textSection.WriteString("    subq $32, %rsp\n")
	cg.textSection.WriteString("    movq %r13, %rax\n")     // number to convert
	cg.textSection.WriteString("    leaq 31(%rsp), %rdi\n") // end of buffer
	cg.textSection.WriteString("    movb $0, (%rdi)\n")     // null terminator
	cg.textSection.WriteString("    movq %rdi, %r14\n")     // save end position
	lblConvLoop := cg.getLabel("post_conv")
	lblConvDone := cg.getLabel("post_conv_done")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblConvLoop))
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    movq $10, %rcx\n")
	cg.textSection.WriteString("    divq %rcx\n")        // rax = quotient, rdx = remainder
	cg.textSection.WriteString("    addb $48, %dl\n")    // convert to ASCII
	cg.textSection.WriteString("    decq %rdi\n")        // move buffer position back
	cg.textSection.WriteString("    movb %dl, (%rdi)\n") // store digit
	cg.textSection.WriteString("    testq %rax, %rax\n") // check if done
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblConvLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblConvDone))
	// Now rdi points to start of number string, r14 points to end
	cg.textSection.WriteString("    movq %rdi, %rsi\n") // source = start of number
	cg.textSection.WriteString("    movq %r14, %rdx\n") // compute length
	cg.textSection.WriteString("    subq %rsi, %rdx\n") // rdx = length
	cg.textSection.WriteString("    movq %r12, %rdi\n") // fd
	cg.textSection.WriteString("    movq $1, %rax\n")   // write syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $32, %rsp\n")

	writeLiteral(lblEnd, 24)

	// write body
	cg.generateExpressionToReg(args[5], "rsi")          // body ptr
	cg.textSection.WriteString("    movq %r13, %rdx\n") // body len from r13
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	// read response into caller buffer
	cg.generateExpressionToReg(args[7], "rsi") // buf ptr
	cg.generateExpressionToReg(args[8], "rdx") // buf len
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// ============================================================================
// HTTP Response Parsing Functions
// ============================================================================

// generateHTTPParseStatus parses HTTP response status code
// Args: response_buffer, buffer_len -> returns status code (e.g., 200, 404) or 0 on error
func generateHTTPParseStatus(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // buffer ptr
	cg.generateExpressionToReg(args[1], "r12") // buffer len

	lblParse := cg.getLabel("http_parse")
	lblSkipSpace := cg.getLabel("http_skip")
	lblDigits := cg.getLabel("http_digits")
	lblDone := cg.getLabel("http_done")
	lblError := cg.getLabel("http_err")

	// HTTP response: "HTTP/1.x STATUS MESSAGE"
	// Skip past "HTTP/X.X " (9 chars minimum)
	cg.textSection.WriteString("    cmpq $12, %r12\n") // need at least 12 chars
	cg.textSection.WriteString(fmt.Sprintf("    jl %s\n", lblError))

	// Find first space after HTTP/X.X
	cg.textSection.WriteString("    movq %rbx, %rcx\n") // current pos
	cg.textSection.WriteString("    addq %r12, %rbx\n") // end pos
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblParse))
	cg.textSection.WriteString("    cmpq %rbx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblError))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $32, %al\n") // space?
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblSkipSpace))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblParse))

	// Skip spaces
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSkipSpace))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    cmpq %rbx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblError))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblSkipSpace))

	// Now parse 3-digit status code
	cg.textSection.WriteString("    xorq %r13, %r13\n") // result = 0
	cg.textSection.WriteString("    movq $3, %r14\n")   // digit count
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDigits))
	cg.textSection.WriteString("    testq %r14, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    cmpq %rbx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblError))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $48, %al\n") // < '0'?
	cg.textSection.WriteString(fmt.Sprintf("    jl %s\n", lblError))
	cg.textSection.WriteString("    cmpb $57, %al\n") // > '9'?
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblError))
	cg.textSection.WriteString("    subb $48, %al\n") // digit value
	cg.textSection.WriteString("    imulq $10, %r13\n")
	cg.textSection.WriteString("    addq %rax, %r13\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    decq %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblDigits))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movq %r13, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblError))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateHTTPGetHeader extracts a header value from HTTP response
// Args: response_buffer, buffer_len, header_name, out_value_ptr
// Returns: length of value or 0 if not found
func generateHTTPGetHeader(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 4 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // buffer ptr
	cg.generateExpressionToReg(args[1], "r12") // buffer len
	cg.generateExpressionToReg(args[2], "r13") // header name ptr
	cg.generateExpressionToReg(args[3], "r15") // out value ptr

	lblLoop := cg.getLabel("hdr_loop")
	lblLineStart := cg.getLabel("hdr_line")
	lblCompare := cg.getLabel("hdr_cmp")
	lblCopyValue := cg.getLabel("hdr_copy")
	lblNextLine := cg.getLabel("hdr_next")
	lblNotFound := cg.getLabel("hdr_nf")
	lblDone := cg.getLabel("hdr_done")

	// r14 = current pos, rbx = buffer end
	cg.textSection.WriteString("    movq %rbx, %r14\n")
	cg.textSection.WriteString("    addq %r12, %rbx\n")

	// Find first line (skip status line - find \n)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    cmpq %rbx, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblNotFound))
	cg.textSection.WriteString("    movzbq (%r14), %rax\n")
	cg.textSection.WriteString("    incq %r14\n")
	cg.textSection.WriteString("    cmpb $10, %al\n") // \n?
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblLoop))

	// Now we're at start of headers
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLineStart))
	cg.textSection.WriteString("    cmpq %rbx, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblNotFound))
	// Check for empty line (end of headers)
	cg.textSection.WriteString("    movzbq (%r14), %rax\n")
	cg.textSection.WriteString("    cmpb $13, %al\n") // \r?
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblNotFound))
	cg.textSection.WriteString("    cmpb $10, %al\n") // \n?
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblNotFound))

	// Compare header name (case-insensitive would be better, but simplified)
	cg.textSection.WriteString("    movq %r14, %rcx\n") // current line
	cg.textSection.WriteString("    movq %r13, %rdx\n") // header name
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCompare))
	cg.textSection.WriteString("    movzbq (%rdx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n") // end of header name?
	cg.textSection.WriteString(fmt.Sprintf("    jz %s_colon\n", lblCompare))
	cg.textSection.WriteString("    movzbq (%rcx), %rsi\n")
	// Simple lowercase conversion for comparison
	cg.textSection.WriteString("    cmpb $65, %al\n") // >= 'A'?
	cg.textSection.WriteString(fmt.Sprintf("    jl %s_cmp2\n", lblCompare))
	cg.textSection.WriteString("    cmpb $90, %al\n") // <= 'Z'?
	cg.textSection.WriteString(fmt.Sprintf("    jg %s_cmp2\n", lblCompare))
	cg.textSection.WriteString("    addb $32, %al\n") // to lowercase
	cg.textSection.WriteString(fmt.Sprintf("%s_cmp2:\n", lblCompare))
	cg.textSection.WriteString("    cmpb $65, %sil\n")
	cg.textSection.WriteString(fmt.Sprintf("    jl %s_cmp3\n", lblCompare))
	cg.textSection.WriteString("    cmpb $90, %sil\n")
	cg.textSection.WriteString(fmt.Sprintf("    jg %s_cmp3\n", lblCompare))
	cg.textSection.WriteString("    addb $32, %sil\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_cmp3:\n", lblCompare))
	cg.textSection.WriteString("    cmpb %al, %sil\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNextLine))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    incq %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCompare))

	// Check for colon after header name
	cg.textSection.WriteString(fmt.Sprintf("%s_colon:\n", lblCompare))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $58, %al\n") // ':'?
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNextLine))
	cg.textSection.WriteString("    incq %rcx\n")
	// Skip spaces after colon
	cg.textSection.WriteString(fmt.Sprintf("%s_skip:\n", lblCompare))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblCopyValue))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_skip\n", lblCompare))

	// Copy header value until \r or \n
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopyValue))
	cg.textSection.WriteString("    movq %r15, %rdi\n") // out ptr
	cg.textSection.WriteString("    xorq %r8, %r8\n")   // length counter
	cg.textSection.WriteString(fmt.Sprintf("%s_loop:\n", lblCopyValue))
	cg.textSection.WriteString("    cmpq %rbx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblDone))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $13, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblDone))
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblDone))
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString("    incq %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_loop\n", lblCopyValue))

	// Skip to next line
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNextLine))
	cg.textSection.WriteString(fmt.Sprintf("%s_find:\n", lblNextLine))
	cg.textSection.WriteString("    cmpq %rbx, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblNotFound))
	cg.textSection.WriteString("    movzbq (%r14), %rax\n")
	cg.textSection.WriteString("    incq %r14\n")
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s_find\n", lblNextLine))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLineStart))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // null terminate
	cg.textSection.WriteString("    movq %r8, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateHTTPGetBody returns pointer to body (after \r\n\r\n)
// Args: response_buffer, buffer_len -> returns pointer to body or 0 if not found
func generateHTTPGetBody(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // buffer ptr
	cg.generateExpressionToReg(args[1], "r12") // buffer len

	lblLoop := cg.getLabel("body_loop")
	lblCheck := cg.getLabel("body_check")
	lblFound := cg.getLabel("body_found")
	lblNotFound := cg.getLabel("body_nf")

	// Look for \r\n\r\n (0x0D 0x0A 0x0D 0x0A)
	cg.textSection.WriteString("    movq %rbx, %rcx\n") // current pos
	cg.textSection.WriteString("    addq %r12, %rbx\n") // end pos
	cg.textSection.WriteString("    subq $4, %rbx\n")   // need at least 4 bytes

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    cmpq %rbx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblNotFound))
	// Check for \r\n\r\n
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $13, %al\n") // \r?
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblCheck))
	cg.textSection.WriteString("    movzbq 1(%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $10, %al\n") // \n?
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblCheck))
	cg.textSection.WriteString("    movzbq 2(%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $13, %al\n") // \r?
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblCheck))
	cg.textSection.WriteString("    movzbq 3(%rcx), %rax\n")
	cg.textSection.WriteString("    cmpb $10, %al\n") // \n?
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCheck))
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    addq $4, %rcx\n") // skip \r\n\r\n
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateHTTPParseHeaders populates a buffer with header pointers
// Args: response_buffer, buffer_len, out_headers_array
// This is a simplified version that returns count of headers
func generateHTTPParseHeaders(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // buffer ptr
	cg.generateExpressionToReg(args[1], "r12") // buffer len
	cg.generateExpressionToReg(args[2], "r13") // out array ptr

	lblLoop := cg.getLabel("hdrs_loop")
	lblLine := cg.getLabel("hdrs_line")
	lblNext := cg.getLabel("hdrs_next")
	lblDone := cg.getLabel("hdrs_done")

	// Skip status line
	cg.textSection.WriteString("    movq %rbx, %r14\n") // current pos
	cg.textSection.WriteString("    addq %r12, %rbx\n") // end
	cg.textSection.WriteString("    xorq %r15, %r15\n") // header count

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    cmpq %rbx, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblDone))
	cg.textSection.WriteString("    movzbq (%r14), %rax\n")
	cg.textSection.WriteString("    incq %r14\n")
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblLoop))

	// Process headers
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLine))
	cg.textSection.WriteString("    cmpq %rbx, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblDone))
	cg.textSection.WriteString("    movzbq (%r14), %rax\n")
	// Check for empty line
	cg.textSection.WriteString("    cmpb $13, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblDone))
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblDone))

	// Store pointer to this header line
	cg.textSection.WriteString("    movq %r14, (%r13)\n")
	cg.textSection.WriteString("    addq $8, %r13\n")
	cg.textSection.WriteString("    incq %r15\n")

	// Find end of line
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    cmpq %rbx, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lblDone))
	cg.textSection.WriteString("    movzbq (%r14), %rax\n")
	cg.textSection.WriteString("    incq %r14\n")
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLine))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movq %r15, %rax\n")
}

// =============================================================================
// HTTP Connection Pooling Implementation
// =============================================================================
// Pool structure (per slot - 24 bytes):
//   offset 0: fd (8 bytes) - -1 if unused
//   offset 8: host hash (8 bytes) - djb2 hash of host string
//   offset 16: port (8 bytes)
// Pool header (16 bytes):
//   offset 0: max_slots (8 bytes)
//   offset 8: used_count (8 bytes)
// Total size = 16 + max_slots * 24

const httpPoolSlotSize = 24
const httpPoolHeaderSize = 16

// generateHTTPPoolNew creates a new connection pool
// Args: max_connections
// Returns: pool pointer
func generateHTTPPoolNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	cg.generateExpressionToReg(args[0], "rbx") // max_conns

	// Calculate size: 16 + max_conns * 24
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    imulq $%d, %%rbx, %%rsi\n", httpPoolSlotSize))
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", httpPoolHeaderSize))

	// mmap allocation
	cg.textSection.WriteString("    movq $9, %rax\n")   // mmap
	cg.textSection.WriteString("    xorq %rdi, %rdi\n") // addr = NULL
	cg.textSection.WriteString("    movq $3, %rdx\n")   // PROT_READ | PROT_WRITE
	cg.textSection.WriteString("    movq $34, %r10\n")  // MAP_PRIVATE | MAP_ANONYMOUS
	cg.textSection.WriteString("    movq $-1, %r8\n")   // fd = -1
	cg.textSection.WriteString("    xorq %r9, %r9\n")   // offset = 0
	cg.textSection.WriteString("    syscall\n")

	cg.textSection.WriteString("    popq %rbx\n")
	cg.textSection.WriteString("    movq %rbx, (%rax)\n") // store max_slots
	cg.textSection.WriteString("    movq $0, 8(%rax)\n")  // used_count = 0

	// Initialize all slots to -1 (unused)
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", httpPoolHeaderSize))

	lblLoop := cg.getLabel("pool_init")
	lblDone := cg.getLabel("pool_init_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    testq %rbx, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq $-1, (%rcx)\n")  // fd = -1
	cg.textSection.WriteString("    movq $0, 8(%rcx)\n")  // host_hash = 0
	cg.textSection.WriteString("    movq $0, 16(%rcx)\n") // port = 0
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", httpPoolSlotSize))
	cg.textSection.WriteString("    decq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	// rax already has pool pointer
}

// generateHTTPPoolGet retrieves a connection from the pool
// Args: pool_ptr, host_ptr, port
// Returns: fd if found, -1 if not found
func generateHTTPPoolGet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    movq $-1, %rax\n")
		return
	}

	cg.generateExpressionToReg(args[0], "rbx") // pool ptr
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12") // host ptr
	cg.textSection.WriteString("    pushq %r12\n")
	cg.generateExpressionToReg(args[2], "r13") // port
	cg.textSection.WriteString("    popq %r12\n")
	cg.textSection.WriteString("    popq %rbx\n")

	// Compute djb2 hash of host
	cg.textSection.WriteString("    movq $5381, %r14\n") // hash
	cg.textSection.WriteString("    movq %r12, %rdi\n")

	lblHashLoop := cg.getLabel("pool_hash")
	lblHashDone := cg.getLabel("pool_hash_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblHashLoop))
	cg.textSection.WriteString("    movzbq (%rdi), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblHashDone))
	cg.textSection.WriteString("    imulq $33, %r14\n")
	cg.textSection.WriteString("    addq %rax, %r14\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblHashLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblHashDone))

	// r14 = host hash, r13 = port, rbx = pool
	// Search for matching slot
	cg.textSection.WriteString("    movq (%rbx), %rcx\n")                                      // max_slots
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%rbx), %%rdi\n", httpPoolHeaderSize)) // slot ptr

	lblSearch := cg.getLabel("pool_search")
	lblFound := cg.getLabel("pool_found")
	lblNotFound := cg.getLabel("pool_not_found")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNotFound))

	cg.textSection.WriteString("    movq (%rdi), %rax\n") // fd
	cg.textSection.WriteString("    cmpq $-1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je 1f\n")) // skip if unused

	// Check host hash
	cg.textSection.WriteString("    cmpq %r14, 8(%rdi)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne 1f\n"))
	// Check port
	cg.textSection.WriteString("    cmpq %r13, 16(%rdi)\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))

	cg.textSection.WriteString("1:\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rdi\n", httpPoolSlotSize))
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq (%rdi), %rax\n") // return fd
	cg.textSection.WriteString("    movq $-1, (%rdi)\n")  // mark as unused
	cg.textSection.WriteString("    decq 8(%rbx)\n")      // decrement used_count
	cg.textSection.WriteString(fmt.Sprintf("    jmp 2f\n"))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

// generateHTTPPoolPut stores a connection in the pool
// Args: pool_ptr, fd, host_ptr, port
// Returns: 1 if stored, 0 if pool full
func generateHTTPPoolPut(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 4 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Evaluate all args with proper register preservation
	cg.generateExpressionToReg(args[0], "rbx") // pool ptr
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12") // fd
	cg.textSection.WriteString("    pushq %r12\n")
	cg.generateExpressionToReg(args[2], "r13") // host ptr
	cg.textSection.WriteString("    pushq %r13\n")
	cg.generateExpressionToReg(args[3], "r14") // port
	cg.textSection.WriteString("    pushq %r14\n")
	// Now pop in reverse order
	cg.textSection.WriteString("    popq %r14\n")
	cg.textSection.WriteString("    popq %r13\n")
	cg.textSection.WriteString("    popq %r12\n")
	cg.textSection.WriteString("    popq %rbx\n")

	// Compute djb2 hash of host (r13 -> r15)
	cg.textSection.WriteString("    movq $5381, %r15\n")
	cg.textSection.WriteString("    movq %r13, %rdi\n")

	lblHashLoop := cg.getLabel("poolput_hash")
	lblHashDone := cg.getLabel("poolput_hash_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblHashLoop))
	cg.textSection.WriteString("    movzbq (%rdi), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblHashDone))
	cg.textSection.WriteString("    imulq $33, %r15\n")
	cg.textSection.WriteString("    addq %rax, %r15\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblHashLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblHashDone))

	// r15 = host hash, r14 = port, r12 = fd, rbx = pool
	// Find first empty slot
	cg.textSection.WriteString("    movq (%rbx), %rcx\n") // max_slots
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%rbx), %%rdi\n", httpPoolHeaderSize))

	lblSearch := cg.getLabel("poolput_search")
	lblFound := cg.getLabel("poolput_found")
	lblFull := cg.getLabel("poolput_full")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblFull))

	cg.textSection.WriteString("    cmpq $-1, (%rdi)\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))

	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rdi\n", httpPoolSlotSize))
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq %r12, (%rdi)\n")   // fd
	cg.textSection.WriteString("    movq %r15, 8(%rdi)\n")  // host hash
	cg.textSection.WriteString("    movq %r14, 16(%rdi)\n") // port
	cg.textSection.WriteString("    incq 8(%rbx)\n")        // increment used_count
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp 1f\n"))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFull))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("1:\n")
}

// generateHTTPPoolClose closes all connections and frees the pool
// Args: pool_ptr
// Returns: number of connections closed
func generateHTTPPoolClose(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	cg.generateExpressionToReg(args[0], "rbx") // pool ptr

	cg.textSection.WriteString("    movq (%rbx), %rcx\n") // max_slots
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%rbx), %%rdi\n", httpPoolHeaderSize))
	cg.textSection.WriteString("    xorq %r12, %r12\n") // closed count
	cg.textSection.WriteString("    pushq %rbx\n")      // save pool ptr for munmap

	lblLoop := cg.getLabel("poolclose_loop")
	lblSkip := cg.getLabel("poolclose_skip")
	lblDone := cg.getLabel("poolclose_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))

	cg.textSection.WriteString("    movq (%rdi), %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblSkip))

	// Close the socket
	cg.textSection.WriteString("    pushq %rcx\n")
	cg.textSection.WriteString("    pushq %rdi\n")
	cg.textSection.WriteString("    movq %rax, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rax\n") // close syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %rdi\n")
	cg.textSection.WriteString("    popq %rcx\n")
	cg.textSection.WriteString("    incq %r12\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSkip))
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rdi\n", httpPoolSlotSize))
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))

	// munmap the pool
	cg.textSection.WriteString("    popq %rdi\n")         // pool ptr
	cg.textSection.WriteString("    movq (%rdi), %rsi\n") // max_slots
	cg.textSection.WriteString(fmt.Sprintf("    imulq $%d, %%rsi\n", httpPoolSlotSize))
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", httpPoolHeaderSize))
	cg.textSection.WriteString("    pushq %r12\n")
	cg.textSection.WriteString("    movq $11, %rax\n") // munmap
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %rax\n") // return closed count
}

func generateStringStartsWith(cg *CodeGenerator, args []ASTNode) {
	// startsWith(s, prefix): compare from start; return 1/0
	if len(args) != 2 {
		return
	}
	// Resolve pointers in rbx (s) and rcx (prefix)
	// Reuse logic from compare/indexOf
	// s -> rbx
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rbx\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %rbx, %rbx\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rbx, %rbx\n")
	}
	// prefix -> rcx
	switch v := args[1].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rcx\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rcx\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %rcx, %rcx\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rcx, %rcx\n")
	}
	lLoop := cg.getLabel("sw_loop")
	lEnd := cg.getLabel("sw_end")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    movzbq (%rcx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    movzbq (%rbx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rax, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s_false\n", lEnd))
	cg.textSection.WriteString("    inc %rbx\n")
	cg.textSection.WriteString("    inc %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s_false:\n", lEnd))
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_done\n", lEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_done:\n", lEnd))
}

func generateStringEndsWith(cg *CodeGenerator, args []ASTNode) {
	// endsWith(s, suffix): naive by computing lengths and comparing from back
	if len(args) != 2 {
		return
	}
	// Compute len(s) in r8 and len(suffix) in r9, pointers in rbx (s) rcx (suffix)
	// Resolve pointers
	var lbl string
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ = emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rbx\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
			ln := 0
			if l, ok := cg.stringLengths[v.Name]; ok {
				ln = l
			}
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", ln))
		} else {
			cg.textSection.WriteString("    xorq %rbx, %rbx\n    xorq %r8, %r8\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rbx, %rbx\n    xorq %r8, %r8\n")
	}
	switch v := args[1].(type) {
	case *StringLiteral:
		lbl, _ = emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rcx\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rcx\n", varInfo.Offset))
			ln := 0
			if l, ok := cg.stringLengths[v.Name]; ok {
				ln = l
			}
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", ln))
		} else {
			cg.textSection.WriteString("    xorq %rcx, %rcx\n    xorq %r9, %r9\n")
		}
	default:
		cg.textSection.WriteString("    xorq %rcx, %rcx\n    xorq %r9, %r9\n")
	}
	// if r9 > r8 -> false
	lFalse := cg.getLabel("ew_false")
	lTrue := cg.getLabel("ew_true")
	lEnd := cg.getLabel("ew_end")
	cg.textSection.WriteString("    cmpq %r9, %r8\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lFalse))
	// set pointers to ends
	cg.textSection.WriteString("    movq %rbx, %r10\n    addq %r8, %r10\n    dec %r10\n")
	cg.textSection.WriteString("    movq %rcx, %r11\n    addq %r9, %r11\n    dec %r11\n")
	lLoop := cg.getLabel("ew_loop")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    testq %r9, %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lTrue))
	cg.textSection.WriteString("    movzbq (%r10), %rax\n")
	cg.textSection.WriteString("    movzbq (%r11), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lFalse))
	cg.textSection.WriteString("    dec %r10\n    dec %r11\n    dec %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrue))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lFalse))
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// printfuncs.go wrappers (will use existing implementations from printfuncs.go)
func generatePrintCode(cg *CodeGenerator, args []ASTNode) {
	// Print each argument
	for _, arg := range args {
		emitPrintValue(cg, arg)
	}
}

// Note: Printf and Printf in module system delegate to printfuncs implementations
// The io module functions reference the existing printf/fprintf/sprintf implementations

// ========================= Number conversions =============================
func generateNumToInt8(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    salq $56, %rax\n    sarq $56, %rax\n")
}
func generateNumToUint8(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    andq $255, %rax\n")
}
func generateNumToInt16(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    salq $48, %rax\n    sarq $48, %rax\n")
}
func generateNumToUint16(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    andq $65535, %rax\n")
}
func generateNumToInt32(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    cltq\n")
}
func generateNumToUint32(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movl %eax, %eax\n")
}
func generateNumToInt64(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	// no-op, already 64-bit
}
func generateNumToUint64(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	// no-op for now
}
func generateNumToBool(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    testq %rax, %rax\n    setne %al\n    movzbq %al, %rax\n")
}

// ========================= Memory implementations ==========================

// generateMemSizeof: sizeof(x) -> returns byte size of variable's type if identifier
func generateMemSizeof(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	switch v := args[0].(type) {
	case *Identifier:
		if info, ok := cg.variables[v.Name]; ok {
			sz := GetTypeSize(info.Type)
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", sz))
			return
		}
	case *StringLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", len(v.Value)))
		return
	case *IntLiteral:
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rax\n", GetTypeSize(TokenTypeInt32)))
		return
	}
	cg.textSection.WriteString("    movq $0, %rax\n")
}

// memcpy(dst, src, n): returns dst
func generateMemMemcpy(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	// dst -> rdi, src -> rsi, n -> rcx
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.generateExpressionToReg(args[2], "rcx")
	lLoop := cg.getLabel("memcpy_loop")
	lEnd := cg.getLabel("memcpy_end")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    movb (%rsi), %al\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    inc %rsi\n    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    movq %rdi, %rax\n    subq %rcx, %rax\n")
}

// memset(dst, value, n): returns dst (value treated as byte)
func generateMemMemset(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rax")
	cg.generateExpressionToReg(args[2], "rcx")
	lLoop := cg.getLabel("memset_loop")
	lEnd := cg.getLabel("memset_end")
	cg.textSection.WriteString("    movb %al, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
	cg.textSection.WriteString("    movb %dl, (%rdi)\n")
	cg.textSection.WriteString("    inc %rdi\n    dec %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	cg.textSection.WriteString("    movq %rdi, %rax\n")
}

// malloc(size): implement via mmap(size) and return pointer
func generateMemMalloc(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	// Evaluate size first before setting up syscall registers
	cg.generateExpressionToReg(args[0], "rax")

	// Save current stack pointer and align to 16 bytes for syscall
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n") // Align stack to 16 bytes

	// Setup mmap syscall: rax=9, rdi=NULL, rsi=len, rdx=PROT_READ|PROT_WRITE(3), r10=MAP_PRIVATE|MAP_ANONYMOUS(0x22), r8=-1, r9=0
	cg.textSection.WriteString("    movq %rax, %rsi\n") // Move size from rax to rsi
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	// Restore stack pointer
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")
	// ptr in rax
}

// mmap(size): same as malloc
func generateMemMmap(cg *CodeGenerator, args []ASTNode) { generateMemMalloc(cg, args) }

// free(ptr): currently no-op; recommend mem.munmap
func generateMemFreeNoop(cg *CodeGenerator, args []ASTNode) {
	// Intentionally a no-op; munmap should be used
}

// munmap(ptr, size): rax=11, rdi=addr, rsi=len
func generateMemMunmap(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Save and align stack for syscall
	cg.textSection.WriteString("    pushq %rbp\n")
	cg.textSection.WriteString("    movq %rsp, %rbp\n")
	cg.textSection.WriteString("    andq $-16, %rsp\n")

	cg.generateExpressionToReg(args[0], "rdi")
	cg.generateExpressionToReg(args[1], "rsi")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")

	// Restore stack
	cg.textSection.WriteString("    movq %rbp, %rsp\n")
	cg.textSection.WriteString("    popq %rbp\n")
}

// IO module wrapper functions - delegate to printfuncs.go implementations
func generateIOPrintf(cg *CodeGenerator, args []ASTNode) {
	generatePrintfCode(cg, args)
}

func generateIOFprintf(cg *CodeGenerator, args []ASTNode) {
	generateFprintfCode(cg, args)
}

func generateIOSprint(cg *CodeGenerator, args []ASTNode) {
	generateSprintCode(cg, args)
}

func generateIOSprintf(cg *CodeGenerator, args []ASTNode) {
	generateSprintfCode(cg, args)
}

func generateIOSprintln(cg *CodeGenerator, args []ASTNode) {
	generateSprintlnCode(cg, args)
}

// ============================================================================
// Collections module - first-pass implementations
// ============================================================================

const collectionsHeaderSize = 40 // len(0), cap(8), head(16), tail(24), data ptr(32)

// allocate header + backing store (cap * elemSize bytes) using mmap; cap is preserved in %rbx
func collectionsAllocSized(cg *CodeGenerator, capExpr ASTNode, elemSize int) {
	cg.generateExpressionToReg(capExpr, "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n") // cap
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    imulq $%d, %%rcx\n", elemSize))
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
}

func collectionsAlloc(cg *CodeGenerator, capExpr ASTNode) { collectionsAllocSized(cg, capExpr, 8) }

// store common header fields and return base pointer in %rax
func collectionsInitHeader(cg *CodeGenerator) {
	cg.textSection.WriteString("    movq %rax, %rdx\n")    // base
	cg.textSection.WriteString("    movq $0, (%rdx)\n")    // len
	cg.textSection.WriteString("    movq %rbx, 8(%rdx)\n") // cap
	cg.textSection.WriteString("    movq $0, 16(%rdx)\n")  // head
	cg.textSection.WriteString("    movq $0, 24(%rdx)\n")  // tail
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%rdx), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%rdx)\n") // data ptr
	cg.textSection.WriteString("    movq %rdx, %rax\n")
}

// Dynamic array (int)
func generateCollectionsArrayIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	collectionsAlloc(cg, args[0])
	collectionsInitHeader(cg)
}

func generateCollectionsArrayIntPush(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// arr ptr -> rbx, value -> rcx
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n") // cap
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 1f\n")             // full
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n") // data ptr
	cg.textSection.WriteString("    movq %rcx, (%r8,%rax,8)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n") // return new len
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n") // overflow
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsArrayIntPop(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8,%rax,8), %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsArrayIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

// generateCollectionsArrayIntCapacity returns the current capacity
func generateCollectionsArrayIntCapacity(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rax\n")
}

// generateCollectionsArrayIntGet gets element at index
// Args: array_ptr, index
// Returns: value at index, or 0 if out of bounds
func generateCollectionsArrayIntGet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    popq %rbx\n")

	cg.textSection.WriteString("    movq (%rbx), %rax\n") // len
	cg.textSection.WriteString("    cmpq %rax, %rcx\n")
	cg.textSection.WriteString("    jae 1f\n")             // index >= len
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n") // data ptr
	cg.textSection.WriteString("    movq (%r8,%rcx,8), %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  xorq %rax, %rax\n")
	cg.textSection.WriteString("2:\n")
}

// generateCollectionsArrayIntSet sets element at index
// Args: array_ptr, index, value
// Returns: 1 if success, 0 if out of bounds
func generateCollectionsArrayIntSet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    pushq %r12\n")
	cg.generateExpressionToReg(args[2], "r13")
	cg.textSection.WriteString("    popq %r12\n")
	cg.textSection.WriteString("    popq %rbx\n")

	cg.textSection.WriteString("    movq (%rbx), %rax\n") // len
	cg.textSection.WriteString("    cmpq %rax, %r12\n")
	cg.textSection.WriteString("    jae 1f\n")             // index >= len
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n") // data ptr
	cg.textSection.WriteString("    movq %r13, (%r8,%r12,8)\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  xorq %rax, %rax\n")
	cg.textSection.WriteString("2:\n")
}

// generateCollectionsArrayIntResize resizes the array to new capacity
// Args: array_ptr, new_capacity
// Returns: new array pointer (may be different if reallocated)
// Uses mmap for new allocation, copies data, munmaps old
func generateCollectionsArrayIntResize(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // old array ptr
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12") // new capacity
	cg.textSection.WriteString("    popq %rbx\n")

	// Save old array info
	cg.textSection.WriteString("    pushq %rbx\n")         // old ptr
	cg.textSection.WriteString("    movq (%rbx), %r13\n")  // old len
	cg.textSection.WriteString("    movq 8(%rbx), %r14\n") // old cap
	cg.textSection.WriteString("    pushq %r13\n")
	cg.textSection.WriteString("    pushq %r14\n")

	// Calculate new size: header (40) + new_cap * 8
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    imulq $8, %r12, %rax\n")
	cg.textSection.WriteString("    addq %rax, %rsi\n")

	// mmap new allocation
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    pushq %r12\n") // save new cap
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %r12\n")

	// rax = new ptr, set up header
	cg.textSection.WriteString("    movq %rax, %r15\n") // save new ptr
	cg.textSection.WriteString("    popq %r14\n")       // old cap
	cg.textSection.WriteString("    popq %r13\n")       // old len

	// Set new len = min(old_len, new_cap)
	cg.textSection.WriteString("    cmpq %r12, %r13\n")
	cg.textSection.WriteString("    jbe 1f\n")
	cg.textSection.WriteString("    movq %r12, %r13\n") // truncate
	cg.textSection.WriteString("1:\n")

	cg.textSection.WriteString("    movq %r13, (%r15)\n")  // len
	cg.textSection.WriteString("    movq %r12, 8(%r15)\n") // cap
	cg.textSection.WriteString("    movq $0, 16(%r15)\n")  // head
	cg.textSection.WriteString("    movq $0, 24(%r15)\n")  // tail
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r15), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r15)\n") // data ptr

	// Copy old data
	cg.textSection.WriteString("    popq %rbx\n")           // old ptr
	cg.textSection.WriteString("    movq 32(%rbx), %rsi\n") // old data ptr
	cg.textSection.WriteString("    movq %r13, %rdx\n")     // count to copy

	lblCopy := cg.getLabel("arr_copy")
	lblCopyDone := cg.getLabel("arr_copy_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopy))
	cg.textSection.WriteString("    testq %rdx, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblCopyDone))
	cg.textSection.WriteString("    movq (%rsi), %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rcx)\n")
	cg.textSection.WriteString("    addq $8, %rsi\n")
	cg.textSection.WriteString("    addq $8, %rcx\n")
	cg.textSection.WriteString("    decq %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCopy))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopyDone))

	// munmap old array
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    imulq $8, %r14, %rax\n")
	cg.textSection.WriteString("    addq %rax, %rsi\n")
	cg.textSection.WriteString("    pushq %r15\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %rax\n") // return new ptr
}

// generateCollectionsArrayIntReserve ensures capacity >= new_capacity
// Args: array_ptr, min_capacity
// Returns: array pointer (may be new if resized)
func generateCollectionsArrayIntReserve(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    popq %rbx\n")

	cg.textSection.WriteString("    movq 8(%rbx), %rax\n") // current cap
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString("    jae 1f\n") // already big enough

	// Need to resize - use 2x growth or min_capacity, whichever is larger
	cg.textSection.WriteString("    shlq $1, %rax\n") // 2x current cap
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    movq %r12, %rax\n") // use min_capacity
	cg.textSection.WriteString("2:\n")

	// Call resize with new capacity in rax
	cg.textSection.WriteString("    pushq %rax\n")
	cg.textSection.WriteString("    pushq %rbx\n")

	// Calculate new size
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    imulq $8, %rax, %rcx\n")
	cg.textSection.WriteString("    addq %rcx, %rsi\n")

	// mmap
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	cg.textSection.WriteString("    popq %rbx\n")       // old ptr
	cg.textSection.WriteString("    popq %r12\n")       // new cap
	cg.textSection.WriteString("    movq %rax, %r15\n") // new ptr

	// Copy header
	cg.textSection.WriteString("    movq (%rbx), %r13\n")  // old len
	cg.textSection.WriteString("    movq 8(%rbx), %r14\n") // old cap
	cg.textSection.WriteString("    movq %r13, (%r15)\n")  // len
	cg.textSection.WriteString("    movq %r12, 8(%r15)\n") // new cap
	cg.textSection.WriteString("    movq $0, 16(%r15)\n")
	cg.textSection.WriteString("    movq $0, 24(%r15)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r15), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r15)\n")

	// Copy data
	cg.textSection.WriteString("    movq 32(%rbx), %rsi\n")
	cg.textSection.WriteString("    movq %r13, %rdx\n")

	lblCopy2 := cg.getLabel("res_copy")
	lblCopyDone2 := cg.getLabel("res_copy_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopy2))
	cg.textSection.WriteString("    testq %rdx, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblCopyDone2))
	cg.textSection.WriteString("    movq (%rsi), %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rcx)\n")
	cg.textSection.WriteString("    addq $8, %rsi\n")
	cg.textSection.WriteString("    addq $8, %rcx\n")
	cg.textSection.WriteString("    decq %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCopy2))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopyDone2))

	// munmap old
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    imulq $8, %r14, %rax\n")
	cg.textSection.WriteString("    addq %rax, %rsi\n")
	cg.textSection.WriteString("    pushq %r15\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString("    jmp 3f\n")

	cg.textSection.WriteString("1:  movq %rbx, %rax\n") // no resize needed
	cg.textSection.WriteString("3:\n")
}

// generateCollectionsArrayIntShrink shrinks capacity to match length
// Args: array_ptr
// Returns: new array pointer
func generateCollectionsArrayIntShrink(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")

	cg.textSection.WriteString("    movq (%rbx), %r12\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %r13\n") // cap

	// If len == cap, nothing to do
	cg.textSection.WriteString("    cmpq %r13, %r12\n")
	cg.textSection.WriteString("    je 1f\n")

	// If len == 0, set to minimum capacity of 1
	cg.textSection.WriteString("    testq %r12, %r12\n")
	cg.textSection.WriteString("    jnz 2f\n")
	cg.textSection.WriteString("    movq $1, %r12\n")
	cg.textSection.WriteString("2:\n")

	// Allocate new with len as capacity
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.textSection.WriteString("    pushq %r12\n")
	cg.textSection.WriteString("    pushq %r13\n")

	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    imulq $8, %r12, %rax\n")
	cg.textSection.WriteString("    addq %rax, %rsi\n")

	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	cg.textSection.WriteString("    popq %r13\n") // old cap
	cg.textSection.WriteString("    popq %r12\n") // len (new cap)
	cg.textSection.WriteString("    popq %rbx\n") // old ptr

	cg.textSection.WriteString("    movq %rax, %r15\n") // new ptr

	// Set header
	cg.textSection.WriteString("    movq (%rbx), %r14\n")  // actual len
	cg.textSection.WriteString("    movq %r14, (%r15)\n")  // len
	cg.textSection.WriteString("    movq %r12, 8(%r15)\n") // cap = len
	cg.textSection.WriteString("    movq $0, 16(%r15)\n")
	cg.textSection.WriteString("    movq $0, 24(%r15)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r15), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r15)\n")

	// Copy data
	cg.textSection.WriteString("    movq 32(%rbx), %rsi\n")
	cg.textSection.WriteString("    movq %r14, %rdx\n")

	lblCopy3 := cg.getLabel("shrink_copy")
	lblCopyDone3 := cg.getLabel("shrink_copy_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopy3))
	cg.textSection.WriteString("    testq %rdx, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblCopyDone3))
	cg.textSection.WriteString("    movq (%rsi), %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rcx)\n")
	cg.textSection.WriteString("    addq $8, %rsi\n")
	cg.textSection.WriteString("    addq $8, %rcx\n")
	cg.textSection.WriteString("    decq %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCopy3))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCopyDone3))

	// munmap old
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    imulq $8, %r13, %rax\n")
	cg.textSection.WriteString("    addq %rax, %rsi\n")
	cg.textSection.WriteString("    pushq %r15\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString("    jmp 3f\n")

	cg.textSection.WriteString("1:  movq %rbx, %rax\n") // no shrink needed
	cg.textSection.WriteString("3:\n")
}

// generateCollectionsArrayIntFree frees the array
// Args: array_ptr
// Returns: 0 on success
func generateCollectionsArrayIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")

	cg.textSection.WriteString("    movq 8(%rdi), %rsi\n") // cap
	cg.textSection.WriteString(fmt.Sprintf("    imulq $8, %%rsi\n"))
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq $11, %rax\n") // munmap
	cg.textSection.WriteString("    syscall\n")
}

// Stack (int) - reuse array layout
func generateCollectionsStackIntNew(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntNew(cg, args)
}
func generateCollectionsStackIntPush(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntPush(cg, args)
}
func generateCollectionsStackIntPop(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntPop(cg, args)
}
func generateCollectionsStackIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

// Queue / Deque (int)
func generateCollectionsQueueIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	collectionsAlloc(cg, args[0])
	collectionsInitHeader(cg)
}

func generateCollectionsQueueIntEnqueue(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n") // cap
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 1f\n")
	cg.textSection.WriteString("    movq 24(%rbx), %r9\n") // tail
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n") // data
	cg.textSection.WriteString("    movq %rcx, (%r8,%r9,8)\n")
	cg.textSection.WriteString("    inc %r9\n")
	cg.textSection.WriteString("    cmpq %rdx, %r9\n")
	cg.textSection.WriteString("    jb 3f\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("3:  movq %r9, 24(%rbx)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsQueueIntDequeue(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r9\n") // head
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8,%r9,8), %rcx\n")
	cg.textSection.WriteString("    inc %r9\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %r9\n")
	cg.textSection.WriteString("    jb 3f\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("3:  movq %r9, 16(%rbx)\n")
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsQueueIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

func generateCollectionsDequeIntNew(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsQueueIntNew(cg, args)
}
func generateCollectionsDequeIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

func generateCollectionsDequeIntPushFront(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 1f\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r9\n") // head
	cg.textSection.WriteString("    testq %r9, %r9\n")
	cg.textSection.WriteString("    jne 3f\n")
	cg.textSection.WriteString("    movq %rdx, %r9\n")
	cg.textSection.WriteString("3:  dec %r9\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq %rcx, (%r8,%r9,8)\n")
	cg.textSection.WriteString("    movq %r9, 16(%rbx)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsDequeIntPushBack(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	generateCollectionsQueueIntEnqueue(cg, args)
}

func generateCollectionsDequeIntPopFront(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsQueueIntDequeue(cg, args)
}

func generateCollectionsDequeIntPopBack(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movq 24(%rbx), %r9\n") // tail
	cg.textSection.WriteString("    testq %r9, %r9\n")
	cg.textSection.WriteString("    jne 3f\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r9\n")
	cg.textSection.WriteString("3:  dec %r9\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8,%r9,8), %rcx\n")
	cg.textSection.WriteString("    movq %r9, 24(%rbx)\n")
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

// Heap (min-heap, int) backed by array layout
func generateCollectionsHeapIntNew(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntNew(cg, args)
}

func generateCollectionsHeapIntPush(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")  // len
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n") // cap
	cg.textSection.WriteString("    cmpq %rdx, %rax\n")
	cg.textSection.WriteString("    jae 5f\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq %rcx, (%r8,%rax,8)\n") // place at end
	cg.textSection.WriteString("    movq %rax, %r9\n")          // idx
	cg.textSection.WriteString("1:  testq %r9, %r9\n")
	cg.textSection.WriteString("    je 2f\n")
	cg.textSection.WriteString("    movq %r9, %r10\n")
	cg.textSection.WriteString("    dec %r10\n")
	cg.textSection.WriteString("    shrq $1, %r10\n")           // parent = (idx-1)/2
	cg.textSection.WriteString("    movq (%r8,%r10,8), %r11\n") // parent val
	cg.textSection.WriteString("    movq (%r8,%r9,8), %r12\n")  // current val
	cg.textSection.WriteString("    cmpq %r11, %r12\n")
	cg.textSection.WriteString("    jge 2f\n")
	cg.textSection.WriteString("    movq %r11, (%r8,%r9,8)\n")
	cg.textSection.WriteString("    movq %r12, (%r8,%r10,8)\n")
	cg.textSection.WriteString("    movq %r10, %r9\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    movq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("5:  movq $-1, %rax\n")
	cg.textSection.WriteString("6:\n")
}

func generateCollectionsHeapIntPop(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 9f\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8), %rcx\n") // result
	cg.textSection.WriteString("    dec %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 8f\n")
	cg.textSection.WriteString("    movq (%r8,%rax,8), %r9\n") // last element
	cg.textSection.WriteString("    movq %r9, (%r8)\n")
	cg.textSection.WriteString("    movq $0, %r10\n") // idx
	cg.textSection.WriteString("3:  movq %r10, %r11\n")
	cg.textSection.WriteString("    shlq $1, %r11\n")
	cg.textSection.WriteString("    inc %r11\n") // left child
	cg.textSection.WriteString("    movq %rax, %r12\n")
	cg.textSection.WriteString("    cmpq %r12, %r11\n")
	cg.textSection.WriteString("    jae 8f\n")
	cg.textSection.WriteString("    movq %r11, %r13\n")
	cg.textSection.WriteString("    inc %r13\n")        // right child
	cg.textSection.WriteString("    movq %r11, %r14\n") // smallest = left
	cg.textSection.WriteString("    cmpq %r12, %r13\n")
	cg.textSection.WriteString("    jae 4f\n")
	cg.textSection.WriteString("    movq (%r8,%r13,8), %r15\n")
	cg.textSection.WriteString("    movq (%r8,%r14,8), %rdi\n")
	cg.textSection.WriteString("    cmpq %rdi, %r15\n")
	cg.textSection.WriteString("    jge 4f\n")
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("4:  movq (%r8,%r14,8), %r15\n") // smallest val
	cg.textSection.WriteString("    movq (%r8,%r10,8), %rdi\n") // current val
	cg.textSection.WriteString("    cmpq %rdi, %r15\n")
	cg.textSection.WriteString("    jge 8f\n")
	cg.textSection.WriteString("    movq %r15, (%r8,%r10,8)\n")
	cg.textSection.WriteString("    movq %rdi, (%r8,%r14,8)\n")
	cg.textSection.WriteString("    movq %r14, %r10\n")
	cg.textSection.WriteString("    jmp 3b\n")
	cg.textSection.WriteString("8:  movq %rcx, %rax\n")
	cg.textSection.WriteString("    jmp 10f\n")
	cg.textSection.WriteString("9:  movq $-1, %rax\n")
	cg.textSection.WriteString("10:\n")
}

func generateCollectionsHeapIntPeek(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString("    je 1f\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r8\n")
	cg.textSection.WriteString("    movq (%r8), %rax\n")
	cg.textSection.WriteString("    jmp 2f\n")
	cg.textSection.WriteString("1:  movq $-1, %rax\n")
	cg.textSection.WriteString("2:\n")
}

func generateCollectionsHeapIntLen(cg *CodeGenerator, args []ASTNode) {
	generateCollectionsArrayIntLen(cg, args)
}

// Hash map (int -> int) with hashing, open addressing, and resize (power-of-two cap)
func generateCollectionsHashmapIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// round capacity to power-of-two >= requested
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("1:  cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    shlq $1, %rax\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  movq %rax, %rbx\n") // rbX = cap pow2
	// size = header + buckets(16*cap) + states(1*cap padded to 8)
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    imulq $16, %rcx\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n") // pad states to 8-byte
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// init header
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString("    movq $0, 24(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n") // data ptr
	// Compute states ptr: rcx + rbx*16 (using shlq since scale 16 is invalid)
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    shlq $4, %rdx\n")
	cg.textSection.WriteString("    addq %rcx, %rdx\n")
	cg.textSection.WriteString("    movq %rdx, 16(%r11)\n") // states ptr
	cg.textSection.WriteString("    movq %r11, %rax\n")
}

func hashmapHash(cg *CodeGenerator, keyReg string) {
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rax\n", keyReg))
	cg.textSection.WriteString("    movq %rax, %r14\n")
	cg.textSection.WriteString("    shrq $33, %r14\n")
	cg.textSection.WriteString("    xorq %r14, %rax\n")
	// Use movabsq to load the large constant, then imulq with register
	cg.textSection.WriteString("    movabsq $-352301462168404107, %r14\n")
	cg.textSection.WriteString("    imulq %r14, %rax\n")
	cg.textSection.WriteString("    movq %rax, %r14\n")
	cg.textSection.WriteString("    shrq $33, %r14\n")
	cg.textSection.WriteString("    xorq %r14, %rax\n")
}

func hashmapEnsureCapacity(cg *CodeGenerator, mapReg string) {
	// if len*10 >= cap*7 -> resize (cap *=2)
	cg.textSection.WriteString(fmt.Sprintf("    movq (%%%s), %%r8\n", mapReg))  // len
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r9\n", mapReg)) // cap
	cg.textSection.WriteString("    movq %r8, %rax\n")
	cg.textSection.WriteString("    imulq $10, %rax\n")
	cg.textSection.WriteString("    movq %r9, %rbx\n")
	cg.textSection.WriteString("    imulq $7, %rbx\n")
	cg.textSection.WriteString("    cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jb 9f\n")
	// resize: newCap = cap*2
	cg.textSection.WriteString("    shlq $1, %r9\n")
	// compute new size
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    imulq $16, %rcx\n")
	cg.textSection.WriteString("    movq %r9, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// rax = new base, old base in mapReg
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%r15\n", mapReg))
	// new header init
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rbx\n", mapReg))
	cg.textSection.WriteString("    shlq $1, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    imulq $16, %rdx\n")
	cg.textSection.WriteString("    leaq (%rcx,%rdx,1), %rsi\n")
	cg.textSection.WriteString("    movq %rsi, 16(%r11)\n")
	// reinsert old entries
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r8\n", mapReg))   // old cap
	cg.textSection.WriteString(fmt.Sprintf("    movq 32(%%%s), %%r9\n", mapReg))  // old data
	cg.textSection.WriteString(fmt.Sprintf("    movq 16(%%%s), %%r10\n", mapReg)) // old states
	cg.textSection.WriteString("    xorq %r12, %r12\n")                           // idx = 0
	cg.textSection.WriteString("1:  cmpq %r8, %r12\n")
	cg.textSection.WriteString("    jae 4f\n")
	cg.textSection.WriteString("    movzbq (%r10,%r12,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    jne 3f\n")
	// load key/value - compute r12*16 offset
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %rcx\n")
	cg.textSection.WriteString("    movq 8(%r9,%rdi), %r15\n")
	// insert into new table (reuse put logic using r11 as base)
	cg.textSection.WriteString("    movq %r11, %rbx\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("3:  inc %r12\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("4:  jmp 8f\n")
	// insertion fragment shared with runtime put
	cg.textSection.WriteString("6:  movq (%rbx), %rax\n")   // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")   // cap
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")  // data
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n") // states
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    movq %r15, %rsi\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r13\n") // idx
	cg.textSection.WriteString("    movq $-1, %r14\n")  // tomb = -1
	cg.textSection.WriteString("7:  movzbq (%r10,%r13,1), %r15\n")
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString("    je 10f\n")
	cg.textSection.WriteString("    cmpq $2, %r15\n")
	cg.textSection.WriteString("    jne 9f\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("    jmp 11f\n")
	cg.textSection.WriteString("9:  movq %r13, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 12f\n")
	cg.textSection.WriteString("    jmp 13f\n")
	// compute r13*16 offset into rdi
	cg.textSection.WriteString("10: movq %r13, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %r15\n")
	cg.textSection.WriteString("    cmpq %rcx, %r15\n")
	cg.textSection.WriteString("    jne 11f\n")
	// compute r13*16 offset into rdi again
	cg.textSection.WriteString("    movq %r13, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq %rsi, 8(%r9,%rdi)\n")
	cg.textSection.WriteString("    jmp 12f\n")
	cg.textSection.WriteString("11: inc %r13\n")
	cg.textSection.WriteString("    cmpq %r8, %r13\n")
	cg.textSection.WriteString("    jb 7b\n")
	cg.textSection.WriteString("    xorq %r13, %r13\n")
	cg.textSection.WriteString("    jmp 7b\n")
	cg.textSection.WriteString("12: jmp 14f\n")
	cg.textSection.WriteString("13: movq %r13, %r14\n")
	cg.textSection.WriteString("14: cmpq $-1, %r14\n")
	cg.textSection.WriteString("    je 15f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r14,1)\n")
	// compute r14*16 offset into rdi
	cg.textSection.WriteString("    movq %r14, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%rdi)\n")
	cg.textSection.WriteString("    movq %rsi, 8(%r9,%rdi)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("15: jmp 3b\n")
	cg.textSection.WriteString("8:  movq %r11, %rax\n")
	// free old
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rsi\n", mapReg))
	cg.textSection.WriteString("    imulq $16, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rdx\n", mapReg))
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rdi\n", mapReg))
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	// update mapReg to new base
	cg.textSection.WriteString(fmt.Sprintf("    movq %%r11, %%%s\n", mapReg))
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashmapIntPut(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // map base
	// ensure capacity
	hashmapEnsureCapacity(cg, "rbx")
	cg.generateExpressionToReg(args[1], "rcx")              // key
	cg.generateExpressionToReg(args[2], "rdx")              // value
	cg.textSection.WriteString("    movq %rdx, %r15\n")     // stash value
	cg.textSection.WriteString("    movq (%rbx), %rax\n")   // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")   // cap
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")  // data ptr
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n") // states ptr
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n") // idx
	cg.textSection.WriteString("    movq $-1, %r12\n")  // tomb
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $2, %r13\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 4f\n")
	cg.textSection.WriteString("    movq %r11, %r12\n")
	cg.textSection.WriteString("    jmp 4f\n")
	cg.textSection.WriteString("2:  movq %r11, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 6f\n")
	cg.textSection.WriteString("    jmp 7f\n")
	// compute r11*16 offset into rdi
	cg.textSection.WriteString("3:  movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %r14\n")
	cg.textSection.WriteString("    cmpq %rcx, %r14\n")
	cg.textSection.WriteString("    jne 4f\n")
	// compute r11*16 offset into rdi again
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq %r15, 8(%r9,%rdi)\n")
	cg.textSection.WriteString("    movq %r15, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("4:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("6:  jmp 8f\n")
	cg.textSection.WriteString("7:  movq %r11, %r12\n")
	cg.textSection.WriteString("8:  cmpq $-1, %r12\n")
	cg.textSection.WriteString("    je 9f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r12,1)\n")
	// compute r12*16 offset into rdi
	cg.textSection.WriteString("    movq %r12, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%rdi)\n")
	cg.textSection.WriteString("    movq %r15, 8(%r9,%rdi)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashmapIntGet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $1, %r12\n")
	cg.textSection.WriteString("    jne 2f\n")
	// compute r11*16 offset into rdi
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("2:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $-1, %rax\n")
	cg.textSection.WriteString("    jmp 5f\n")
	// compute r11*16 offset into rdi
	cg.textSection.WriteString("4:  movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq 8(%r9,%rdi), %rax\n")
	cg.textSection.WriteString("5:\n")
}

func generateCollectionsHashmapIntRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %r14\n") // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	// compute r11*16 offset into rdi
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("    inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $-1, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	// compute r11*16 offset into rdi
	cg.textSection.WriteString("4:  movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq 8(%r9,%rdi), %rax\n")
	cg.textSection.WriteString("    movb $2, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    dec %r14\n")
	cg.textSection.WriteString("    movq %r14, (%rbx)\n")
	cg.textSection.WriteString("6:\n")
}

func generateCollectionsHashmapIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

func generateCollectionsHashmapIntClear(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq $0, (%rbx)\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rcx\n")
	cg.textSection.WriteString("    movq 16(%rbx), %rdi\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
}

func generateCollectionsHashmapIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rsi\n")
	cg.textSection.WriteString("    imulq $16, %rsi\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// Hash set (int) with hashing, open addressing, resize
func generateCollectionsHashsetIntNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("1:  cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    shlq $1, %rax\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    imulq $8, %rcx\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString("    movq $0, 24(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    leaq (%rcx,%rbx,8), %rdx\n")
	cg.textSection.WriteString("    movq %rdx, 16(%r11)\n")
	cg.textSection.WriteString("    movq %r11, %rax\n")
}

func hashsetEnsureCapacity(cg *CodeGenerator, mapReg string) {
	cg.textSection.WriteString(fmt.Sprintf("    movq (%%%s), %%r8\n", mapReg))  // len
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r9\n", mapReg)) // cap
	cg.textSection.WriteString("    movq %r8, %rax\n")
	cg.textSection.WriteString("    imulq $10, %rax\n")
	cg.textSection.WriteString("    movq %r9, %rbx\n")
	cg.textSection.WriteString("    imulq $7, %rbx\n")
	cg.textSection.WriteString("    cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jb 9f\n")
	cg.textSection.WriteString("    shlq $1, %r9\n")
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    imulq $8, %rcx\n")
	cg.textSection.WriteString("    movq %r9, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%r15\n", mapReg))
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rbx\n", mapReg))
	cg.textSection.WriteString("    shlq $1, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    imulq $8, %rdx\n")
	cg.textSection.WriteString("    leaq (%rcx,%rdx,1), %rsi\n")
	cg.textSection.WriteString("    movq %rsi, 16(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%r8\n", mapReg))
	cg.textSection.WriteString(fmt.Sprintf("    movq 32(%%%s), %%r9\n", mapReg))
	cg.textSection.WriteString(fmt.Sprintf("    movq 16(%%%s), %%r10\n", mapReg))
	cg.textSection.WriteString("    xorq %r12, %r12\n")
	cg.textSection.WriteString("1:  cmpq %r8, %r12\n")
	cg.textSection.WriteString("    jae 4f\n")
	cg.textSection.WriteString("    movzbq (%r10,%r12,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    jne 3f\n")
	cg.textSection.WriteString("    movq (%r9,%r12,8), %rcx\n")
	cg.textSection.WriteString("    movq %r11, %rbx\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("3:  inc %r12\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("4:  jmp 8f\n")
	cg.textSection.WriteString("6:  movq (%rbx), %rax\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r13\n")
	cg.textSection.WriteString("    movq $-1, %r14\n")
	cg.textSection.WriteString("7:  movzbq (%r10,%r13,1), %r15\n")
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString("    je 10f\n")
	cg.textSection.WriteString("    cmpq $2, %r15\n")
	cg.textSection.WriteString("    jne 9f\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("    jmp 11f\n")
	cg.textSection.WriteString("9:  movq %r13, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    jne 12f\n")
	cg.textSection.WriteString("    jmp 13f\n")
	cg.textSection.WriteString("10: movq (%r9,%r13,8), %r15\n")
	cg.textSection.WriteString("    cmpq %rcx, %r15\n")
	cg.textSection.WriteString("    jne 11f\n")
	cg.textSection.WriteString("    jmp 12f\n")
	cg.textSection.WriteString("11: inc %r13\n")
	cg.textSection.WriteString("    cmpq %r8, %r13\n")
	cg.textSection.WriteString("    jb 7b\n")
	cg.textSection.WriteString("    xorq %r13, %r13\n")
	cg.textSection.WriteString("    jmp 7b\n")
	cg.textSection.WriteString("12: jmp 14f\n")
	cg.textSection.WriteString("13: movq %r13, %r14\n")
	cg.textSection.WriteString("14: cmpq $-1, %r14\n")
	cg.textSection.WriteString("    je 15f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r14,1)\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%r14,8)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("15: jmp 3b\n")
	cg.textSection.WriteString("8:  movq %r11, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rsi\n", mapReg))
	cg.textSection.WriteString("    imulq $8, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq 8(%%%s), %%rdx\n", mapReg))
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rdi\n", mapReg))
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq %%r11, %%%s\n", mapReg))
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashsetIntAdd(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	hashsetEnsureCapacity(cg, "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("    movq $-1, %r12\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r13\n")
	cg.textSection.WriteString("    cmpq $1, %r13\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $2, %r13\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 4f\n")
	cg.textSection.WriteString("    movq %r11, %r12\n")
	cg.textSection.WriteString("    jmp 4f\n")
	cg.textSection.WriteString("2:  movq %r11, %rax\n")
	cg.textSection.WriteString("    cmpq $-1, %r12\n")
	cg.textSection.WriteString("    jne 6f\n")
	cg.textSection.WriteString("    jmp 7f\n")
	cg.textSection.WriteString("3:  movq (%r9,%r11,8), %r14\n")
	cg.textSection.WriteString("    cmpq %rcx, %r14\n")
	cg.textSection.WriteString("    je 6f\n")
	cg.textSection.WriteString("4:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("6:  jmp 8f\n")
	cg.textSection.WriteString("7:  movq %r11, %r12\n")
	cg.textSection.WriteString("8:  cmpq $-1, %r12\n")
	cg.textSection.WriteString("    je 9f\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r12,1)\n")
	cg.textSection.WriteString("    movq %rcx, (%r9,%r12,8)\n")
	cg.textSection.WriteString("    inc %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString("9:\n")
}

func generateCollectionsHashsetIntContains(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    cmpq $1, %r12\n")
	cg.textSection.WriteString("    jne 2f\n")
	cg.textSection.WriteString("    movq (%r9,%r11,8), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("2:  inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $0, %rax\n")
	cg.textSection.WriteString("    jmp 5f\n")
	cg.textSection.WriteString("4:  movq $1, %rax\n")
	cg.textSection.WriteString("5:\n")
}

func generateCollectionsHashsetIntRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "rcx")
	cg.textSection.WriteString("    movq (%rbx), %r14\n")
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")
	hashmapHash(cg, "rcx")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("1:  movzbq (%r10,%r11,1), %r12\n")
	cg.textSection.WriteString("    cmpq $0, %r12\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    movq (%r9,%r11,8), %r13\n")
	cg.textSection.WriteString("    cmpq %rcx, %r13\n")
	cg.textSection.WriteString("    je 4f\n")
	cg.textSection.WriteString("    inc %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString("    jb 1b\n")
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq $-1, %rax\n")
	cg.textSection.WriteString("    jmp 6f\n")
	cg.textSection.WriteString("4:  movq %rcx, %rax\n")
	cg.textSection.WriteString("    movb $2, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    dec %r14\n")
	cg.textSection.WriteString("    movq %r14, (%rbx)\n")
	cg.textSection.WriteString("6:\n")
}

func generateCollectionsHashsetIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

func generateCollectionsHashsetIntClear(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq $0, (%rbx)\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rcx\n")
	cg.textSection.WriteString("    movq 16(%rbx), %rdi\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
}

func generateCollectionsHashsetIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rsi\n")
	cg.textSection.WriteString("    imulq $8, %rsi\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// ============================================================================
// String-key HashMap and HashSet implementations
// ============================================================================
// These use the djb2 hash algorithm for string keys and strcmp for comparison
// Layout is similar to int variants but stores string pointers as keys

// Helper: compute djb2 hash of null-terminated string
// Input: string ptr in specified register
// Output: hash in %rax
func hashmapStrHash(cg *CodeGenerator, ptrReg string) {
	cg.textSection.WriteString(fmt.Sprintf("    movq %%%s, %%rsi\n", ptrReg))
	cg.textSection.WriteString("    movq $5381, %rax\n") // djb2 starting value
	lblLoop := cg.getLabel("djb2_loop")
	lblDone := cg.getLabel("djb2_done")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    movzbl (%rsi), %edx\n")
	cg.textSection.WriteString("    testb %dl, %dl\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	// hash = hash * 33 + c = (hash << 5) + hash + c
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    shlq $5, %rax\n")
	cg.textSection.WriteString("    addq %rcx, %rax\n")
	cg.textSection.WriteString("    addq %rdx, %rax\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// Helper: compare two null-terminated strings
// Input: str1 in rdi, str2 in rsi
// Output: rax = 1 if equal, 0 if not
func hashmapStrCmp(cg *CodeGenerator) {
	lblLoop := cg.getLabel("strcmp_loop")
	lblNotEqual := cg.getLabel("strcmp_ne")
	lblEqual := cg.getLabel("strcmp_eq")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    movzbl (%rdi), %eax\n")
	cg.textSection.WriteString("    movzbl (%rsi), %edx\n")
	cg.textSection.WriteString("    cmpb %dl, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNotEqual))
	cg.textSection.WriteString("    testb %al, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblEqual))
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotEqual))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblEqual))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsHashmapStrNew creates a string-keyed hashmap
// Args: initial capacity
// Returns: pointer to hashmap structure
func generateCollectionsHashmapStrNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Same structure as int hashmap: header + buckets(16*cap) + states(1*cap)
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("1:  cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    shlq $1, %rax\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    imulq $16, %rcx\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString("    movq $0, 24(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    shlq $4, %rdx\n")
	cg.textSection.WriteString("    addq %rcx, %rdx\n")
	cg.textSection.WriteString("    movq %rdx, 16(%r11)\n")
	cg.textSection.WriteString("    movq %r11, %rax\n")
}

// generateCollectionsHashmapStrPut inserts/updates a key-value pair
// Args: map_ptr, string_key, value
func generateCollectionsHashmapStrPut(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // map base
	cg.generateExpressionToReg(args[1], "r12") // key (string ptr)
	cg.generateExpressionToReg(args[2], "r13") // value

	// Get map fields
	cg.textSection.WriteString("    movq (%rbx), %rax\n")   // len
	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")   // cap
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")  // data ptr
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n") // states ptr

	// Hash the key
	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n") // idx
	cg.textSection.WriteString("    movq $-1, %r14\n")  // tomb slot

	// Probe loop
	lblProbe := cg.getLabel("hm_str_probe")
	lblFound := cg.getLabel("hm_str_found")
	lblInsert := cg.getLabel("hm_str_insert")
	lblNext := cg.getLabel("hm_str_next")
	lblDone := cg.getLabel("hm_str_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    cmpq $1, %r15\n") // occupied
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString("    cmpq $2, %r15\n") // tombstone
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblInsert))
	// tombstone - save if first
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))
	cg.textSection.WriteString("    movq %r11, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblNext))

	// Check if key matches
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %rdi\n") // existing key
	cg.textSection.WriteString("    movq %r12, %rsi\n")       // new key
	cg.textSection.WriteString("    pushq %r11\n")
	cg.textSection.WriteString("    pushq %rax\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %rcx\n") // restore old rax (len) to rcx
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNext))
	// Keys match - update value
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq %r13, 8(%r9,%rdi)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblDone))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	// Empty slot - insert here (or at tombstone)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblInsert))
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    cmovneq %r14, %r11\n") // use tomb if available
	cg.textSection.WriteString("    movb $1, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq %r12, (%r9,%rdi)\n")  // store key
	cg.textSection.WriteString("    movq %r13, 8(%r9,%rdi)\n") // store value
	cg.textSection.WriteString("    incq (%rbx)\n")            // len++

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsHashmapStrGet retrieves value for a string key
// Args: map_ptr, string_key
// Returns: value or 0 if not found
func generateCollectionsHashmapStrGet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx") // map base
	cg.generateExpressionToReg(args[1], "r12") // key

	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")   // cap
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")  // data
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n") // states

	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")

	lblProbe := cg.getLabel("hm_str_get_probe")
	lblCheck := cg.getLabel("hm_str_get_check")
	lblNext := cg.getLabel("hm_str_get_next")
	lblNotFound := cg.getLabel("hm_str_get_nf")
	lblFound := cg.getLabel("hm_str_get_found")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    testq %r15, %r15\n") // empty
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNotFound))
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))

	// Compare keys
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCheck))
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %rdi\n")
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    pushq %r11\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNext))
	// Found
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq 8(%r9,%rdi), %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblFound))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
}

// generateCollectionsHashmapStrContains checks if key exists
// Args: map_ptr, string_key
// Returns: 1 if exists, 0 otherwise
func generateCollectionsHashmapStrContains(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "r12")

	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")

	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")

	lblProbe := cg.getLabel("hm_str_has_probe")
	lblCheck := cg.getLabel("hm_str_has_check")
	lblNext := cg.getLabel("hm_str_has_next")
	lblNo := cg.getLabel("hm_str_has_no")
	lblYes := cg.getLabel("hm_str_has_yes")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    testq %r15, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNo))
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCheck))
	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %rdi\n")
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    pushq %r11\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblYes))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNo))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblYes))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsHashmapStrRemove removes a key
// Args: map_ptr, string_key
func generateCollectionsHashmapStrRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "r12")

	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")

	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")

	lblProbe := cg.getLabel("hm_str_rm_probe")
	lblNext := cg.getLabel("hm_str_rm_next")
	lblDone := cg.getLabel("hm_str_rm_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    testq %r15, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))

	cg.textSection.WriteString("    movq %r11, %rdi\n")
	cg.textSection.WriteString("    shlq $4, %rdi\n")
	cg.textSection.WriteString("    movq (%r9,%rdi), %rdi\n")
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    pushq %r11\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNext))
	// Found - mark as tombstone
	cg.textSection.WriteString("    movb $2, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    decq (%rbx)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblDone))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsHashmapStrLen returns number of entries
func generateCollectionsHashmapStrLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

// generateCollectionsHashmapStrClear clears all entries
func generateCollectionsHashmapStrClear(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq $0, (%rbx)\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rcx\n")
	cg.textSection.WriteString("    movq 16(%rbx), %rdi\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
}

// generateCollectionsHashmapStrFree deallocates the hashmap
func generateCollectionsHashmapStrFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rsi\n")
	cg.textSection.WriteString("    imulq $16, %rsi\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// ============================================================================
// String HashSet implementations
// ============================================================================

// generateCollectionsHashsetStrNew creates a string hashset
func generateCollectionsHashsetStrNew(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Structure: header + elements(8*cap for string ptrs) + states(1*cap)
	cg.generateExpressionToReg(args[0], "rax")
	cg.textSection.WriteString("    movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("1:  cmpq %rbx, %rax\n")
	cg.textSection.WriteString("    jae 2f\n")
	cg.textSection.WriteString("    shlq $1, %rax\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  movq %rax, %rbx\n")
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    imulq $8, %rcx\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq %rax, %r11\n")
	cg.textSection.WriteString("    movq $0, (%r11)\n")
	cg.textSection.WriteString("    movq %rbx, 8(%r11)\n")
	cg.textSection.WriteString(fmt.Sprintf("    leaq %d(%%r11), %%rcx\n", collectionsHeaderSize))
	cg.textSection.WriteString("    movq %rcx, 32(%r11)\n")
	cg.textSection.WriteString("    movq %rbx, %rdx\n")
	cg.textSection.WriteString("    shlq $3, %rdx\n")
	cg.textSection.WriteString("    addq %rcx, %rdx\n")
	cg.textSection.WriteString("    movq %rdx, 16(%r11)\n")
	cg.textSection.WriteString("    movq %r11, %rax\n")
}

// generateCollectionsHashsetStrAdd adds a string to the set
func generateCollectionsHashsetStrAdd(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "r12")

	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")

	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")
	cg.textSection.WriteString("    movq $-1, %r14\n")

	lblProbe := cg.getLabel("hs_str_add_probe")
	lblCheck := cg.getLabel("hs_str_add_check")
	lblNext := cg.getLabel("hs_str_add_next")
	lblInsert := cg.getLabel("hs_str_add_ins")
	lblDone := cg.getLabel("hs_str_add_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblCheck))
	cg.textSection.WriteString("    cmpq $2, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblInsert))
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))
	cg.textSection.WriteString("    movq %r11, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblNext))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCheck))
	cg.textSection.WriteString("    movq (%r9,%r11,8), %rdi\n")
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    pushq %r11\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblDone)) // already exists

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblInsert))
	cg.textSection.WriteString("    cmpq $-1, %r14\n")
	cg.textSection.WriteString("    cmovneq %r14, %r11\n")
	cg.textSection.WriteString("    movb $1, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    movq %r12, (%r9,%r11,8)\n")
	cg.textSection.WriteString("    incq (%rbx)\n")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsHashsetStrContains checks if string is in set
func generateCollectionsHashsetStrContains(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "r12")

	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")

	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")

	lblProbe := cg.getLabel("hs_str_has_probe")
	lblCheck := cg.getLabel("hs_str_has_check")
	lblNext := cg.getLabel("hs_str_has_next")
	lblNo := cg.getLabel("hs_str_has_no")
	lblYes := cg.getLabel("hs_str_has_yes")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    testq %r15, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNo))
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCheck))
	cg.textSection.WriteString("    movq (%r9,%r11,8), %rdi\n")
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    pushq %r11\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblYes))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNo))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblYes))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsHashsetStrRemove removes a string from set
func generateCollectionsHashsetStrRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.generateExpressionToReg(args[1], "r12")

	cg.textSection.WriteString("    movq 8(%rbx), %r8\n")
	cg.textSection.WriteString("    movq 32(%rbx), %r9\n")
	cg.textSection.WriteString("    movq 16(%rbx), %r10\n")

	hashmapStrHash(cg, "r12")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %r8\n")
	cg.textSection.WriteString("    movq %rdx, %r11\n")

	lblProbe := cg.getLabel("hs_str_rm_probe")
	lblNext := cg.getLabel("hs_str_rm_next")
	lblDone := cg.getLabel("hs_str_rm_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblProbe))
	cg.textSection.WriteString("    movzbq (%r10,%r11,1), %r15\n")
	cg.textSection.WriteString("    testq %r15, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    cmpq $1, %r15\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lblNext))

	cg.textSection.WriteString("    movq (%r9,%r11,8), %rdi\n")
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    pushq %r11\n")
	hashmapStrCmp(cg)
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNext))
	cg.textSection.WriteString("    movb $2, (%r10,%r11,1)\n")
	cg.textSection.WriteString("    decq (%rbx)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblDone))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNext))
	cg.textSection.WriteString("    incq %r11\n")
	cg.textSection.WriteString("    cmpq %r8, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblProbe))
	cg.textSection.WriteString("    xorq %r11, %r11\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblProbe))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsHashsetStrLen returns size of set
func generateCollectionsHashsetStrLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq (%rbx), %rax\n")
}

// generateCollectionsHashsetStrClear clears the set
func generateCollectionsHashsetStrClear(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq $0, (%rbx)\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rcx\n")
	cg.textSection.WriteString("    movq 16(%rbx), %rdi\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
}

// generateCollectionsHashsetStrFree frees the set
func generateCollectionsHashsetStrFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rsi\n")
	cg.textSection.WriteString("    imulq $8, %rsi\n")
	cg.textSection.WriteString("    movq 8(%rbx), %rdx\n")
	cg.textSection.WriteString("    addq $7, %rdx\n")
	cg.textSection.WriteString("    andq $-8, %rdx\n")
	cg.textSection.WriteString(fmt.Sprintf("    addq $%d, %%rsi\n", collectionsHeaderSize))
	cg.textSection.WriteString("    addq %rdx, %rsi\n")
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
}

// ============================================================================
// Sorted Set (BST-based) implementations
// ============================================================================
// BST node structure (24 bytes):
//   offset 0: key (8 bytes)
//   offset 8: left child ptr (8 bytes)
//   offset 16: right child ptr (8 bytes)
// Sorted set structure (16 bytes):
//   offset 0: root ptr (8 bytes)
//   offset 8: count (8 bytes)

const bstNodeSize = 24

// allocBSTNode allocates a new BST node and returns ptr in rax
// Preserves: rbx, r12, r13, r14, r15
func allocBSTNode(cg *CodeGenerator) {
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.textSection.WriteString("    pushq %r12\n")
	cg.textSection.WriteString("    pushq %r13\n")
	cg.textSection.WriteString("    pushq %r14\n")
	cg.textSection.WriteString("    movq $9, %rax\n")   // mmap
	cg.textSection.WriteString("    xorq %rdi, %rdi\n") // addr = NULL
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", bstNodeSize))
	cg.textSection.WriteString("    movq $3, %rdx\n")  // PROT_READ | PROT_WRITE
	cg.textSection.WriteString("    movq $34, %r10\n") // MAP_PRIVATE | MAP_ANONYMOUS
	cg.textSection.WriteString("    movq $-1, %r8\n")  // fd = -1
	cg.textSection.WriteString("    xorq %r9, %r9\n")  // offset = 0
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %r14\n")
	cg.textSection.WriteString("    popq %r13\n")
	cg.textSection.WriteString("    popq %r12\n")
	cg.textSection.WriteString("    popq %rbx\n")
}

// generateCollectionsSortedsetIntNew creates a new sorted set
// Returns: pointer to set structure (root=NULL, count=0)
func generateCollectionsSortedsetIntNew(cg *CodeGenerator, args []ASTNode) {
	// Allocate 16-byte header
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $16, %rsi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// Initialize: root = NULL, count = 0
	cg.textSection.WriteString("    movq $0, (%rax)\n")
	cg.textSection.WriteString("    movq $0, 8(%rax)\n")
}

// generateCollectionsSortedsetIntAdd adds a value to the sorted set
// Args: set_ptr, value
func generateCollectionsSortedsetIntAdd(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Save rbx across the second generateExpressionToReg call
	cg.generateExpressionToReg(args[0], "rbx")     // set ptr
	cg.textSection.WriteString("    pushq %rbx\n") // save set ptr
	cg.generateExpressionToReg(args[1], "r12")     // value to add
	cg.textSection.WriteString("    popq %rbx\n")  // restore set ptr

	lblInsert := cg.getLabel("ss_add")
	lblFound := cg.getLabel("ss_found")
	lblLeft := cg.getLabel("ss_left")
	lblRight := cg.getLabel("ss_right")
	lblNewNode := cg.getLabel("ss_new")
	lblDone := cg.getLabel("ss_done")

	// Initialize parent to NULL (needed if root is NULL)
	cg.textSection.WriteString("    xorq %r13, %r13\n") // parent = NULL

	// Check if root is NULL
	cg.textSection.WriteString("    movq (%rbx), %rcx\n") // root
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNewNode))

	// Traverse BST to find insertion point
	// rcx = current node, r13 = parent, r14 = direction (0=left, 1=right)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblInsert))
	cg.textSection.WriteString("    movq (%rcx), %rax\n") // current.key
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound)) // already exists
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	// Go right
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq %rcx, %r13\n")     // parent = current
	cg.textSection.WriteString("    movq $1, %r14\n")       // direction = right
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n") // current = current.right
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblInsert))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblNewNode))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq %rcx, %r13\n")    // parent = current
	cg.textSection.WriteString("    xorq %r14, %r14\n")    // direction = left
	cg.textSection.WriteString("    movq 8(%rcx), %rcx\n") // current = current.left
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblInsert))

	// Create new node
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNewNode))
	allocBSTNode(cg)
	// rax = new node ptr
	cg.textSection.WriteString("    movq %r12, (%rax)\n") // node.key = value
	cg.textSection.WriteString("    movq $0, 8(%rax)\n")  // node.left = NULL
	cg.textSection.WriteString("    movq $0, 16(%rax)\n") // node.right = NULL

	// Link to parent
	cg.textSection.WriteString("    testq %r13, %r13\n") // parent == NULL?
	cg.textSection.WriteString(fmt.Sprintf("    jz %s_root\n", lblNewNode))
	cg.textSection.WriteString("    testq %r14, %r14\n") // direction?
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s_linkr\n", lblNewNode))
	cg.textSection.WriteString("    movq %rax, 8(%r13)\n") // parent.left = new
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_inc\n", lblNewNode))
	cg.textSection.WriteString(fmt.Sprintf("%s_linkr:\n", lblNewNode))
	cg.textSection.WriteString("    movq %rax, 16(%r13)\n") // parent.right = new
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_inc\n", lblNewNode))
	cg.textSection.WriteString(fmt.Sprintf("%s_root:\n", lblNewNode))
	cg.textSection.WriteString("    movq %rax, (%rbx)\n") // set.root = new
	cg.textSection.WriteString(fmt.Sprintf("%s_inc:\n", lblNewNode))
	cg.textSection.WriteString("    incq 8(%rbx)\n") // count++
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblDone))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	// Value already exists, do nothing

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsSortedsetIntContains checks if value exists
// Args: set_ptr, value -> returns 1 if found, 0 otherwise
func generateCollectionsSortedsetIntContains(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    popq %rbx\n")

	lblSearch := cg.getLabel("ss_search")
	lblLeft := cg.getLabel("ss_sleft")
	lblRight := cg.getLabel("ss_sright")
	lblFound := cg.getLabel("ss_sfound")
	lblNotFound := cg.getLabel("ss_snf")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n") // root
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNotFound))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq 8(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedsetIntRemove removes a value (simplified - marks as deleted)
// For simplicity, we don't actually restructure the tree
func generateCollectionsSortedsetIntRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// Simplified: just decrement count if found (doesn't actually remove node)
	// Full BST deletion is complex, so this is a placeholder
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    popq %rbx\n")

	lblSearch := cg.getLabel("ss_rm_search")
	lblLeft := cg.getLabel("ss_rm_left")
	lblRight := cg.getLabel("ss_rm_right")
	lblFound := cg.getLabel("ss_rm_found")
	lblDone := cg.getLabel("ss_rm_done")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq 8(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	// Mark node as removed by setting key to special value (min int64)
	cg.textSection.WriteString("    movabsq $-9223372036854775808, %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rcx)\n")
	cg.textSection.WriteString("    decq 8(%rbx)\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsSortedsetIntMin finds the minimum value
func generateCollectionsSortedsetIntMin(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")

	lblLoop := cg.getLabel("ss_min")
	lblDone := cg.getLabel("ss_min_done")
	lblEmpty := cg.getLabel("ss_min_empty")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n") // root
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblEmpty))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    movq 8(%rcx), %rax\n") // left
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movq (%rcx), %rax\n") // return key
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblEmpty))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedsetIntMax finds the maximum value
func generateCollectionsSortedsetIntMax(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")

	lblLoop := cg.getLabel("ss_max")
	lblDone := cg.getLabel("ss_max_done")
	lblEmpty := cg.getLabel("ss_max_empty")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblEmpty))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    movq 16(%rcx), %rax\n") // right
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblEmpty))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedsetIntLen returns the count
func generateCollectionsSortedsetIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rax\n")
}

// generateCollectionsSortedsetIntFree frees the set (simplified - just frees header)
func generateCollectionsSortedsetIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq $16, %rsi\n")
	cg.textSection.WriteString("    movq $11, %rax\n") // munmap
	cg.textSection.WriteString("    syscall\n")
}

// ============================================================================
// Sorted Map (BST-based) implementations
// ============================================================================
// BST node structure for map (32 bytes):
//   offset 0: key (8 bytes)
//   offset 8: value (8 bytes)
//   offset 16: left child ptr (8 bytes)
//   offset 24: right child ptr (8 bytes)

const bstMapNodeSize = 32

// generateCollectionsSortedmapIntNew creates a new sorted map
func generateCollectionsSortedmapIntNew(cg *CodeGenerator, args []ASTNode) {
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $16, %rsi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq $0, (%rax)\n")
	cg.textSection.WriteString("    movq $0, 8(%rax)\n")
}

// generateCollectionsSortedmapIntPut inserts or updates a key-value pair
func generateCollectionsSortedmapIntPut(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	// Save registers across generateExpressionToReg calls
	cg.generateExpressionToReg(args[0], "rbx") // map ptr
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12") // key
	cg.textSection.WriteString("    pushq %r12\n")
	cg.generateExpressionToReg(args[2], "r15") // value
	cg.textSection.WriteString("    popq %r12\n")
	cg.textSection.WriteString("    popq %rbx\n")

	lblInsert := cg.getLabel("sm_add")
	lblFound := cg.getLabel("sm_found")
	lblLeft := cg.getLabel("sm_left")
	lblRight := cg.getLabel("sm_right")
	lblNewNode := cg.getLabel("sm_new")
	lblDone := cg.getLabel("sm_done")

	// Initialize parent to NULL (needed if root is NULL)
	cg.textSection.WriteString("    xorq %r13, %r13\n")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNewNode))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblInsert))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq %rcx, %r13\n")
	cg.textSection.WriteString("    movq $1, %r14\n")
	cg.textSection.WriteString("    movq 24(%rcx), %rcx\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblInsert))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblNewNode))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq %rcx, %r13\n")
	cg.textSection.WriteString("    xorq %r14, %r14\n")
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lblInsert))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNewNode))
	// Allocate new map node - save registers
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.textSection.WriteString("    pushq %r12\n")
	cg.textSection.WriteString("    pushq %r13\n")
	cg.textSection.WriteString("    pushq %r14\n")
	cg.textSection.WriteString("    pushq %r15\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%rsi\n", bstMapNodeSize))
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    popq %r15\n")
	cg.textSection.WriteString("    popq %r14\n")
	cg.textSection.WriteString("    popq %r13\n")
	cg.textSection.WriteString("    popq %r12\n")
	cg.textSection.WriteString("    popq %rbx\n")
	cg.textSection.WriteString("    movq %r12, (%rax)\n")
	cg.textSection.WriteString("    movq %r15, 8(%rax)\n")
	cg.textSection.WriteString("    movq $0, 16(%rax)\n")
	cg.textSection.WriteString("    movq $0, 24(%rax)\n")

	cg.textSection.WriteString("    testq %r13, %r13\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s_root\n", lblNewNode))
	cg.textSection.WriteString("    testq %r14, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s_linkr\n", lblNewNode))
	cg.textSection.WriteString("    movq %rax, 16(%r13)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_inc\n", lblNewNode))
	cg.textSection.WriteString(fmt.Sprintf("%s_linkr:\n", lblNewNode))
	cg.textSection.WriteString("    movq %rax, 24(%r13)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_inc\n", lblNewNode))
	cg.textSection.WriteString(fmt.Sprintf("%s_root:\n", lblNewNode))
	cg.textSection.WriteString("    movq %rax, (%rbx)\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_inc:\n", lblNewNode))
	cg.textSection.WriteString("    incq 8(%rbx)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblDone))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq %r15, 8(%rcx)\n") // update value

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsSortedmapIntGet retrieves value for key
func generateCollectionsSortedmapIntGet(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    popq %rbx\n")

	lblSearch := cg.getLabel("sm_get")
	lblLeft := cg.getLabel("sm_getl")
	lblRight := cg.getLabel("sm_getr")
	lblFound := cg.getLabel("sm_getf")
	lblNotFound := cg.getLabel("sm_getnf")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNotFound))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq 24(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq 8(%rcx), %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedmapIntContains checks if key exists
func generateCollectionsSortedmapIntContains(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    popq %rbx\n")

	lblSearch := cg.getLabel("sm_has")
	lblLeft := cg.getLabel("sm_hasl")
	lblRight := cg.getLabel("sm_hasr")
	lblFound := cg.getLabel("sm_hasf")
	lblNotFound := cg.getLabel("sm_hasnf")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblNotFound))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq 24(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblNotFound))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedmapIntRemove removes a key (simplified)
func generateCollectionsSortedmapIntRemove(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.generateExpressionToReg(args[1], "r12")
	cg.textSection.WriteString("    popq %rbx\n")

	lblSearch := cg.getLabel("sm_rm")
	lblLeft := cg.getLabel("sm_rml")
	lblRight := cg.getLabel("sm_rmr")
	lblFound := cg.getLabel("sm_rmf")
	lblDone := cg.getLabel("sm_rmd")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSearch))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    cmpq %r12, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lblFound))
	cg.textSection.WriteString(fmt.Sprintf("    jg %s\n", lblLeft))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRight))
	cg.textSection.WriteString("    movq 24(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLeft))
	cg.textSection.WriteString("    movq 16(%rcx), %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblSearch))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFound))
	cg.textSection.WriteString("    movabsq $-9223372036854775808, %rax\n")
	cg.textSection.WriteString("    movq %rax, (%rcx)\n")
	cg.textSection.WriteString("    decq 8(%rbx)\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
}

// generateCollectionsSortedmapIntMinKey returns the minimum key
func generateCollectionsSortedmapIntMinKey(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")

	lblLoop := cg.getLabel("sm_min")
	lblDone := cg.getLabel("sm_min_done")
	lblEmpty := cg.getLabel("sm_min_empty")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblEmpty))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    movq 16(%rcx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblEmpty))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedmapIntMaxKey returns the maximum key
func generateCollectionsSortedmapIntMaxKey(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")

	lblLoop := cg.getLabel("sm_max")
	lblDone := cg.getLabel("sm_max_done")
	lblEmpty := cg.getLabel("sm_max_empty")

	cg.textSection.WriteString("    movq (%rbx), %rcx\n")
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblEmpty))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblLoop))
	cg.textSection.WriteString("    movq 24(%rcx), %rax\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", lblDone))
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblDone))
	cg.textSection.WriteString("    movq (%rcx), %rax\n")
	cg.textSection.WriteString("    jmp 99f\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblEmpty))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("99:\n")
}

// generateCollectionsSortedmapIntLen returns the count
func generateCollectionsSortedmapIntLen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")
	cg.textSection.WriteString("    movq 8(%rbx), %rax\n")
}

// generateCollectionsSortedmapIntFree frees the map header
func generateCollectionsSortedmapIntFree(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi")
	cg.textSection.WriteString("    movq $16, %rsi\n")
	cg.textSection.WriteString("    movq $11, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// Binary search helper (int) expects args: base pointer, length, target
func generateCollectionsBinarySearchInt(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}
	cg.generateExpressionToReg(args[0], "rbx")         // base pointer
	cg.generateExpressionToReg(args[1], "rcx")         // len
	cg.generateExpressionToReg(args[2], "rdx")         // target
	cg.textSection.WriteString("    xorq %r8, %r8\n")  // low
	cg.textSection.WriteString("    movq %rcx, %r9\n") // high
	cg.textSection.WriteString("1:  cmpq %r9, %r8\n")
	cg.textSection.WriteString("    jge 4f\n")
	cg.textSection.WriteString("    movq %r8, %r10\n")
	cg.textSection.WriteString("    addq %r9, %r10\n")
	cg.textSection.WriteString("    shrq $1, %r10\n") // mid
	cg.textSection.WriteString("    movq (%rbx,%r10,8), %r11\n")
	cg.textSection.WriteString("    cmpq %rdx, %r11\n")
	cg.textSection.WriteString("    je 3f\n")
	cg.textSection.WriteString("    jl 2f\n")
	cg.textSection.WriteString("    movq %r10, %r9\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("2:  inc %r10\n")
	cg.textSection.WriteString("    movq %r10, %r8\n")
	cg.textSection.WriteString("    jmp 1b\n")
	cg.textSection.WriteString("3:  movq %r10, %rax\n")
	cg.textSection.WriteString("    jmp 5f\n")
	cg.textSection.WriteString("4:  movq $-1, %rax\n")
	cg.textSection.WriteString("5:\n")
}

// ============================================================================
// Hash module implementations
// ============================================================================

// generateHashCRC32 computes CRC32 checksum (IEEE 802.3 polynomial)
// Args: data_ptr, len -> returns uint32 in rax
func generateHashCRC32(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// CRC32 table-driven implementation
	// Polynomial: 0xEDB88320 (IEEE 802.3)
	tableLbl := cg.getLabel("crc32_table")
	loopLbl := cg.getLabel("crc32_loop")
	endLbl := cg.getLabel("crc32_end")

	// Generate CRC32 lookup table in .data section
	if !cg.hasLabel(tableLbl) {
		cg.dataSection.WriteString(fmt.Sprintf("%s:\n", tableLbl))
		poly := uint32(0xEDB88320)
		for i := 0; i < 256; i++ {
			crc := uint32(i)
			for j := 0; j < 8; j++ {
				if crc&1 != 0 {
					crc = (crc >> 1) ^ poly
				} else {
					crc >>= 1
				}
			}
			cg.dataSection.WriteString(fmt.Sprintf("    .long 0x%08x\n", crc))
		}
		cg.markLabel(tableLbl)
	}

	cg.generateExpressionToReg(args[0], "rsi")                 // data pointer
	cg.generateExpressionToReg(args[1], "rcx")                 // length
	cg.textSection.WriteString("    movl $0xFFFFFFFF, %eax\n") // crc = ~0
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rdi\n", tableLbl))

	// Main loop
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n")      // byte
	cg.textSection.WriteString("    xorb %al, %bl\n")            // index = (crc ^ byte) & 0xFF
	cg.textSection.WriteString("    shrl $8, %eax\n")            // crc >>= 8
	cg.textSection.WriteString("    movl (%rdi,%rbx,4), %edx\n") // table[index]
	cg.textSection.WriteString("    xorl %edx, %eax\n")          // crc ^= table[index]
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
	cg.textSection.WriteString("    notl %eax\n") // final XOR with 0xFFFFFFFF
}

// generateHashFNV1a computes FNV-1a hash (64-bit)
// Args: data_ptr, len -> returns uint64 in rax
func generateHashFNV1a(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	loopLbl := cg.getLabel("fnv1a_loop")
	endLbl := cg.getLabel("fnv1a_end")

	cg.generateExpressionToReg(args[0], "rsi") // data pointer
	cg.generateExpressionToReg(args[1], "rcx") // length

	// FNV-1a 64-bit offset basis and prime
	cg.textSection.WriteString("    movabsq $0xCBF29CE484222325, %rax\n") // offset_basis
	cg.textSection.WriteString("    movabsq $0x100000001B3, %r8\n")       // FNV prime

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n") // byte
	cg.textSection.WriteString("    xorq %rbx, %rax\n")     // hash ^= byte
	cg.textSection.WriteString("    imulq %r8, %rax\n")     // hash *= prime
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
}

// generateHashDJB2 computes DJB2 hash (simple string hash)
// Args: string_ptr -> returns uint64 in rax
func generateHashDJB2(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	loopLbl := cg.getLabel("djb2_loop")
	endLbl := cg.getLabel("djb2_end")

	cg.generateExpressionToReg(args[0], "rsi")           // string pointer
	cg.textSection.WriteString("    movq $5381, %rax\n") // hash = 5381

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n") // byte
	cg.textSection.WriteString("    testq %rbx, %rbx\n")    // check for null terminator
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	// hash = hash * 33 + c
	cg.textSection.WriteString("    movq %rax, %rdx\n") // save hash
	cg.textSection.WriteString("    shlq $5, %rax\n")   // hash << 5 (hash * 32)
	cg.textSection.WriteString("    addq %rdx, %rax\n") // + hash (now hash * 33)
	cg.textSection.WriteString("    addq %rbx, %rax\n") // + c
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
}

// generateHashMurmur3 computes MurmurHash3 (32-bit)
// Args: data_ptr, len, seed -> returns uint32 in rax
func generateHashMurmur3(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Simplified MurmurHash3 implementation
	// For full implementation, we'd need proper block processing
	loopLbl := cg.getLabel("murmur_loop")
	endLbl := cg.getLabel("murmur_end")

	cg.generateExpressionToReg(args[0], "rsi") // data pointer
	cg.generateExpressionToReg(args[1], "rcx") // length
	cg.generateExpressionToReg(args[2], "rax") // seed

	cg.textSection.WriteString("    movl %eax, %r8d\n") // hash = seed

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", loopLbl))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s\n", endLbl))
	cg.textSection.WriteString("    movzbq (%rsi), %rbx\n")
	cg.textSection.WriteString("    imull $0xCC9E2D51, %ebx, %ebx\n")
	cg.textSection.WriteString("    roll $15, %ebx\n")
	cg.textSection.WriteString("    imull $0x1B873593, %ebx, %ebx\n")
	cg.textSection.WriteString("    xorl %ebx, %r8d\n")
	cg.textSection.WriteString("    roll $13, %r8d\n")
	cg.textSection.WriteString("    imull $5, %r8d, %r8d\n")
	cg.textSection.WriteString("    addl $0xE6546B64, %r8d\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", loopLbl))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", endLbl))
	// Final mix
	cg.textSection.WriteString("    xorl %ecx, %r8d\n")
	cg.textSection.WriteString("    movl %r8d, %eax\n")
	cg.textSection.WriteString("    shrl $16, %eax\n")
	cg.textSection.WriteString("    xorl %r8d, %eax\n")
	cg.textSection.WriteString("    imull $0x85EBCA6B, %eax, %eax\n")
	cg.textSection.WriteString("    movl %eax, %edx\n")
	cg.textSection.WriteString("    shrl $13, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    imull $0xC2B2AE35, %eax, %eax\n")
	cg.textSection.WriteString("    movl %eax, %edx\n")
	cg.textSection.WriteString("    shrl $16, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
}

// generateHashSHA256 computes SHA-256 hash
// Args: data_ptr, len, out_buf (32 bytes) -> void
// Implements the full SHA-256 algorithm per FIPS 180-4
func generateHashSHA256(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}

	// SHA-256 requires significant state - we'll use stack space
	// Stack layout:
	//   [rsp+0..255]   - message schedule W[64] (256 bytes for 64 u32s)
	//   [rsp+256..287] - working variables a-h (32 bytes for 8 u32s)
	//   [rsp+288..319] - initial hash H[8] (32 bytes)
	//   [rsp+320]      - data_ptr
	//   [rsp+328]      - data_len
	//   [rsp+336]      - out_ptr
	//   [rsp+344]      - bytes_processed
	//   Total: 352 bytes, align to 368 for 16-byte alignment

	// Save arguments
	cg.generateExpressionToReg(args[0], "r12") // data_ptr
	cg.generateExpressionToReg(args[1], "r13") // data_len
	cg.generateExpressionToReg(args[2], "r14") // out_buf

	// Allocate stack space
	cg.textSection.WriteString("    subq $368, %rsp\n")
	cg.textSection.WriteString("    movq %r12, 320(%rsp)\n") // save data_ptr
	cg.textSection.WriteString("    movq %r13, 328(%rsp)\n") // save data_len
	cg.textSection.WriteString("    movq %r14, 336(%rsp)\n") // save out_ptr
	cg.textSection.WriteString("    movq $0, 344(%rsp)\n")   // bytes_processed = 0

	// Initialize hash values H[0..7] (SHA-256 initial values)
	// These are the first 32 bits of the fractional parts of the square roots of the first 8 primes
	cg.textSection.WriteString("    movl $0x6a09e667, 288(%rsp)\n") // H[0]
	cg.textSection.WriteString("    movl $0xbb67ae85, 292(%rsp)\n") // H[1]
	cg.textSection.WriteString("    movl $0x3c6ef372, 296(%rsp)\n") // H[2]
	cg.textSection.WriteString("    movl $0xa54ff53a, 300(%rsp)\n") // H[3]
	cg.textSection.WriteString("    movl $0x510e527f, 304(%rsp)\n") // H[4]
	cg.textSection.WriteString("    movl $0x9b05688c, 308(%rsp)\n") // H[5]
	cg.textSection.WriteString("    movl $0x1f83d9ab, 312(%rsp)\n") // H[6]
	cg.textSection.WriteString("    movl $0x5be0cd19, 316(%rsp)\n") // H[7]

	// Emit the K constants table in data section
	lblK := cg.getLabel("sha256_k")
	cg.dataSection.WriteString(fmt.Sprintf("%s:\n", lblK))
	sha256K := []uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
	}
	for _, k := range sha256K {
		cg.dataSection.WriteString(fmt.Sprintf("    .long 0x%08x\n", k))
	}

	// Process complete 64-byte blocks
	lblBlockLoop := cg.getLabel("sha256_block")
	lblPadding := cg.getLabel("sha256_pad")
	lblCompress := cg.getLabel("sha256_compress")
	lblRoundLoop := cg.getLabel("sha256_round")
	lblSchedule := cg.getLabel("sha256_schedule")
	lblFinalize := cg.getLabel("sha256_final")

	// Main block processing loop
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblBlockLoop))
	cg.textSection.WriteString("    movq 328(%rsp), %rax\n") // remaining data_len
	cg.textSection.WriteString("    movq 344(%rsp), %rbx\n") // bytes_processed
	cg.textSection.WriteString("    subq %rbx, %rax\n")      // remaining = len - processed
	cg.textSection.WriteString("    cmpq $64, %rax\n")       // at least 64 bytes?
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblPadding))

	// Copy 64 bytes to message schedule W[0..15] with big-endian conversion
	cg.textSection.WriteString("    movq 320(%rsp), %rsi\n") // data_ptr
	cg.textSection.WriteString("    addq 344(%rsp), %rsi\n") // + bytes_processed
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")    // W array
	cg.textSection.WriteString("    movq $16, %rcx\n")       // 16 words
	cg.textSection.WriteString(fmt.Sprintf("%s_copy:\n", lblBlockLoop))
	cg.textSection.WriteString("    movl (%rsi), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n") // big-endian conversion
	cg.textSection.WriteString("    movl %eax, (%rdi)\n")
	cg.textSection.WriteString("    addq $4, %rsi\n")
	cg.textSection.WriteString("    addq $4, %rdi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s_copy\n", lblBlockLoop))

	cg.textSection.WriteString("    addq $64, 344(%rsp)\n") // bytes_processed += 64
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCompress))

	// Padding phase - handle final block(s)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblPadding))
	// Create padded block on stack at W
	// Clear W[0..15]
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")
	cg.textSection.WriteString("    movq $64, %rcx\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")

	// Copy remaining bytes
	cg.textSection.WriteString("    movq 328(%rsp), %rax\n") // total len
	cg.textSection.WriteString("    movq 344(%rsp), %rbx\n") // processed
	cg.textSection.WriteString("    subq %rbx, %rax\n")      // remaining
	cg.textSection.WriteString("    movq %rax, %rcx\n")      // save remaining count
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s_append1\n", lblPadding))

	cg.textSection.WriteString("    movq 320(%rsp), %rsi\n")
	cg.textSection.WriteString("    addq %rbx, %rsi\n")   // source = data + processed
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n") // dest = W
	cg.textSection.WriteString("    rep movsb\n")         // copy remaining bytes

	// Append 0x80 byte
	cg.textSection.WriteString(fmt.Sprintf("%s_append1:\n", lblPadding))
	cg.textSection.WriteString("    movq 328(%rsp), %rax\n")    // total len
	cg.textSection.WriteString("    movq 344(%rsp), %rbx\n")    // processed
	cg.textSection.WriteString("    subq %rbx, %rax\n")         // remaining
	cg.textSection.WriteString("    movb $0x80, (%rsp,%rax)\n") // append 0x80

	// Check if we have room for length (need 8 bytes at end)
	// If remaining+1 > 56, we need two blocks
	cg.textSection.WriteString("    incq %rax\n") // remaining + 1
	cg.textSection.WriteString("    cmpq $56, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    ja %s_twoblock\n", lblPadding))

	// Single block: append length at bytes 56-63
	cg.textSection.WriteString("    movq 328(%rsp), %rax\n") // total length in bytes
	cg.textSection.WriteString("    shlq $3, %rax\n")        // convert to bits
	cg.textSection.WriteString("    bswap %rax\n")           // big-endian
	cg.textSection.WriteString("    movq %rax, 56(%rsp)\n")  // store at offset 56

	// Convert W[0..15] to big-endian and compress
	cg.textSection.WriteString("    leaq (%rsp), %rsi\n")
	cg.textSection.WriteString("    movq $16, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_bswap:\n", lblPadding))
	cg.textSection.WriteString("    movl (%rsi), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, (%rsi)\n")
	cg.textSection.WriteString("    addq $4, %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s_bswap\n", lblPadding))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCompress))

	// Two blocks needed
	cg.textSection.WriteString(fmt.Sprintf("%s_twoblock:\n", lblPadding))
	// First: compress current block (already has 0x80)
	cg.textSection.WriteString("    leaq (%rsp), %rsi\n")
	cg.textSection.WriteString("    movq $16, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_bswap2:\n", lblPadding))
	cg.textSection.WriteString("    movl (%rsi), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, (%rsi)\n")
	cg.textSection.WriteString("    addq $4, %rsi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s_bswap2\n", lblPadding))

	// Compress first padded block
	cg.textSection.WriteString(fmt.Sprintf("    call %s_func\n", lblCompress))

	// Create second block with just length
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")
	cg.textSection.WriteString("    movq $64, %rcx\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
	cg.textSection.WriteString("    movq 328(%rsp), %rax\n")
	cg.textSection.WriteString("    shlq $3, %rax\n")
	cg.textSection.WriteString("    bswap %rax\n")
	cg.textSection.WriteString("    movq %rax, 56(%rsp)\n")
	// W is already in memory format (zeros + length)
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_final\n", lblCompress))

	// Compression function
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCompress))
	cg.textSection.WriteString(fmt.Sprintf("%s_func:\n", lblCompress))

	// Extend W[0..15] to W[0..63]
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblSchedule))
	cg.textSection.WriteString("    movq $16, %rcx\n") // i = 16
	cg.textSection.WriteString(fmt.Sprintf("%s_ext:\n", lblSchedule))
	cg.textSection.WriteString("    cmpq $64, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s_ext_done\n", lblSchedule))

	// W[i] = sig1(W[i-2]) + W[i-7] + sig0(W[i-15]) + W[i-16]
	// sig0(x) = ROTR7(x) ^ ROTR18(x) ^ SHR3(x)
	// sig1(x) = ROTR17(x) ^ ROTR19(x) ^ SHR10(x)

	// Get W[i-2] for sig1
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    subq $2, %rax\n")
	cg.textSection.WriteString("    movl (%rsp,%rax,4), %r8d\n") // W[i-2]

	// sig1(W[i-2]): ROTR17 ^ ROTR19 ^ SHR10
	cg.textSection.WriteString("    movl %r8d, %eax\n")
	cg.textSection.WriteString("    rorl $17, %eax\n")
	cg.textSection.WriteString("    movl %r8d, %edx\n")
	cg.textSection.WriteString("    rorl $19, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    movl %r8d, %edx\n")
	cg.textSection.WriteString("    shrl $10, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    movl %eax, %r9d\n") // r9d = sig1

	// Get W[i-7]
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    subq $7, %rax\n")
	cg.textSection.WriteString("    addl (%rsp,%rax,4), %r9d\n") // += W[i-7]

	// Get W[i-15] for sig0
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    subq $15, %rax\n")
	cg.textSection.WriteString("    movl (%rsp,%rax,4), %r8d\n") // W[i-15]

	// sig0(W[i-15]): ROTR7 ^ ROTR18 ^ SHR3
	cg.textSection.WriteString("    movl %r8d, %eax\n")
	cg.textSection.WriteString("    rorl $7, %eax\n")
	cg.textSection.WriteString("    movl %r8d, %edx\n")
	cg.textSection.WriteString("    rorl $18, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    movl %r8d, %edx\n")
	cg.textSection.WriteString("    shrl $3, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    addl %eax, %r9d\n") // += sig0

	// Get W[i-16]
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    subq $16, %rax\n")
	cg.textSection.WriteString("    addl (%rsp,%rax,4), %r9d\n") // += W[i-16]

	// Store W[i]
	cg.textSection.WriteString("    movl %r9d, (%rsp,%rcx,4)\n")
	cg.textSection.WriteString("    incq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_ext\n", lblSchedule))
	cg.textSection.WriteString(fmt.Sprintf("%s_ext_done:\n", lblSchedule))

	// Initialize working variables from current hash
	cg.textSection.WriteString("    movl 288(%rsp), %r8d\n")  // a = H[0]
	cg.textSection.WriteString("    movl 292(%rsp), %r9d\n")  // b = H[1]
	cg.textSection.WriteString("    movl 296(%rsp), %r10d\n") // c = H[2]
	cg.textSection.WriteString("    movl 300(%rsp), %r11d\n") // d = H[3]
	cg.textSection.WriteString("    movl 304(%rsp), %r12d\n") // e = H[4]
	cg.textSection.WriteString("    movl 308(%rsp), %r13d\n") // f = H[5]
	cg.textSection.WriteString("    movl 312(%rsp), %r14d\n") // g = H[6]
	cg.textSection.WriteString("    movl 316(%rsp), %r15d\n") // h = H[7]

	// 64 rounds
	cg.textSection.WriteString("    xorq %rbx, %rbx\n") // i = 0
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRoundLoop))
	cg.textSection.WriteString("    cmpq $64, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s_done\n", lblRoundLoop))

	// T1 = h + Sig1(e) + Ch(e,f,g) + K[i] + W[i]
	// Sig1(e) = ROTR6(e) ^ ROTR11(e) ^ ROTR25(e)
	// Ch(e,f,g) = (e & f) ^ (~e & g)

	// Sig1(e)
	cg.textSection.WriteString("    movl %r12d, %eax\n") // e
	cg.textSection.WriteString("    rorl $6, %eax\n")
	cg.textSection.WriteString("    movl %r12d, %edx\n")
	cg.textSection.WriteString("    rorl $11, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    movl %r12d, %edx\n")
	cg.textSection.WriteString("    rorl $25, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n") // eax = Sig1(e)

	// Add h
	cg.textSection.WriteString("    addl %r15d, %eax\n")

	// Ch(e,f,g) = (e & f) ^ (~e & g)
	cg.textSection.WriteString("    movl %r12d, %edx\n") // e
	cg.textSection.WriteString("    andl %r13d, %edx\n") // e & f
	cg.textSection.WriteString("    movl %r12d, %ecx\n") // e
	cg.textSection.WriteString("    notl %ecx\n")        // ~e
	cg.textSection.WriteString("    andl %r14d, %ecx\n") // ~e & g
	cg.textSection.WriteString("    xorl %ecx, %edx\n")  // Ch
	cg.textSection.WriteString("    addl %edx, %eax\n")

	// Add K[i]
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rdi\n", lblK))
	cg.textSection.WriteString("    addl (%rdi,%rbx,4), %eax\n")

	// Add W[i]
	cg.textSection.WriteString("    addl (%rsp,%rbx,4), %eax\n")
	cg.textSection.WriteString("    movl %eax, %edi\n") // edi = T1

	// T2 = Sig0(a) + Maj(a,b,c)
	// Sig0(a) = ROTR2(a) ^ ROTR13(a) ^ ROTR22(a)
	// Maj(a,b,c) = (a & b) ^ (a & c) ^ (b & c)

	cg.textSection.WriteString("    movl %r8d, %eax\n") // a
	cg.textSection.WriteString("    rorl $2, %eax\n")
	cg.textSection.WriteString("    movl %r8d, %edx\n")
	cg.textSection.WriteString("    rorl $13, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n")
	cg.textSection.WriteString("    movl %r8d, %edx\n")
	cg.textSection.WriteString("    rorl $22, %edx\n")
	cg.textSection.WriteString("    xorl %edx, %eax\n") // eax = Sig0(a)

	// Maj(a,b,c)
	cg.textSection.WriteString("    movl %r8d, %edx\n")  // a
	cg.textSection.WriteString("    andl %r9d, %edx\n")  // a & b
	cg.textSection.WriteString("    movl %r8d, %ecx\n")  // a
	cg.textSection.WriteString("    andl %r10d, %ecx\n") // a & c
	cg.textSection.WriteString("    xorl %ecx, %edx\n")
	cg.textSection.WriteString("    movl %r9d, %ecx\n")  // b
	cg.textSection.WriteString("    andl %r10d, %ecx\n") // b & c
	cg.textSection.WriteString("    xorl %ecx, %edx\n")  // Maj
	cg.textSection.WriteString("    addl %edx, %eax\n")  // eax = T2

	// Update working variables
	cg.textSection.WriteString("    movl %r14d, %r15d\n") // h = g
	cg.textSection.WriteString("    movl %r13d, %r14d\n") // g = f
	cg.textSection.WriteString("    movl %r12d, %r13d\n") // f = e
	cg.textSection.WriteString("    movl %r11d, %r12d\n") // e = d + T1
	cg.textSection.WriteString("    addl %edi, %r12d\n")
	cg.textSection.WriteString("    movl %r10d, %r11d\n") // d = c
	cg.textSection.WriteString("    movl %r9d, %r10d\n")  // c = b
	cg.textSection.WriteString("    movl %r8d, %r9d\n")   // b = a
	cg.textSection.WriteString("    movl %edi, %r8d\n")   // a = T1 + T2
	cg.textSection.WriteString("    addl %eax, %r8d\n")

	cg.textSection.WriteString("    incq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblRoundLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s_done:\n", lblRoundLoop))

	// Add compressed chunk to current hash
	cg.textSection.WriteString("    addl %r8d, 288(%rsp)\n")
	cg.textSection.WriteString("    addl %r9d, 292(%rsp)\n")
	cg.textSection.WriteString("    addl %r10d, 296(%rsp)\n")
	cg.textSection.WriteString("    addl %r11d, 300(%rsp)\n")
	cg.textSection.WriteString("    addl %r12d, 304(%rsp)\n")
	cg.textSection.WriteString("    addl %r13d, 308(%rsp)\n")
	cg.textSection.WriteString("    addl %r14d, 312(%rsp)\n")
	cg.textSection.WriteString("    addl %r15d, 316(%rsp)\n")

	cg.textSection.WriteString("    ret\n")
	cg.textSection.WriteString(fmt.Sprintf("%s_final:\n", lblCompress))

	// Final compression for last padded block
	cg.textSection.WriteString(fmt.Sprintf("    call %s_func\n", lblCompress))

	// Finalize - copy hash to output buffer
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFinalize))
	cg.textSection.WriteString("    movq 336(%rsp), %rdi\n") // out_ptr
	// Copy H[0..7] to output with big-endian conversion
	cg.textSection.WriteString("    movl 288(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, (%rdi)\n")
	cg.textSection.WriteString("    movl 292(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 4(%rdi)\n")
	cg.textSection.WriteString("    movl 296(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 8(%rdi)\n")
	cg.textSection.WriteString("    movl 300(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 12(%rdi)\n")
	cg.textSection.WriteString("    movl 304(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 16(%rdi)\n")
	cg.textSection.WriteString("    movl 308(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 20(%rdi)\n")
	cg.textSection.WriteString("    movl 312(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 24(%rdi)\n")
	cg.textSection.WriteString("    movl 316(%rsp), %eax\n")
	cg.textSection.WriteString("    bswap %eax\n")
	cg.textSection.WriteString("    movl %eax, 28(%rdi)\n")

	// Restore stack
	cg.textSection.WriteString("    addq $368, %rsp\n")
}

// generateHashMD5 computes MD5 hash
// Args: data_ptr, len, out_buf (16 bytes) -> void
// Implements the full MD5 algorithm per RFC 1321
func generateHashMD5(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		return
	}

	// Stack layout for MD5:
	//   [rsp+0..63]    - message block M[16] (64 bytes)
	//   [rsp+64..79]   - hash state A,B,C,D (16 bytes)
	//   [rsp+80]       - data_ptr
	//   [rsp+88]       - data_len
	//   [rsp+96]       - out_ptr
	//   [rsp+104]      - bytes_processed
	//   Total: 112 bytes, align to 128 for 16-byte alignment

	// Save arguments
	cg.generateExpressionToReg(args[0], "r12") // data_ptr
	cg.generateExpressionToReg(args[1], "r13") // data_len
	cg.generateExpressionToReg(args[2], "r14") // out_buf

	// Allocate stack space
	cg.textSection.WriteString("    subq $128, %rsp\n")
	cg.textSection.WriteString("    movq %r12, 80(%rsp)\n") // save data_ptr
	cg.textSection.WriteString("    movq %r13, 88(%rsp)\n") // save data_len
	cg.textSection.WriteString("    movq %r14, 96(%rsp)\n") // save out_ptr
	cg.textSection.WriteString("    movq $0, 104(%rsp)\n")  // bytes_processed = 0

	// Initialize hash state (MD5 magic constants - little endian)
	cg.textSection.WriteString("    movl $0x67452301, 64(%rsp)\n") // A
	cg.textSection.WriteString("    movl $0xefcdab89, 68(%rsp)\n") // B
	cg.textSection.WriteString("    movl $0x98badcfe, 72(%rsp)\n") // C
	cg.textSection.WriteString("    movl $0x10325476, 76(%rsp)\n") // D

	// Emit T table (precomputed sin values) in data section
	lblT := cg.getLabel("md5_t")
	cg.dataSection.WriteString(fmt.Sprintf("%s:\n", lblT))
	md5T := []uint32{
		0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee,
		0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
		0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
		0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
		0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa,
		0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
		0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed,
		0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
		0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
		0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
		0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05,
		0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
		0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039,
		0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
		0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
		0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
	}
	for _, t := range md5T {
		cg.dataSection.WriteString(fmt.Sprintf("    .long 0x%08x\n", t))
	}

	// Shift amounts per round
	lblS := cg.getLabel("md5_s")
	cg.dataSection.WriteString(fmt.Sprintf("%s:\n", lblS))
	md5S := []byte{
		7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
		5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
		4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
		6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
	}
	for i := 0; i < len(md5S); i += 16 {
		cg.dataSection.WriteString("    .byte ")
		for j := 0; j < 16; j++ {
			if j > 0 {
				cg.dataSection.WriteString(", ")
			}
			cg.dataSection.WriteString(fmt.Sprintf("%d", md5S[i+j]))
		}
		cg.dataSection.WriteString("\n")
	}

	lblBlockLoop := cg.getLabel("md5_block")
	lblPadding := cg.getLabel("md5_pad")
	lblCompress := cg.getLabel("md5_compress")
	lblRound := cg.getLabel("md5_round")
	lblFinalize := cg.getLabel("md5_final")

	// Main block processing loop
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblBlockLoop))
	cg.textSection.WriteString("    movq 88(%rsp), %rax\n")  // data_len
	cg.textSection.WriteString("    movq 104(%rsp), %rbx\n") // bytes_processed
	cg.textSection.WriteString("    subq %rbx, %rax\n")      // remaining
	cg.textSection.WriteString("    cmpq $64, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lblPadding))

	// Copy 64 bytes to M (no endian conversion - MD5 is little-endian)
	cg.textSection.WriteString("    movq 80(%rsp), %rsi\n")
	cg.textSection.WriteString("    addq 104(%rsp), %rsi\n")
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")
	cg.textSection.WriteString("    movq $64, %rcx\n")
	cg.textSection.WriteString("    rep movsb\n")

	cg.textSection.WriteString("    addq $64, 104(%rsp)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCompress))

	// Padding phase
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblPadding))
	// Clear M
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")
	cg.textSection.WriteString("    movq $64, %rcx\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")

	// Copy remaining bytes
	cg.textSection.WriteString("    movq 88(%rsp), %rax\n")
	cg.textSection.WriteString("    movq 104(%rsp), %rbx\n")
	cg.textSection.WriteString("    subq %rbx, %rax\n") // remaining
	cg.textSection.WriteString("    movq %rax, %rcx\n")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jz %s_append\n", lblPadding))

	cg.textSection.WriteString("    movq 80(%rsp), %rsi\n")
	cg.textSection.WriteString("    addq %rbx, %rsi\n")
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")
	cg.textSection.WriteString("    rep movsb\n")

	// Append 0x80
	cg.textSection.WriteString(fmt.Sprintf("%s_append:\n", lblPadding))
	cg.textSection.WriteString("    movq 88(%rsp), %rax\n")
	cg.textSection.WriteString("    movq 104(%rsp), %rbx\n")
	cg.textSection.WriteString("    subq %rbx, %rax\n")
	cg.textSection.WriteString("    movb $0x80, (%rsp,%rax)\n")

	// Check if room for length (need 8 bytes at offset 56)
	cg.textSection.WriteString("    incq %rax\n")
	cg.textSection.WriteString("    cmpq $56, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    ja %s_twoblock\n", lblPadding))

	// Single block - append length in bits (little-endian)
	cg.textSection.WriteString("    movq 88(%rsp), %rax\n")
	cg.textSection.WriteString("    shlq $3, %rax\n")       // bits
	cg.textSection.WriteString("    movq %rax, 56(%rsp)\n") // little-endian length
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblCompress))

	// Two blocks needed
	cg.textSection.WriteString(fmt.Sprintf("%s_twoblock:\n", lblPadding))
	// Compress first padded block
	cg.textSection.WriteString(fmt.Sprintf("    call %s_func\n", lblCompress))

	// Create second block with just length
	cg.textSection.WriteString("    leaq (%rsp), %rdi\n")
	cg.textSection.WriteString("    movq $64, %rcx\n")
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString("    rep stosb\n")
	cg.textSection.WriteString("    movq 88(%rsp), %rax\n")
	cg.textSection.WriteString("    shlq $3, %rax\n")
	cg.textSection.WriteString("    movq %rax, 56(%rsp)\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_final\n", lblCompress))

	// Compression function
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblCompress))
	cg.textSection.WriteString(fmt.Sprintf("%s_func:\n", lblCompress))

	// Load current hash into registers
	cg.textSection.WriteString("    movl 64(%rsp), %r8d\n")  // A
	cg.textSection.WriteString("    movl 68(%rsp), %r9d\n")  // B
	cg.textSection.WriteString("    movl 72(%rsp), %r10d\n") // C
	cg.textSection.WriteString("    movl 76(%rsp), %r11d\n") // D

	// 64 rounds
	cg.textSection.WriteString("    xorq %rbx, %rbx\n") // i = 0
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblRound))
	cg.textSection.WriteString("    cmpq $64, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s_done\n", lblRound))

	// Compute F and g based on round number
	// Round 0-15: F = (B & C) | (~B & D), g = i
	// Round 16-31: F = (D & B) | (~D & C), g = (5*i + 1) % 16
	// Round 32-47: F = B ^ C ^ D, g = (3*i + 5) % 16
	// Round 48-63: F = C ^ (B | ~D), g = (7*i) % 16

	cg.textSection.WriteString("    cmpq $16, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s_r16\n", lblRound))

	// Round 0-15: F = (B & C) | (~B & D)
	cg.textSection.WriteString("    movl %r9d, %eax\n")  // B
	cg.textSection.WriteString("    andl %r10d, %eax\n") // B & C
	cg.textSection.WriteString("    movl %r9d, %edx\n")
	cg.textSection.WriteString("    notl %edx\n")        // ~B
	cg.textSection.WriteString("    andl %r11d, %edx\n") // ~B & D
	cg.textSection.WriteString("    orl %edx, %eax\n")   // F
	cg.textSection.WriteString("    movq %rbx, %rdi\n")  // g = i
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_compute\n", lblRound))

	cg.textSection.WriteString(fmt.Sprintf("%s_r16:\n", lblRound))
	cg.textSection.WriteString("    cmpq $32, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s_r32\n", lblRound))

	// Round 16-31: F = (D & B) | (~D & C)
	cg.textSection.WriteString("    movl %r11d, %eax\n") // D
	cg.textSection.WriteString("    andl %r9d, %eax\n")  // D & B
	cg.textSection.WriteString("    movl %r11d, %edx\n")
	cg.textSection.WriteString("    notl %edx\n")        // ~D
	cg.textSection.WriteString("    andl %r10d, %edx\n") // ~D & C
	cg.textSection.WriteString("    orl %edx, %eax\n")   // F
	// g = (5*i + 1) % 16
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    imulq $5, %rdi\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    andq $15, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_compute\n", lblRound))

	cg.textSection.WriteString(fmt.Sprintf("%s_r32:\n", lblRound))
	cg.textSection.WriteString("    cmpq $48, %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jge %s_r48\n", lblRound))

	// Round 32-47: F = B ^ C ^ D
	cg.textSection.WriteString("    movl %r9d, %eax\n")
	cg.textSection.WriteString("    xorl %r10d, %eax\n")
	cg.textSection.WriteString("    xorl %r11d, %eax\n")
	// g = (3*i + 5) % 16
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    imulq $3, %rdi\n")
	cg.textSection.WriteString("    addq $5, %rdi\n")
	cg.textSection.WriteString("    andq $15, %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s_compute\n", lblRound))

	cg.textSection.WriteString(fmt.Sprintf("%s_r48:\n", lblRound))
	// Round 48-63: F = C ^ (B | ~D)
	cg.textSection.WriteString("    movl %r11d, %eax\n") // D
	cg.textSection.WriteString("    notl %eax\n")        // ~D
	cg.textSection.WriteString("    orl %r9d, %eax\n")   // B | ~D
	cg.textSection.WriteString("    xorl %r10d, %eax\n") // C ^ ...
	// g = (7*i) % 16
	cg.textSection.WriteString("    movq %rbx, %rdi\n")
	cg.textSection.WriteString("    imulq $7, %rdi\n")
	cg.textSection.WriteString("    andq $15, %rdi\n")

	cg.textSection.WriteString(fmt.Sprintf("%s_compute:\n", lblRound))
	// F is in %eax, g is in %rdi

	// F = F + A + K[i] + M[g]
	cg.textSection.WriteString("    addl %r8d, %eax\n") // + A
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", lblT))
	cg.textSection.WriteString("    addl (%rsi,%rbx,4), %eax\n") // + K[i]
	cg.textSection.WriteString("    addl (%rsp,%rdi,4), %eax\n") // + M[g]

	// Rotate left by S[i]
	cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%rsi\n", lblS))
	cg.textSection.WriteString("    movzbl (%rsi,%rbx), %ecx\n") // S[i]
	cg.textSection.WriteString("    roll %cl, %eax\n")

	// F = F + B
	cg.textSection.WriteString("    addl %r9d, %eax\n")

	// Rotate registers: A=D, D=C, C=B, B=F
	cg.textSection.WriteString("    movl %r11d, %r8d\n")  // A = D
	cg.textSection.WriteString("    movl %r10d, %r11d\n") // D = C
	cg.textSection.WriteString("    movl %r9d, %r10d\n")  // C = B
	cg.textSection.WriteString("    movl %eax, %r9d\n")   // B = F

	cg.textSection.WriteString("    incq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lblRound))

	cg.textSection.WriteString(fmt.Sprintf("%s_done:\n", lblRound))

	// Add to hash state
	cg.textSection.WriteString("    addl %r8d, 64(%rsp)\n")
	cg.textSection.WriteString("    addl %r9d, 68(%rsp)\n")
	cg.textSection.WriteString("    addl %r10d, 72(%rsp)\n")
	cg.textSection.WriteString("    addl %r11d, 76(%rsp)\n")

	cg.textSection.WriteString("    ret\n")

	cg.textSection.WriteString(fmt.Sprintf("%s_final:\n", lblCompress))
	// Compress last block
	cg.textSection.WriteString(fmt.Sprintf("    call %s_func\n", lblCompress))

	// Finalize - copy hash to output (already little-endian)
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lblFinalize))
	cg.textSection.WriteString("    movq 96(%rsp), %rdi\n")
	cg.textSection.WriteString("    movl 64(%rsp), %eax\n")
	cg.textSection.WriteString("    movl %eax, (%rdi)\n")
	cg.textSection.WriteString("    movl 68(%rsp), %eax\n")
	cg.textSection.WriteString("    movl %eax, 4(%rdi)\n")
	cg.textSection.WriteString("    movl 72(%rsp), %eax\n")
	cg.textSection.WriteString("    movl %eax, 8(%rdi)\n")
	cg.textSection.WriteString("    movl 76(%rsp), %eax\n")
	cg.textSection.WriteString("    movl %eax, 12(%rdi)\n")

	// Restore stack
	cg.textSection.WriteString("    addq $128, %rsp\n")
}

// ============================================================================
// String Extension Functions (Phase 4)
// ============================================================================

// generateStringSubstring(s, start, len) -> new allocated string
func generateStringSubstring(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Allocate buffer, copy substring from start with length len
	cg.generateExpressionToReg(args[1], "r8")  // start offset
	cg.generateExpressionToReg(args[2], "r9")  // length
	cg.generateExpressionToReg(args[0], "rbx") // source string ptr

	// size = len + 1 (for NUL)
	cg.textSection.WriteString("    movq %r9, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	// mmap
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq %rsi, %rdx\n") // size
	cg.textSection.WriteString("    movq $3, %r10\n")
	cg.textSection.WriteString("    movq $34, %rcx\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")
	// Copy substring
	cg.textSection.WriteString("    movq %rax, %rdi\n")
	cg.textSection.WriteString("    addq %r8, %rbx\n") // source + start
	cg.textSection.WriteString("    movq %r9, %rcx\n") // length
	cg.textSection.WriteString("    rep movsb\n")
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    mov %rax, %rax\n")  // result in rax
}

// generateStringSplit(str, delim) -> array ptr
// Returns a pointer to a structure: [count:i64][ptr1][ptr2]...[ptrN]
// Each pointer points to an allocated substring
func generateStringSplit(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Get source string into r10
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n")
		}
	default:
		cg.generateExpressionToReg(args[0], "r10")
	}

	// Get delimiter into r11 (single char for simplicity)
	switch v := args[1].(type) {
	case *StringLiteral:
		if len(v.Value) > 0 {
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r11\n", byte(v.Value[0])))
		} else {
			cg.textSection.WriteString("    xorq %r11, %r11\n")
		}
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
			cg.textSection.WriteString("    movzbl (%rbx), %r11d\n")
		} else {
			cg.textSection.WriteString("    xorq %r11, %r11\n")
		}
	default:
		cg.generateExpressionToReg(args[1], "rbx")
		cg.textSection.WriteString("    movzbl (%rbx), %r11d\n")
	}

	// Save string ptr and delimiter
	cg.textSection.WriteString("    pushq %r10\n") // save src
	cg.textSection.WriteString("    pushq %r11\n") // save delim

	// First pass: count delimiters to know array size
	cg.textSection.WriteString("    movq %r10, %rbx\n")
	cg.textSection.WriteString("    movq $1, %r12\n") // count = 1 (at least one segment)

	lCountLoop := cg.getLabel("split_count_loop")
	lCountDone := cg.getLabel("split_count_done")
	lCountInc := cg.getLabel("split_count_inc")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCountLoop))
	cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
	cg.textSection.WriteString("    testb %al, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCountDone))
	cg.textSection.WriteString("    cmpb %r11b, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lCountInc))
	cg.textSection.WriteString("    incq %r12\n") // found delimiter
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCountInc))
	cg.textSection.WriteString("    incq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCountLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCountDone))

	// Save count
	cg.textSection.WriteString("    pushq %r12\n")

	// Allocate array: 8 bytes for count + 8 bytes per pointer
	cg.textSection.WriteString("    movq %r12, %rsi\n")
	cg.textSection.WriteString("    shlq $3, %rsi\n") // count * 8
	cg.textSection.WriteString("    addq $8, %rsi\n") // + 8 for count header
	cg.textSection.WriteString("    movq $9, %rax\n") // mmap
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	lErr := cg.getLabel("split_err")
	lEnd := cg.getLabel("split_end")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))

	// rax = array ptr, store count at offset 0
	cg.textSection.WriteString("    movq %rax, %r15\n")   // save array ptr
	cg.textSection.WriteString("    popq %r12\n")         // restore count
	cg.textSection.WriteString("    movq %r12, (%r15)\n") // store count

	// Restore delim and src
	cg.textSection.WriteString("    popq %r11\n") // restore delim
	cg.textSection.WriteString("    popq %r10\n") // restore src

	// Second pass: extract substrings
	// r13 = current segment start, r14 = current array slot (after count)
	cg.textSection.WriteString("    movq %r10, %r13\n")    // segment start
	cg.textSection.WriteString("    leaq 8(%r15), %r14\n") // first slot
	cg.textSection.WriteString("    movq %r10, %rbx\n")    // current ptr

	lExtractLoop := cg.getLabel("split_extract_loop")
	lExtractDone := cg.getLabel("split_extract_done")
	lFoundDelim := cg.getLabel("split_found_delim")
	_ = cg.getLabel("split_extract_next") // reserved for future use

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lExtractLoop))
	cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
	cg.textSection.WriteString("    testb %al, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lFoundDelim)) // end of string = final segment
	cg.textSection.WriteString("    cmpb %r11b, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lFoundDelim))
	cg.textSection.WriteString("    incq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lExtractLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lFoundDelim))
	// Calculate segment length: rbx - r13
	cg.textSection.WriteString("    movq %rbx, %rcx\n")
	cg.textSection.WriteString("    subq %r13, %rcx\n") // len = current - start

	// Save registers before mmap
	cg.textSection.WriteString("    pushq %rbx\n")
	cg.textSection.WriteString("    pushq %r10\n")
	cg.textSection.WriteString("    pushq %r11\n")
	cg.textSection.WriteString("    pushq %r13\n")
	cg.textSection.WriteString("    pushq %r14\n")
	cg.textSection.WriteString("    pushq %r15\n")
	cg.textSection.WriteString("    pushq %rcx\n")

	// Allocate substring (len + 1)
	cg.textSection.WriteString("    movq %rcx, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	// Restore registers
	cg.textSection.WriteString("    popq %rcx\n")
	cg.textSection.WriteString("    popq %r15\n")
	cg.textSection.WriteString("    popq %r14\n")
	cg.textSection.WriteString("    popq %r13\n")
	cg.textSection.WriteString("    popq %r11\n")
	cg.textSection.WriteString("    popq %r10\n")
	cg.textSection.WriteString("    popq %rbx\n")

	// Copy segment: rax = dest, r13 = src, rcx = len
	cg.textSection.WriteString("    movq %rax, %rdi\n")
	cg.textSection.WriteString("    movq %r13, %rsi\n")
	cg.textSection.WriteString("    pushq %rcx\n")
	cg.textSection.WriteString("    rep movsb\n")
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    popq %rcx\n")

	// Store pointer in array
	cg.textSection.WriteString("    movq %rax, (%r14)\n")
	cg.textSection.WriteString("    addq $8, %r14\n") // next slot

	// Check if done
	cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
	cg.textSection.WriteString("    testb %al, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lExtractDone))

	// Move past delimiter
	cg.textSection.WriteString("    incq %rbx\n")
	cg.textSection.WriteString("    movq %rbx, %r13\n") // new segment start
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lExtractLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lExtractDone))
	cg.textSection.WriteString("    movq %r15, %rax\n") // return array ptr
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    addq $24, %rsp\n") // clean stack
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// generateStringJoin(array, sep) -> joined string
// array format: [count:i64][ptr1][ptr2]...[ptrN]
func generateStringJoin(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Get array ptr into r10
	cg.generateExpressionToReg(args[0], "r10")

	// Get separator into r11
	switch v := args[1].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r11\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", len(v.Value))) // sep len
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r11\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r9\n", l))
			} else {
				// Compute sep length at runtime
				lLoop := cg.getLabel("join_seplen_loop")
				lEnd := cg.getLabel("join_seplen_end")
				cg.textSection.WriteString("    movq %r11, %rbx\n")
				cg.textSection.WriteString("    xorq %r9, %r9\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
				cg.textSection.WriteString("    testb %al, %al\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    incq %r9\n")
				cg.textSection.WriteString("    incq %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r11, %r11\n")
			cg.textSection.WriteString("    xorq %r9, %r9\n")
		}
	default:
		cg.generateExpressionToReg(args[1], "r11")
		// Compute sep length at runtime
		lLoop := cg.getLabel("join_seplen_loop")
		lEnd := cg.getLabel("join_seplen_end")
		cg.textSection.WriteString("    movq %r11, %rbx\n")
		cg.textSection.WriteString("    xorq %r9, %r9\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
		cg.textSection.WriteString("    testb %al, %al\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    incq %r9\n")
		cg.textSection.WriteString("    incq %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	}

	// r10 = array, r11 = sep, r9 = sep_len
	// Get count from array
	cg.textSection.WriteString("    movq (%r10), %r12\n") // count

	// Calculate total length needed
	// total = sum(strlen of each) + (count-1) * sep_len
	cg.textSection.WriteString("    xorq %r8, %r8\n")      // total = 0
	cg.textSection.WriteString("    leaq 8(%r10), %r13\n") // first ptr slot
	cg.textSection.WriteString("    movq %r12, %r14\n")    // loop counter

	lLenLoop := cg.getLabel("join_len_loop")
	lLenDone := cg.getLabel("join_len_done")
	lStrLen := cg.getLabel("join_strlen")
	lStrLenDone := cg.getLabel("join_strlen_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLenLoop))
	cg.textSection.WriteString("    testq %r14, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lLenDone))
	cg.textSection.WriteString("    movq (%r13), %rbx\n") // get string ptr
	// strlen loop
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lStrLen))
	cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
	cg.textSection.WriteString("    testb %al, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lStrLenDone))
	cg.textSection.WriteString("    incq %r8\n")
	cg.textSection.WriteString("    incq %rbx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lStrLen))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lStrLenDone))
	cg.textSection.WriteString("    addq $8, %r13\n") // next slot
	cg.textSection.WriteString("    decq %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLenLoop))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLenDone))

	// Add separator lengths: (count-1) * sep_len
	cg.textSection.WriteString("    movq %r12, %rax\n")
	cg.textSection.WriteString("    decq %rax\n")       // count - 1
	cg.textSection.WriteString("    imulq %r9, %rax\n") // * sep_len
	cg.textSection.WriteString("    addq %rax, %r8\n")  // add to total

	// Save values before mmap
	cg.textSection.WriteString("    pushq %r10\n") // array
	cg.textSection.WriteString("    pushq %r11\n") // sep
	cg.textSection.WriteString("    pushq %r12\n") // count
	cg.textSection.WriteString("    pushq %r9\n")  // sep_len
	cg.textSection.WriteString("    pushq %r8\n")  // total

	// Allocate result (total + 1)
	cg.textSection.WriteString("    movq %r8, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	lErr := cg.getLabel("join_err")
	lEnd := cg.getLabel("join_end")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))

	// Restore values
	cg.textSection.WriteString("    popq %r8\n")  // total (unused now)
	cg.textSection.WriteString("    popq %r9\n")  // sep_len
	cg.textSection.WriteString("    popq %r12\n") // count
	cg.textSection.WriteString("    popq %r11\n") // sep
	cg.textSection.WriteString("    popq %r10\n") // array

	// rax = dest buffer
	cg.textSection.WriteString("    movq %rax, %r15\n")    // save result ptr
	cg.textSection.WriteString("    movq %rax, %rdi\n")    // dest ptr
	cg.textSection.WriteString("    leaq 8(%r10), %r13\n") // first string slot
	cg.textSection.WriteString("    movq %r12, %r14\n")    // loop counter

	lCopyLoop := cg.getLabel("join_copy_loop")
	lCopyDone := cg.getLabel("join_copy_done")
	lCopyStr := cg.getLabel("join_copy_str")
	lCopyStrDone := cg.getLabel("join_copy_str_done")
	lCopySep := cg.getLabel("join_copy_sep")
	lCopySepDone := cg.getLabel("join_copy_sep_done")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyLoop))
	cg.textSection.WriteString("    testq %r14, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopyDone))

	// Copy current string
	cg.textSection.WriteString("    movq (%r13), %rsi\n") // get string ptr
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyStr))
	cg.textSection.WriteString("    movzbl (%rsi), %eax\n")
	cg.textSection.WriteString("    testb %al, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopyStrDone))
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCopyStr))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyStrDone))

	// If not last, add separator
	cg.textSection.WriteString("    cmpq $1, %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopySepDone))
	cg.textSection.WriteString("    movq %r11, %rsi\n") // sep ptr
	cg.textSection.WriteString("    movq %r9, %rcx\n")  // sep len
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopySep))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopySepDone))
	cg.textSection.WriteString("    movzbl (%rsi), %eax\n")
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCopySep))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopySepDone))

	cg.textSection.WriteString("    addq $8, %r13\n") // next slot
	cg.textSection.WriteString("    decq %r14\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCopyLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyDone))
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    movq %r15, %rax\n") // return result
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    addq $40, %rsp\n") // clean stack
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// generateStringReplace(str, old, new) -> new string with all occurrences replaced
func generateStringReplace(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// For simplicity, implement single-character replacement
	// Get source string
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", l))
			} else {
				// Compute length
				lLoop := cg.getLabel("replace_len_loop")
				lEnd := cg.getLabel("replace_len_end")
				cg.textSection.WriteString("    movq %r10, %rbx\n")
				cg.textSection.WriteString("    xorq %r8, %r8\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
				cg.textSection.WriteString("    testb %al, %al\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    incq %r8\n")
				cg.textSection.WriteString("    incq %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n")
			cg.textSection.WriteString("    xorq %r8, %r8\n")
		}
	default:
		cg.generateExpressionToReg(args[0], "r10")
		// Compute length
		lLoop := cg.getLabel("replace_len_loop")
		lEnd := cg.getLabel("replace_len_end")
		cg.textSection.WriteString("    movq %r10, %rbx\n")
		cg.textSection.WriteString("    xorq %r8, %r8\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
		cg.textSection.WriteString("    testb %al, %al\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    incq %r8\n")
		cg.textSection.WriteString("    incq %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	}

	// Get old char (first char of string)
	switch v := args[1].(type) {
	case *StringLiteral:
		if len(v.Value) > 0 {
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r11\n", byte(v.Value[0])))
		} else {
			cg.textSection.WriteString("    xorq %r11, %r11\n")
		}
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
			cg.textSection.WriteString("    movzbl (%rbx), %r11d\n")
		} else {
			cg.textSection.WriteString("    xorq %r11, %r11\n")
		}
	default:
		cg.generateExpressionToReg(args[1], "rbx")
		cg.textSection.WriteString("    movzbl (%rbx), %r11d\n")
	}

	// Get new char (first char of string)
	switch v := args[2].(type) {
	case *StringLiteral:
		if len(v.Value) > 0 {
			cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r12\n", byte(v.Value[0])))
		} else {
			cg.textSection.WriteString("    xorq %r12, %r12\n")
		}
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%rbx\n", varInfo.Offset))
			cg.textSection.WriteString("    movzbl (%rbx), %r12d\n")
		} else {
			cg.textSection.WriteString("    xorq %r12, %r12\n")
		}
	default:
		cg.generateExpressionToReg(args[2], "rbx")
		cg.textSection.WriteString("    movzbl (%rbx), %r12d\n")
	}

	// r10 = src, r8 = len, r11 = old char, r12 = new char
	// Save values before mmap
	cg.textSection.WriteString("    pushq %r10\n")
	cg.textSection.WriteString("    pushq %r8\n")
	cg.textSection.WriteString("    pushq %r11\n")
	cg.textSection.WriteString("    pushq %r12\n")

	// Allocate buffer (len + 1)
	cg.textSection.WriteString("    movq %r8, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    movq $3, %rdx\n")
	cg.textSection.WriteString("    movq $34, %r10\n")
	cg.textSection.WriteString("    movq $-1, %r8\n")
	cg.textSection.WriteString("    xorq %r9, %r9\n")
	cg.textSection.WriteString("    syscall\n")

	lErr := cg.getLabel("replace_err")
	lEnd := cg.getLabel("replace_end")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))

	// Restore values
	cg.textSection.WriteString("    popq %r12\n") // new char
	cg.textSection.WriteString("    popq %r11\n") // old char
	cg.textSection.WriteString("    popq %r8\n")  // len
	cg.textSection.WriteString("    popq %r10\n") // src

	cg.textSection.WriteString("    movq %rax, %r15\n") // save dest

	// Copy with replacement
	cg.textSection.WriteString("    movq %r10, %rsi\n") // src
	cg.textSection.WriteString("    movq %r15, %rdi\n") // dest
	cg.textSection.WriteString("    movq %r8, %rcx\n")  // len

	lCopyLoop := cg.getLabel("replace_copy_loop")
	lCopyDone := cg.getLabel("replace_copy_done")
	lNoReplace := cg.getLabel("replace_no_replace")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopyDone))
	cg.textSection.WriteString("    movzbl (%rsi), %eax\n")
	cg.textSection.WriteString("    cmpb %r11b, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    jne %s\n", lNoReplace))
	cg.textSection.WriteString("    movb %r12b, %al\n") // replace
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNoReplace))
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCopyLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyDone))
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    movq %r15, %rax\n") // return result
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    addq $32, %rsp\n") // clean stack
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// generateStringToLower(str) -> lowercased copy
func generateStringToLower(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Step 1: Get source string pointer into r10
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", l))
			} else {
				// Compute length at runtime
				lLoop := cg.getLabel("lower_len_loop")
				lEnd := cg.getLabel("lower_len_end")
				cg.textSection.WriteString("    movq %r10, %rbx\n")
				cg.textSection.WriteString("    xorq %r8, %r8\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    incq %r8\n")
				cg.textSection.WriteString("    incq %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n")
			cg.textSection.WriteString("    xorq %r8, %r8\n")
		}
	default:
		cg.generateExpressionToReg(args[0], "r10")
		// Compute length at runtime
		lLoop := cg.getLabel("lower_len_loop")
		lEnd := cg.getLabel("lower_len_end")
		cg.textSection.WriteString("    movq %r10, %rbx\n")
		cg.textSection.WriteString("    xorq %r8, %r8\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    incq %r8\n")
		cg.textSection.WriteString("    incq %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	}

	// Save source ptr and length before mmap
	cg.textSection.WriteString("    movq %r10, %r13\n") // save src ptr
	cg.textSection.WriteString("    movq %r8, %r14\n")  // save src len

	// Step 2: Allocate buffer with mmap (len + 1)
	cg.textSection.WriteString("    movq %r8, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")   // mmap syscall
	cg.textSection.WriteString("    xorq %rdi, %rdi\n") // addr = NULL
	cg.textSection.WriteString("    movq $3, %rdx\n")   // PROT_READ | PROT_WRITE
	cg.textSection.WriteString("    movq $34, %r10\n")  // MAP_PRIVATE | MAP_ANONYMOUS
	cg.textSection.WriteString("    movq $-1, %r8\n")   // fd = -1
	cg.textSection.WriteString("    xorq %r9, %r9\n")   // offset = 0
	cg.textSection.WriteString("    syscall\n")

	lErr := cg.getLabel("lower_mmap_err")
	lEnd := cg.getLabel("lower_end")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))

	// Save dest ptr
	cg.textSection.WriteString("    movq %rax, %r15\n") // save dest ptr

	// Step 3: Copy with lowercase conversion
	// rsi = src, rdi = dest, rcx = len
	cg.textSection.WriteString("    movq %r13, %rsi\n") // restore src
	cg.textSection.WriteString("    movq %r15, %rdi\n") // dest
	cg.textSection.WriteString("    movq %r14, %rcx\n") // restore len

	lCopyLoop := cg.getLabel("lower_copy_loop")
	lCopyDone := cg.getLabel("lower_copy_done")
	lNotUpper := cg.getLabel("lower_not_upper")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopyDone))
	cg.textSection.WriteString("    movzbl (%rsi), %eax\n")
	// Check if uppercase (A-Z: 65-90)
	cg.textSection.WriteString("    cmpb $65, %al\n") // 'A'
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lNotUpper))
	cg.textSection.WriteString("    cmpb $90, %al\n") // 'Z'
	cg.textSection.WriteString(fmt.Sprintf("    ja %s\n", lNotUpper))
	// Convert to lowercase: add 32
	cg.textSection.WriteString("    addb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNotUpper))
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCopyLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyDone))
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    movq %r15, %rax\n") // return dest ptr
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// generateStringToUpper(str) -> uppercased copy
func generateStringToUpper(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Step 1: Get source string pointer into r10
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", l))
			} else {
				// Compute length at runtime
				lLoop := cg.getLabel("upper_len_loop")
				lEnd := cg.getLabel("upper_len_end")
				cg.textSection.WriteString("    movq %r10, %rbx\n")
				cg.textSection.WriteString("    xorq %r8, %r8\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    incq %r8\n")
				cg.textSection.WriteString("    incq %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n")
			cg.textSection.WriteString("    xorq %r8, %r8\n")
		}
	default:
		cg.generateExpressionToReg(args[0], "r10")
		// Compute length at runtime
		lLoop := cg.getLabel("upper_len_loop")
		lEnd := cg.getLabel("upper_len_end")
		cg.textSection.WriteString("    movq %r10, %rbx\n")
		cg.textSection.WriteString("    xorq %r8, %r8\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    incq %r8\n")
		cg.textSection.WriteString("    incq %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	}

	// Save source ptr and length before mmap
	cg.textSection.WriteString("    movq %r10, %r13\n") // save src ptr
	cg.textSection.WriteString("    movq %r8, %r14\n")  // save src len

	// Step 2: Allocate buffer with mmap (len + 1)
	cg.textSection.WriteString("    movq %r8, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")   // mmap syscall
	cg.textSection.WriteString("    xorq %rdi, %rdi\n") // addr = NULL
	cg.textSection.WriteString("    movq $3, %rdx\n")   // PROT_READ | PROT_WRITE
	cg.textSection.WriteString("    movq $34, %r10\n")  // MAP_PRIVATE | MAP_ANONYMOUS
	cg.textSection.WriteString("    movq $-1, %r8\n")   // fd = -1
	cg.textSection.WriteString("    xorq %r9, %r9\n")   // offset = 0
	cg.textSection.WriteString("    syscall\n")

	lErr := cg.getLabel("upper_mmap_err")
	lEnd := cg.getLabel("upper_end")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))

	// Save dest ptr
	cg.textSection.WriteString("    movq %rax, %r15\n") // save dest ptr

	// Step 3: Copy with uppercase conversion
	// rsi = src, rdi = dest, rcx = len
	cg.textSection.WriteString("    movq %r13, %rsi\n") // restore src
	cg.textSection.WriteString("    movq %r15, %rdi\n") // dest
	cg.textSection.WriteString("    movq %r14, %rcx\n") // restore len

	lCopyLoop := cg.getLabel("upper_copy_loop")
	lCopyDone := cg.getLabel("upper_copy_done")
	lNotLower := cg.getLabel("upper_not_lower")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyLoop))
	cg.textSection.WriteString("    testq %rcx, %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lCopyDone))
	cg.textSection.WriteString("    movzbl (%rsi), %eax\n")
	// Check if lowercase (a-z: 97-122)
	cg.textSection.WriteString("    cmpb $97, %al\n") // 'a'
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lNotLower))
	cg.textSection.WriteString("    cmpb $122, %al\n") // 'z'
	cg.textSection.WriteString(fmt.Sprintf("    ja %s\n", lNotLower))
	// Convert to uppercase: subtract 32
	cg.textSection.WriteString("    subb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNotLower))
	cg.textSection.WriteString("    movb %al, (%rdi)\n")
	cg.textSection.WriteString("    incq %rsi\n")
	cg.textSection.WriteString("    incq %rdi\n")
	cg.textSection.WriteString("    decq %rcx\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lCopyLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lCopyDone))
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    movq %r15, %rax\n") // return dest ptr
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// generateStringTrim(str) -> trimmed string (whitespace removed from both ends)
func generateStringTrim(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}

	// Step 1: Get source string pointer into r10
	switch v := args[0].(type) {
	case *StringLiteral:
		lbl, _ := emitStringLiteral(cg, v.Value)
		cg.textSection.WriteString(fmt.Sprintf("    leaq %s(%%rip), %%r10\n", lbl))
		cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", len(v.Value)))
	case *Identifier:
		if varInfo, ok := cg.variables[v.Name]; ok {
			cg.textSection.WriteString(fmt.Sprintf("    movq -%d(%%rbp), %%r10\n", varInfo.Offset))
			if l, ok := cg.stringLengths[v.Name]; ok {
				cg.textSection.WriteString(fmt.Sprintf("    movq $%d, %%r8\n", l))
			} else {
				// Compute length at runtime
				lLoop := cg.getLabel("trim_len_loop")
				lEnd := cg.getLabel("trim_len_end")
				cg.textSection.WriteString("    movq %r10, %rbx\n")
				cg.textSection.WriteString("    xorq %r8, %r8\n")
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
				cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
				cg.textSection.WriteString("    testq %rax, %rax\n")
				cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
				cg.textSection.WriteString("    incq %r8\n")
				cg.textSection.WriteString("    incq %rbx\n")
				cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
				cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
			}
		} else {
			cg.textSection.WriteString("    xorq %r10, %r10\n")
			cg.textSection.WriteString("    xorq %r8, %r8\n")
		}
	default:
		cg.generateExpressionToReg(args[0], "r10")
		// Compute length at runtime
		lLoop := cg.getLabel("trim_len_loop")
		lEnd := cg.getLabel("trim_len_end")
		cg.textSection.WriteString("    movq %r10, %rbx\n")
		cg.textSection.WriteString("    xorq %r8, %r8\n")
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLoop))
		cg.textSection.WriteString("    movzbq (%rbx), %rax\n")
		cg.textSection.WriteString("    testq %rax, %rax\n")
		cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lEnd))
		cg.textSection.WriteString("    incq %r8\n")
		cg.textSection.WriteString("    incq %rbx\n")
		cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lLoop))
		cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
	}

	// r10 = src ptr, r8 = len
	// Step 2: Find start (skip leading whitespace)
	// r12 = start offset, r13 = end offset
	cg.textSection.WriteString("    movq %r10, %rbx\n") // current ptr
	cg.textSection.WriteString("    xorq %r12, %r12\n") // start offset = 0

	lTrimStart := cg.getLabel("trim_start_loop")
	lTrimStartDone := cg.getLabel("trim_start_done")
	lIsWhitespace := cg.getLabel("trim_is_ws")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrimStart))
	cg.textSection.WriteString("    cmpq %r8, %r12\n") // if start >= len, done
	cg.textSection.WriteString(fmt.Sprintf("    jge %s\n", lTrimStartDone))
	cg.textSection.WriteString("    movq %r10, %rbx\n")
	cg.textSection.WriteString("    addq %r12, %rbx\n")
	cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
	// Check whitespace: space(32), tab(9), newline(10), carriage return(13)
	cg.textSection.WriteString("    cmpb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWhitespace))
	cg.textSection.WriteString("    cmpb $9, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWhitespace))
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWhitespace))
	cg.textSection.WriteString("    cmpb $13, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWhitespace))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lTrimStartDone))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lIsWhitespace))
	cg.textSection.WriteString("    incq %r12\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lTrimStart))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrimStartDone))

	// Step 3: Find end (skip trailing whitespace)
	// r13 = end position (exclusive)
	cg.textSection.WriteString("    movq %r8, %r13\n") // end = len

	lTrimEnd := cg.getLabel("trim_end_loop")
	lTrimEndDone := cg.getLabel("trim_end_done")
	lIsWsEnd := cg.getLabel("trim_is_ws_end")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrimEnd))
	cg.textSection.WriteString("    cmpq %r12, %r13\n") // if end <= start, done
	cg.textSection.WriteString(fmt.Sprintf("    jle %s\n", lTrimEndDone))
	cg.textSection.WriteString("    movq %r10, %rbx\n")
	cg.textSection.WriteString("    addq %r13, %rbx\n")
	cg.textSection.WriteString("    decq %rbx\n") // check char at end-1
	cg.textSection.WriteString("    movzbl (%rbx), %eax\n")
	cg.textSection.WriteString("    cmpb $32, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWsEnd))
	cg.textSection.WriteString("    cmpb $9, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWsEnd))
	cg.textSection.WriteString("    cmpb $10, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWsEnd))
	cg.textSection.WriteString("    cmpb $13, %al\n")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lIsWsEnd))
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lTrimEndDone))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lIsWsEnd))
	cg.textSection.WriteString("    decq %r13\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lTrimEnd))
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lTrimEndDone))

	// r12 = start, r13 = end, new_len = r13 - r12
	cg.textSection.WriteString("    movq %r13, %r14\n")
	cg.textSection.WriteString("    subq %r12, %r14\n") // r14 = new length

	// Save values before mmap
	cg.textSection.WriteString("    pushq %r10\n") // save src ptr
	cg.textSection.WriteString("    pushq %r12\n") // save start offset
	cg.textSection.WriteString("    pushq %r14\n") // save new length

	// Step 4: Allocate buffer with mmap (new_len + 1)
	cg.textSection.WriteString("    movq %r14, %rsi\n")
	cg.textSection.WriteString("    addq $1, %rsi\n")
	cg.textSection.WriteString("    movq $9, %rax\n")   // mmap syscall
	cg.textSection.WriteString("    xorq %rdi, %rdi\n") // addr = NULL
	cg.textSection.WriteString("    movq $3, %rdx\n")   // PROT_READ | PROT_WRITE
	cg.textSection.WriteString("    movq $34, %r10\n")  // MAP_PRIVATE | MAP_ANONYMOUS
	cg.textSection.WriteString("    movq $-1, %r8\n")   // fd = -1
	cg.textSection.WriteString("    xorq %r9, %r9\n")   // offset = 0
	cg.textSection.WriteString("    syscall\n")

	lErr := cg.getLabel("trim_mmap_err")
	lEnd := cg.getLabel("trim_end")
	cg.textSection.WriteString("    testq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    js %s\n", lErr))

	// Restore values
	cg.textSection.WriteString("    popq %r14\n") // new length
	cg.textSection.WriteString("    popq %r12\n") // start offset
	cg.textSection.WriteString("    popq %r10\n") // src ptr

	// Save dest ptr
	cg.textSection.WriteString("    movq %rax, %r15\n")

	// Step 5: Copy trimmed portion
	cg.textSection.WriteString("    movq %r10, %rsi\n")
	cg.textSection.WriteString("    addq %r12, %rsi\n") // src + start
	cg.textSection.WriteString("    movq %r15, %rdi\n") // dest
	cg.textSection.WriteString("    movq %r14, %rcx\n") // length
	cg.textSection.WriteString("    rep movsb\n")
	cg.textSection.WriteString("    movb $0, (%rdi)\n") // NUL terminate
	cg.textSection.WriteString("    movq %r15, %rax\n") // return dest ptr
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lEnd))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lErr))
	cg.textSection.WriteString("    addq $24, %rsp\n") // clean up stack
	cg.textSection.WriteString("    xorq %rax, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lEnd))
}

// ============================================================================
// File I/O Functions (Phase 4)
// ============================================================================

// generateFileOpen(path_ptr, flags) -> fd or error
func generateFileOpen(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    movq $-1, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // path
	cg.generateExpressionToReg(args[1], "rsi") // flags
	// open(2) syscall: rax=2, rdi=path, rsi=flags, rdx=mode(0)
	cg.textSection.WriteString("    movq $2, %rax\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    syscall\n")
	// rax contains fd or error code
}

// generateFileClose(fd) -> status
func generateFileClose(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	// close(2) syscall: rax=3, rdi=fd
	cg.textSection.WriteString("    movq $3, %rax\n")
	cg.textSection.WriteString("    syscall\n")
}

// generateFileRead(fd, buf_ptr, size) -> bytes_read
func generateFileRead(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	cg.generateExpressionToReg(args[1], "rsi") // buf_ptr
	cg.generateExpressionToReg(args[2], "rdx") // size
	// read(2) syscall: rax=0, rdi=fd, rsi=buf, rdx=count
	cg.textSection.WriteString("    movq $0, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	// rax contains bytes_read or error
}

// generateFileWrite(fd, buf_ptr, size) -> bytes_written
func generateFileWrite(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	cg.generateExpressionToReg(args[1], "rsi") // buf_ptr
	cg.generateExpressionToReg(args[2], "rdx") // size
	// write(2) syscall: rax=1, rdi=fd, rsi=buf, rdx=count
	cg.textSection.WriteString("    movq $1, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	// rax contains bytes_written or error
}

// generateFileSeek(fd, offset, whence) -> new_pos
func generateFileSeek(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 3 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // fd
	cg.generateExpressionToReg(args[1], "rsi") // offset
	cg.generateExpressionToReg(args[2], "rdx") // whence (0=SEEK_SET, 1=SEEK_CUR, 2=SEEK_END)
	// lseek(2) syscall: rax=8, rdi=fd, rsi=offset, rdx=whence
	cg.textSection.WriteString("    movq $8, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	// rax contains new position or error
}

// generateFileStat(path_ptr, stat_buf) -> status
func generateFileStat(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		cg.textSection.WriteString("    movq $-1, %rax\n")
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // path
	cg.generateExpressionToReg(args[1], "rsi") // stat buffer
	// stat(2) syscall: rax=4, rdi=path, rsi=statbuf
	cg.textSection.WriteString("    movq $4, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	// rax contains status (0 on success, < 0 on error)
}

// generateFileExists(path_ptr) -> 0 or 1
func generateFileExists(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		cg.textSection.WriteString("    xorq %rax, %rax\n")
		return
	}
	// Allocate stat buffer on stack and call stat
	cg.generateExpressionToReg(args[0], "rdi")          // path
	cg.textSection.WriteString("    subq $144, %rsp\n") // stat buffer (144 bytes)
	cg.textSection.WriteString("    movq %rsp, %rsi\n") // stat buffer ptr
	cg.textSection.WriteString("    movq $4, %rax\n")   // stat syscall
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $144, %rsp\n")
	cg.textSection.WriteString("    movq $0, %rax\n") // default: doesn't exist
	cg.textSection.WriteString("    cmp $0, %rax\n")  // if stat returned 0, file exists
	lExistsLabel := cg.getLabel("file_exists")
	cg.textSection.WriteString(fmt.Sprintf("    je %s\n", lExistsLabel))
	cg.textSection.WriteString("    movq $1, %rax\n") // file exists
	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lExistsLabel))
}

// ============================================================================
// Time Functions (Phase 4)
// ============================================================================

// generateTimeNow() -> unix timestamp
func generateTimeNow(cg *CodeGenerator, args []ASTNode) {
	// time(2) syscall: rax=201, rdi=NULL, rsi=NULL
	cg.textSection.WriteString("    movq $201, %rax\n")
	cg.textSection.WriteString("    xorq %rdi, %rdi\n")
	cg.textSection.WriteString("    xorq %rsi, %rsi\n")
	cg.textSection.WriteString("    syscall\n")
	// rax contains current unix timestamp
}

// generateTimeSleep(seconds) -> status
func generateTimeSleep(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 1 {
		return
	}
	cg.generateExpressionToReg(args[0], "rdi") // seconds
	// nanosleep(2) syscall: rax=35, rdi=timespec ptr, rsi=remaining
	// For simplicity, allocate timespec on stack: [seconds(8), nanoseconds(8)]
	cg.textSection.WriteString("    subq $16, %rsp\n")     // timespec structure
	cg.textSection.WriteString("    movq %rdi, 0(%rsp)\n") // tv_sec
	cg.textSection.WriteString("    movq $0, 8(%rsp)\n")   // tv_nsec = 0
	cg.textSection.WriteString("    movq %rsp, %rdi\n")
	cg.textSection.WriteString("    xorq %rsi, %rsi\n")
	cg.textSection.WriteString("    movq $35, %rax\n")
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    addq $16, %rsp\n")
}

// generateTimeMillis() -> current milliseconds since epoch
func generateTimeMillis(cg *CodeGenerator, args []ASTNode) {
	// clock_gettime(2) syscall: rax=228, rdi=CLOCK_REALTIME(0), rsi=timespec ptr
	// Returns: timespec { tv_sec, tv_nsec }
	// milliseconds = tv_sec * 1000 + tv_nsec / 1000000
	cg.textSection.WriteString("    subq $16, %rsp\n")  // timespec on stack
	cg.textSection.WriteString("    movq $228, %rax\n") // clock_gettime syscall
	cg.textSection.WriteString("    xorq %rdi, %rdi\n") // CLOCK_REALTIME = 0
	cg.textSection.WriteString("    movq %rsp, %rsi\n") // timespec ptr
	cg.textSection.WriteString("    syscall\n")
	cg.textSection.WriteString("    movq 0(%rsp), %rax\n") // tv_sec
	cg.textSection.WriteString("    imulq $1000, %rax\n")  // tv_sec * 1000
	cg.textSection.WriteString("    movq 8(%rsp), %rcx\n") // tv_nsec
	cg.textSection.WriteString("    movq $1000000, %rdx\n")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    movq $1000000, %rcx\n")
	cg.textSection.WriteString("    divq %rcx\n") // tv_nsec / 1000000
	cg.textSection.WriteString("    popq %rcx\n")
	cg.textSection.WriteString("    addq %rcx, %rax\n") // total milliseconds
	cg.textSection.WriteString("    addq $16, %rsp\n")
}

// generateTimeNanos() -> current nanoseconds (monotonic for timing)
func generateTimeNanos(cg *CodeGenerator, args []ASTNode) {
	// clock_gettime(2) with CLOCK_MONOTONIC(1) for timing purposes
	// Returns raw nanoseconds portion (useful for elapsed time measurement)
	cg.textSection.WriteString("    subq $16, %rsp\n")  // timespec on stack
	cg.textSection.WriteString("    movq $228, %rax\n") // clock_gettime syscall
	cg.textSection.WriteString("    movq $1, %rdi\n")   // CLOCK_MONOTONIC = 1
	cg.textSection.WriteString("    movq %rsp, %rsi\n") // timespec ptr
	cg.textSection.WriteString("    syscall\n")
	// Compute: tv_sec * 1000000000 + tv_nsec
	cg.textSection.WriteString("    movq 0(%rsp), %rax\n") // tv_sec
	cg.textSection.WriteString("    movq $1000000000, %rcx\n")
	cg.textSection.WriteString("    imulq %rcx, %rax\n")   // tv_sec * 1e9
	cg.textSection.WriteString("    addq 8(%rsp), %rax\n") // + tv_nsec
	cg.textSection.WriteString("    addq $16, %rsp\n")
}

// generateTimeClock() -> CPU clock ticks (via RDTSC instruction)
func generateTimeClock(cg *CodeGenerator, args []ASTNode) {
	// Use RDTSC (Read Time-Stamp Counter) for CPU cycles
	cg.textSection.WriteString("    rdtsc\n")          // edx:eax = timestamp counter
	cg.textSection.WriteString("    shlq $32, %rdx\n") // shift high 32 bits
	cg.textSection.WriteString("    orq %rdx, %rax\n") // combine into 64-bit value
}

// generateTimeGMTime(timestamp, tm_buf) -> void
// Converts Unix timestamp to broken-down time (UTC)
// tm_buf layout (all i64): [sec][min][hour][mday][mon][year][wday][yday][isdst]
func generateTimeGMTime(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}

	// Get timestamp into r10
	cg.generateExpressionToReg(args[0], "r10")
	// Get buffer pointer into r11
	cg.generateExpressionToReg(args[1], "r11")

	// Constants
	cg.textSection.WriteString("    movq %r10, %rax\n") // timestamp

	// Calculate days since epoch and time of day
	// seconds_in_day = 86400
	cg.textSection.WriteString("    movq $86400, %rcx\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %rcx\n")       // rax = days, rdx = seconds remaining
	cg.textSection.WriteString("    movq %rax, %r12\n") // r12 = days since epoch
	cg.textSection.WriteString("    movq %rdx, %r13\n") // r13 = seconds in day

	// Calculate hour, minute, second from r13
	cg.textSection.WriteString("    movq %r13, %rax\n")
	cg.textSection.WriteString("    movq $3600, %rcx\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %rcx\n")       // rax = hours, rdx = remaining
	cg.textSection.WriteString("    movq %rax, %r14\n") // r14 = hours
	cg.textSection.WriteString("    movq %rdx, %rax\n")
	cg.textSection.WriteString("    movq $60, %rcx\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %rcx\n")       // rax = minutes, rdx = seconds
	cg.textSection.WriteString("    movq %rax, %r15\n") // r15 = minutes
	cg.textSection.WriteString("    movq %rdx, %r8\n")  // r8 = seconds

	// Store sec, min, hour
	cg.textSection.WriteString("    movq %r8, (%r11)\n")    // tm_sec
	cg.textSection.WriteString("    movq %r15, 8(%r11)\n")  // tm_min
	cg.textSection.WriteString("    movq %r14, 16(%r11)\n") // tm_hour

	// Calculate day of week: (days + 4) % 7 (Jan 1, 1970 was Thursday = 4)
	cg.textSection.WriteString("    movq %r12, %rax\n")
	cg.textSection.WriteString("    addq $4, %rax\n")
	cg.textSection.WriteString("    movq $7, %rcx\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %rcx\n")
	cg.textSection.WriteString("    movq %rdx, 48(%r11)\n") // tm_wday

	// Calculate year, month, day using simplified algorithm
	// This is a simplified version - works for dates after 1970
	cg.textSection.WriteString("    movq %r12, %rax\n") // days since epoch
	cg.textSection.WriteString("    movq $1970, %r9\n") // year counter

	// Year loop - subtract days per year
	lYearLoop := cg.getLabel("gmtime_year_loop")
	lYearDone := cg.getLabel("gmtime_year_done")
	lLeapYear := cg.getLabel("gmtime_leap_year")
	lNotLeap := cg.getLabel("gmtime_not_leap")

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lYearLoop))
	// Check if leap year: divisible by 4 and (not by 100 or by 400)
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    andq $3, %rcx\n") // year % 4
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lNotLeap))
	// Divisible by 4, check if not by 100 or by 400
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.textSection.WriteString("    pushq %rdx\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    movq $100, %rcx\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %rcx\n")
	cg.textSection.WriteString("    testq %rdx, %rdx\n") // year % 100
	cg.textSection.WriteString("    popq %rdx\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lLeapYear)) // not divisible by 100 = leap
	// Divisible by 100, check if by 400
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    pushq %rax\n")
	cg.textSection.WriteString("    pushq %rdx\n")
	cg.textSection.WriteString("    movq %rcx, %rax\n")
	cg.textSection.WriteString("    movq $400, %rcx\n")
	cg.textSection.WriteString("    xorq %rdx, %rdx\n")
	cg.textSection.WriteString("    divq %rcx\n")
	cg.textSection.WriteString("    testq %rdx, %rdx\n")
	cg.textSection.WriteString("    popq %rdx\n")
	cg.textSection.WriteString("    popq %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jnz %s\n", lNotLeap)) // not divisible by 400 = not leap

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lLeapYear))
	cg.textSection.WriteString("    cmpq $366, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lYearDone))
	cg.textSection.WriteString("    subq $366, %rax\n")
	cg.textSection.WriteString("    incq %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lYearLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lNotLeap))
	cg.textSection.WriteString("    cmpq $365, %rax\n")
	cg.textSection.WriteString(fmt.Sprintf("    jb %s\n", lYearDone))
	cg.textSection.WriteString("    subq $365, %rax\n")
	cg.textSection.WriteString("    incq %r9\n")
	cg.textSection.WriteString(fmt.Sprintf("    jmp %s\n", lYearLoop))

	cg.textSection.WriteString(fmt.Sprintf("%s:\n", lYearDone))
	// rax = day of year (0-based), r9 = year
	cg.textSection.WriteString("    movq %rax, 56(%r11)\n") // tm_yday

	// Calculate year - 1900 for tm_year
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    subq $1900, %rcx\n")
	cg.textSection.WriteString("    movq %rcx, 40(%r11)\n") // tm_year

	// Calculate month and day from day of year
	// Use month lengths: 31,28/29,31,30,31,30,31,31,30,31,30,31
	cg.textSection.WriteString("    movq %rax, %r12\n") // day of year
	cg.textSection.WriteString("    xorq %r13, %r13\n") // month = 0

	// Simplified: subtract month lengths
	// January
	cg.textSection.WriteString("    cmpq $31, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $31, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// February (check leap year)
	cg.textSection.WriteString("    movq %r9, %rcx\n")
	cg.textSection.WriteString("    andq $3, %rcx\n")
	cg.textSection.WriteString("    jnz 2f\n")
	cg.textSection.WriteString("    cmpq $29, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $29, %r12\n")
	cg.textSection.WriteString("    jmp 3f\n")
	cg.textSection.WriteString("2:\n")
	cg.textSection.WriteString("    cmpq $28, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $28, %r12\n")
	cg.textSection.WriteString("3:\n")
	cg.textSection.WriteString("    incq %r13\n")
	// March
	cg.textSection.WriteString("    cmpq $31, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $31, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// April
	cg.textSection.WriteString("    cmpq $30, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $30, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// May
	cg.textSection.WriteString("    cmpq $31, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $31, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// June
	cg.textSection.WriteString("    cmpq $30, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $30, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// July
	cg.textSection.WriteString("    cmpq $31, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $31, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// August
	cg.textSection.WriteString("    cmpq $31, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $31, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// September
	cg.textSection.WriteString("    cmpq $30, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $30, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// October
	cg.textSection.WriteString("    cmpq $31, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $31, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// November
	cg.textSection.WriteString("    cmpq $30, %r12\n")
	cg.textSection.WriteString("    jb 1f\n")
	cg.textSection.WriteString("    subq $30, %r12\n")
	cg.textSection.WriteString("    incq %r13\n")
	// December (remaining days)
	cg.textSection.WriteString("1:\n")

	// Store month (0-based) and day (1-based)
	cg.textSection.WriteString("    movq %r13, 32(%r11)\n") // tm_mon (0-11)
	cg.textSection.WriteString("    incq %r12\n")           // day is 1-based
	cg.textSection.WriteString("    movq %r12, 24(%r11)\n") // tm_mday

	// isdst = 0 for UTC
	cg.textSection.WriteString("    movq $0, 64(%r11)\n") // tm_isdst
}

// generateTimeLocalTime(timestamp, tm_buf) -> void
// For simplicity, this is the same as gmtime (no timezone support)
func generateTimeLocalTime(cg *CodeGenerator, args []ASTNode) {
	if len(args) != 2 {
		return
	}
	// For now, localtime is the same as gmtime (UTC)
	// Full implementation would need timezone support
	generateTimeGMTime(cg, args)
}
