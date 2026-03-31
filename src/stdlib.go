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
	"random":      createRandomModule(),
	"regex":       createRegexModule(),
	"json":        createJSONModule(),
	"sdl3":        createSDL3Module(),
	"sdl_mixer":   createSDLMixerModule(),
	"os":          createOSModule(),
	"func":        createFuncModule(),
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
			},
			"println": {
				Name:    "println",
				Module:  "io",
				NumArgs: -1,
			},
			"printf": {
				Name:    "printf",
				Module:  "io",
				NumArgs: -1,
			},
			"fprintf": {
				Name:    "fprintf",
				Module:  "io",
				NumArgs: -1,
			},
			"sprint": {
				Name:    "sprint",
				Module:  "io",
				NumArgs: -1,
			},
			"sprintf": {
				Name:    "sprintf",
				Module:  "io",
				NumArgs: -1,
			},
			"sprintln": {
				Name:    "sprintln",
				Module:  "io",
				NumArgs: -1,
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
			},
			"free": {
				Name:    "free",
				Module:  "mem",
				NumArgs: 1,
			},
			"sizeof": {
				Name:    "sizeof",
				Module:  "mem",
				NumArgs: 1,
			},
			"memcpy": {
				Name:    "memcpy",
				Module:  "mem",
				NumArgs: 3,
			},
			"memset": {
				Name:    "memset",
				Module:  "mem",
				NumArgs: 3,
			},
			"mmap": {
				Name:    "mmap",
				Module:  "mem",
				NumArgs: 1,
			},
			"munmap": {
				Name:    "munmap",
				Module:  "mem",
				NumArgs: 2,
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
			},
			"min": {
				Name:    "min",
				Module:  "math",
				NumArgs: 2,
			},
			"max": {
				Name:    "max",
				Module:  "math",
				NumArgs: 2,
			},
			"sqrt": {
				Name:    "sqrt",
				Module:  "math",
				NumArgs: 1,
			},
			"pow": {
				Name:    "pow",
				Module:  "math",
				NumArgs: 2,
			},
			"floor": {
				Name:    "floor",
				Module:  "math",
				NumArgs: 1,
			},
			"ceil": {
				Name:    "ceil",
				Module:  "math",
				NumArgs: 1,
			},
			"round": {
				Name:    "round",
				Module:  "math",
				NumArgs: 1,
			},
			"gcd": {
				Name:    "gcd",
				Module:  "math",
				NumArgs: 2,
			},
			"lcm": {
				Name:    "lcm",
				Module:  "math",
				NumArgs: 2,
			},
			"sin": {
				Name:    "sin",
				Module:  "math",
				NumArgs: 1,
			},
			"cos": {
				Name:    "cos",
				Module:  "math",
				NumArgs: 1,
			},
			"tan": {
				Name:    "tan",
				Module:  "math",
				NumArgs: 1,
			},
			"atan2": {
				Name:    "atan2",
				Module:  "math",
				NumArgs: 2,
			},
			"asin": {
				Name:    "asin",
				Module:  "math",
				NumArgs: 1,
			},
			"acos": {
				Name:    "acos",
				Module:  "math",
				NumArgs: 1,
			},
			"fmod": {
				Name:    "fmod",
				Module:  "math",
				NumArgs: 2,
			},
			"fabs": {
				Name:    "fabs",
				Module:  "math",
				NumArgs: 1,
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
			},
			"concat": {
				Name:    "concat",
				Module:  "str",
				NumArgs: -1,
			},
			"compare": {
				Name:    "compare",
				Module:  "str",
				NumArgs: 2,
			},
			"copy": {
				Name:    "copy",
				Module:  "str",
				NumArgs: 1,
			},
			"indexOf": {
				Name:    "indexOf",
				Module:  "str",
				NumArgs: 2,
			},
			"contains": {
				Name:    "contains",
				Module:  "str",
				NumArgs: 2,
			},
			"startsWith": {
				Name:    "startsWith",
				Module:  "str",
				NumArgs: 2,
			},
			"endsWith": {
				Name:    "endsWith",
				Module:  "str",
				NumArgs: 2,
			},
			"substring": {
				Name:    "substring",
				Module:  "str",
				NumArgs: 3,
			},
			"split": {
				Name:    "split",
				Module:  "str",
				NumArgs: 2,
			},
			"join": {
				Name:    "join",
				Module:  "str",
				NumArgs: 2,
			},
			"replace": {
				Name:    "replace",
				Module:  "str",
				NumArgs: 3,
			},
			"toLower": {
				Name:    "toLower",
				Module:  "str",
				NumArgs: 1,
			},
			"toUpper": {
				Name:    "toUpper",
				Module:  "str",
				NumArgs: 1,
			},
			"trim": {
				Name:    "trim",
				Module:  "str",
				NumArgs: 1,
			},
		},
		Types: map[string]TokenType{},
	}
}

func createNetModule() *StdlibModule {
	return &StdlibModule{
		Name: "net",
		Functions: map[string]*StdlibFunction{
			"socket":       {Name: "socket", Module: "net", NumArgs: 3},
			"connect_ipv4": {Name: "connect_ipv4", Module: "net", NumArgs: 3},
			"send":         {Name: "send", Module: "net", NumArgs: 3},
			"recv":         {Name: "recv", Module: "net", NumArgs: 3},
			"close":        {Name: "close", Module: "net", NumArgs: 1},
			// TCP server support
			"bind_ipv4":  {Name: "bind_ipv4", Module: "net", NumArgs: 3},
			"listen":     {Name: "listen", Module: "net", NumArgs: 2},
			"accept":     {Name: "accept", Module: "net", NumArgs: 1},
			"setsockopt": {Name: "setsockopt", Module: "net", NumArgs: 4},
			// UDP support
			"sendto_ipv4": {Name: "sendto_ipv4", Module: "net", NumArgs: 5},
			"recvfrom":    {Name: "recvfrom", Module: "net", NumArgs: 3},
			// IPv6 support
			"connect_ipv6": {Name: "connect_ipv6", Module: "net", NumArgs: 3},
			"bind_ipv6":    {Name: "bind_ipv6", Module: "net", NumArgs: 3},
			"sendto_ipv6":  {Name: "sendto_ipv6", Module: "net", NumArgs: 5},
			// DNS resolution
			"resolve":      {Name: "resolve", Module: "net", NumArgs: 2},
			"resolve_ipv6": {Name: "resolve_ipv6", Module: "net", NumArgs: 2},
		},
		Types: map[string]TokenType{},
	}
}

// createHTTPModule creates a minimal HTTP client module built on net helpers
func createHTTPModule() *StdlibModule {
	return &StdlibModule{
		Name: "http",
		Functions: map[string]*StdlibFunction{
			"get":  {Name: "get", Module: "http", NumArgs: 7},
			"post": {Name: "post", Module: "http", NumArgs: 9},
			// Response parsing
			"parse_status":  {Name: "parse_status", Module: "http", NumArgs: 2},
			"get_header":    {Name: "get_header", Module: "http", NumArgs: 4},
			"get_body":      {Name: "get_body", Module: "http", NumArgs: 2},
			"parse_headers": {Name: "parse_headers", Module: "http", NumArgs: 3},
			// Connection pooling
			"pool_new":   {Name: "pool_new", Module: "http", NumArgs: 1},     // pool_new(max_conns) -> pool_ptr
			"pool_get":   {Name: "pool_get", Module: "http", NumArgs: 3},     // pool_get(pool, host_ptr, port) -> fd or -1
			"pool_put":   {Name: "pool_put", Module: "http", NumArgs: 4},     // pool_put(pool, fd, host_ptr, port) -> 0/1
			"pool_close": {Name: "pool_close", Module: "http", NumArgs: 1}, // pool_close(pool) -> void
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
			"array_int_new":      {Name: "array_int_new", Module: "collections", NumArgs: 1},
			"array_int_push":     {Name: "array_int_push", Module: "collections", NumArgs: 2},
			"array_int_pop":      {Name: "array_int_pop", Module: "collections", NumArgs: 1},
			"array_int_len":      {Name: "array_int_len", Module: "collections", NumArgs: 1},
			"array_int_capacity": {Name: "array_int_capacity", Module: "collections", NumArgs: 1},
			"array_int_resize":   {Name: "array_int_resize", Module: "collections", NumArgs: 2},
			"array_int_reserve":  {Name: "array_int_reserve", Module: "collections", NumArgs: 2},
			"array_int_shrink":   {Name: "array_int_shrink", Module: "collections", NumArgs: 1},
			"array_int_get":      {Name: "array_int_get", Module: "collections", NumArgs: 2},
			"array_int_set":      {Name: "array_int_set", Module: "collections", NumArgs: 3},
			"array_int_free":     {Name: "array_int_free", Module: "collections", NumArgs: 1},

			// Stack
			"stack_int_new":  {Name: "stack_int_new", Module: "collections", NumArgs: 1},
			"stack_int_push": {Name: "stack_int_push", Module: "collections", NumArgs: 2},
			"stack_int_pop":  {Name: "stack_int_pop", Module: "collections", NumArgs: 1},
			"stack_int_len":  {Name: "stack_int_len", Module: "collections", NumArgs: 1},

			// Queue / Deque
			"queue_int_new":     {Name: "queue_int_new", Module: "collections", NumArgs: 1},
			"queue_int_enqueue": {Name: "queue_int_enqueue", Module: "collections", NumArgs: 2},
			"queue_int_dequeue": {Name: "queue_int_dequeue", Module: "collections", NumArgs: 1},
			"queue_int_len":     {Name: "queue_int_len", Module: "collections", NumArgs: 1},

			"deque_int_new":        {Name: "deque_int_new", Module: "collections", NumArgs: 1},
			"deque_int_push_front": {Name: "deque_int_push_front", Module: "collections", NumArgs: 2},
			"deque_int_push_back":  {Name: "deque_int_push_back", Module: "collections", NumArgs: 2},
			"deque_int_pop_front":  {Name: "deque_int_pop_front", Module: "collections", NumArgs: 1},
			"deque_int_pop_back":   {Name: "deque_int_pop_back", Module: "collections", NumArgs: 1},
			"deque_int_len":        {Name: "deque_int_len", Module: "collections", NumArgs: 1},

			// Heap (min-heap)
			"heap_int_new":  {Name: "heap_int_new", Module: "collections", NumArgs: 1},
			"heap_int_push": {Name: "heap_int_push", Module: "collections", NumArgs: 2},
			"heap_int_pop":  {Name: "heap_int_pop", Module: "collections", NumArgs: 1},
			"heap_int_peek": {Name: "heap_int_peek", Module: "collections", NumArgs: 1},
			"heap_int_len":  {Name: "heap_int_len", Module: "collections", NumArgs: 1},

			// Hash map & set (int keys)
			"hashmap_int_new":    {Name: "hashmap_int_new", Module: "collections", NumArgs: 1},
			"hashmap_int_put":    {Name: "hashmap_int_put", Module: "collections", NumArgs: 3},
			"hashmap_int_get":    {Name: "hashmap_int_get", Module: "collections", NumArgs: 2},
			"hashmap_int_remove": {Name: "hashmap_int_remove", Module: "collections", NumArgs: 2},
			"hashmap_int_len":    {Name: "hashmap_int_len", Module: "collections", NumArgs: 1},
			"hashmap_int_clear":  {Name: "hashmap_int_clear", Module: "collections", NumArgs: 1},
			"hashmap_int_free":   {Name: "hashmap_int_free", Module: "collections", NumArgs: 1},

			"hashset_int_new":      {Name: "hashset_int_new", Module: "collections", NumArgs: 1},
			"hashset_int_add":      {Name: "hashset_int_add", Module: "collections", NumArgs: 2},
			"hashset_int_contains": {Name: "hashset_int_contains", Module: "collections", NumArgs: 2},
			"hashset_int_remove":   {Name: "hashset_int_remove", Module: "collections", NumArgs: 2},
			"hashset_int_len":      {Name: "hashset_int_len", Module: "collections", NumArgs: 1},
			"hashset_int_clear":    {Name: "hashset_int_clear", Module: "collections", NumArgs: 1},
			"hashset_int_free":     {Name: "hashset_int_free", Module: "collections", NumArgs: 1},

			// Hash map & set (string keys)
			"hashmap_str_new":      {Name: "hashmap_str_new", Module: "collections", NumArgs: 1},
			"hashmap_str_put":      {Name: "hashmap_str_put", Module: "collections", NumArgs: 3},
			"hashmap_str_get":      {Name: "hashmap_str_get", Module: "collections", NumArgs: 2},
			"hashmap_str_contains": {Name: "hashmap_str_contains", Module: "collections", NumArgs: 2},
			"hashmap_str_remove":   {Name: "hashmap_str_remove", Module: "collections", NumArgs: 2},
			"hashmap_str_len":      {Name: "hashmap_str_len", Module: "collections", NumArgs: 1},
			"hashmap_str_clear":    {Name: "hashmap_str_clear", Module: "collections", NumArgs: 1},
			"hashmap_str_free":     {Name: "hashmap_str_free", Module: "collections", NumArgs: 1},

			"hashset_str_new":      {Name: "hashset_str_new", Module: "collections", NumArgs: 1},
			"hashset_str_add":      {Name: "hashset_str_add", Module: "collections", NumArgs: 2},
			"hashset_str_contains": {Name: "hashset_str_contains", Module: "collections", NumArgs: 2},
			"hashset_str_remove":   {Name: "hashset_str_remove", Module: "collections", NumArgs: 2},
			"hashset_str_len":      {Name: "hashset_str_len", Module: "collections", NumArgs: 1},
			"hashset_str_clear":    {Name: "hashset_str_clear", Module: "collections", NumArgs: 1},
			"hashset_str_free":     {Name: "hashset_str_free", Module: "collections", NumArgs: 1},

			// Sorted set (BST-based, maintains sorted order)
			"sortedset_int_new":      {Name: "sortedset_int_new", Module: "collections", NumArgs: 0},
			"sortedset_int_add":      {Name: "sortedset_int_add", Module: "collections", NumArgs: 2},
			"sortedset_int_contains": {Name: "sortedset_int_contains", Module: "collections", NumArgs: 2},
			"sortedset_int_remove":   {Name: "sortedset_int_remove", Module: "collections", NumArgs: 2},
			"sortedset_int_min":      {Name: "sortedset_int_min", Module: "collections", NumArgs: 1},
			"sortedset_int_max":      {Name: "sortedset_int_max", Module: "collections", NumArgs: 1},
			"sortedset_int_len":      {Name: "sortedset_int_len", Module: "collections", NumArgs: 1},
			"sortedset_int_free":     {Name: "sortedset_int_free", Module: "collections", NumArgs: 1},

			// Sorted map (BST-based, maintains sorted order by key)
			"sortedmap_int_new":      {Name: "sortedmap_int_new", Module: "collections", NumArgs: 0},
			"sortedmap_int_put":      {Name: "sortedmap_int_put", Module: "collections", NumArgs: 3},
			"sortedmap_int_get":      {Name: "sortedmap_int_get", Module: "collections", NumArgs: 2},
			"sortedmap_int_contains": {Name: "sortedmap_int_contains", Module: "collections", NumArgs: 2},
			"sortedmap_int_remove":   {Name: "sortedmap_int_remove", Module: "collections", NumArgs: 2},
			"sortedmap_int_min_key":  {Name: "sortedmap_int_min_key", Module: "collections", NumArgs: 1},
			"sortedmap_int_max_key":  {Name: "sortedmap_int_max_key", Module: "collections", NumArgs: 1},
			"sortedmap_int_len":      {Name: "sortedmap_int_len", Module: "collections", NumArgs: 1},
			"sortedmap_int_free":     {Name: "sortedmap_int_free", Module: "collections", NumArgs: 1},

			// Array helper
			"binary_search_int": {Name: "binary_search_int", Module: "collections", NumArgs: 3},
		},
		Types: map[string]TokenType{},
	}
}

// createNumModule creates the numeric conversions stdlib module
func createNumModule() *StdlibModule {
	return &StdlibModule{
		Name: "num",
		Functions: map[string]*StdlibFunction{
			"toInt8":   {Name: "toInt8", Module: "num", NumArgs: 1},
			"toUint8":  {Name: "toUint8", Module: "num", NumArgs: 1},
			"toInt16":  {Name: "toInt16", Module: "num", NumArgs: 1},
			"toUint16": {Name: "toUint16", Module: "num", NumArgs: 1},
			"toInt32":  {Name: "toInt32", Module: "num", NumArgs: 1},
			"toUint32": {Name: "toUint32", Module: "num", NumArgs: 1},
			"toInt64":  {Name: "toInt64", Module: "num", NumArgs: 1},
			"toUint64": {Name: "toUint64", Module: "num", NumArgs: 1},
			"toBool":   {Name: "toBool", Module: "num", NumArgs: 1},
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
			"crc32":  {Name: "crc32", Module: "hash", NumArgs: 2},    // crc32(data_ptr, len) -> uint32
			"fnv1a":  {Name: "fnv1a", Module: "hash", NumArgs: 2},    // fnv1a(data_ptr, len) -> uint64
			"djb2":   {Name: "djb2", Module: "hash", NumArgs: 1},      // djb2(string_ptr) -> uint64
			"murmur": {Name: "murmur", Module: "hash", NumArgs: 3}, // murmur(data_ptr, len, seed) -> uint32

			// Cryptographic hashes
			"sha256": {Name: "sha256", Module: "hash", NumArgs: 3}, // sha256(data_ptr, len, out_buf) -> void
			"md5":    {Name: "md5", Module: "hash", NumArgs: 3},       // md5(data_ptr, len, out_buf) -> void
		},
		Types: map[string]TokenType{},
	}
}

// createFileModule creates a file I/O stdlib module (POSIX file operations)
func createFileModule() *StdlibModule {
	return &StdlibModule{
		Name: "file",
		Functions: map[string]*StdlibFunction{
			"open":   {Name: "open", Module: "file", NumArgs: 2},     // open(path_ptr, flags) -> fd
			"close":  {Name: "close", Module: "file", NumArgs: 1},   // close(fd) -> status
			"read":   {Name: "read", Module: "file", NumArgs: 3},     // read(fd, buf_ptr, size) -> bytes_read
			"write":  {Name: "write", Module: "file", NumArgs: 3},   // write(fd, buf_ptr, size) -> bytes_written
			"seek":   {Name: "seek", Module: "file", NumArgs: 3},     // seek(fd, offset, whence) -> new_pos
			"stat":   {Name: "stat", Module: "file", NumArgs: 2},     // stat(path_ptr, stat_buf) -> status
			"exists": {Name: "exists", Module: "file", NumArgs: 1}, // exists(path_ptr) -> 0/1
		},
		Types: map[string]TokenType{},
	}
}

// createOSModule creates the OS / process-execution stdlib module.
// All functions call POSIX libc helpers (system, popen, pclose, fread,
// getenv, setenv) via the PLT so no extra link flags are required.
func createOSModule() *StdlibModule {
	return &StdlibModule{
		Name: "os",
		Functions: map[string]*StdlibFunction{
			// exec(cmd) -> exit_code  (wraps system(3))
			"exec": {Name: "exec", Module: "os", NumArgs: 1},
			// popen(cmd, mode) -> FILE*  (wraps popen(3); mode 0="r", 1="w")
			"popen": {Name: "popen", Module: "os", NumArgs: 2},
			// pread(fp, buf, size) -> bytes_read  (wraps fread(3) on a popen pipe)
			"pread": {Name: "pread", Module: "os", NumArgs: 3},
			// pclose(fp) -> exit_code  (wraps pclose(3))
			"pclose": {Name: "pclose", Module: "os", NumArgs: 1},
			// getenv(name) -> ptr  (wraps getenv(3); returns 0 if unset)
			"getenv": {Name: "getenv", Module: "os", NumArgs: 1},
			// setenv(name, value, overwrite) -> status  (wraps setenv(3))
			"setenv": {Name: "setenv", Module: "os", NumArgs: 3},
		},
		Types: map[string]TokenType{},
	}
}

// createTimeModule creates a time/date utility stdlib module
func createTimeModule() *StdlibModule {
	return &StdlibModule{
		Name: "time",
		Functions: map[string]*StdlibFunction{
			"now":       {Name: "now", Module: "time", NumArgs: 0},             // now() -> unix_timestamp
			"sleep":     {Name: "sleep", Module: "time", NumArgs: 1},         // sleep(seconds) -> status
			"millis":    {Name: "millis", Module: "time", NumArgs: 0},       // millis() -> milliseconds
			"nanos":     {Name: "nanos", Module: "time", NumArgs: 0},         // nanos() -> nanoseconds
			"clock":     {Name: "clock", Module: "time", NumArgs: 0},         // clock() -> clock_ticks
			"gmtime":    {Name: "gmtime", Module: "time", NumArgs: 2},       // gmtime(timestamp, tm_buf) -> void
			"localtime": {Name: "localtime", Module: "time", NumArgs: 2}, // localtime(timestamp, tm_buf) -> void
		},
		Types: map[string]TokenType{},
	}
}

// createRandomModule creates the random number generation module
func createRandomModule() *StdlibModule {
	return &StdlibModule{
		Name: "random",
		Functions: map[string]*StdlibFunction{
			"rand":        {Name: "rand", Module: "random", NumArgs: 0},              // rand() -> random int
			"rand_range":  {Name: "rand_range", Module: "random", NumArgs: 2},   // rand_range(min, max) -> random int in [min, max]
			"seed":        {Name: "seed", Module: "random", NumArgs: 1},              // seed(n) -> void
			"rand_float":  {Name: "rand_float", Module: "random", NumArgs: 0},   // rand_float() -> random float [0.0, 1.0)
			"rand_bool":   {Name: "rand_bool", Module: "random", NumArgs: 0},     // rand_bool() -> random boolean
			"rand_bytes":  {Name: "rand_bytes", Module: "random", NumArgs: 2},   // rand_bytes(buf, len) -> fills buffer with random bytes
			"shuffle":     {Name: "shuffle", Module: "random", NumArgs: 2},        // shuffle(arr, len) -> shuffles array in place
			"choice":      {Name: "choice", Module: "random", NumArgs: 2},          // choice(arr, len) -> random element from array
			"rand_n":      {Name: "rand_n", Module: "random", NumArgs: 1},           // rand_n(n) -> random int in [0, n)
			"rand_string": {Name: "rand_string", Module: "random", NumArgs: 2}, // rand_string(buf, len) -> random alphanumeric string
		},
		Types: map[string]TokenType{},
	}
}

// createRegexModule creates the regular expression module
// Uses POSIX regex functions from libc (regcomp, regexec, regfree)
func createRegexModule() *StdlibModule {
	return &StdlibModule{
		Name: "regex",
		Functions: map[string]*StdlibFunction{
			"match": {
				Name:    "match",
				Module:  "regex",
				NumArgs: 2,
			}, // match(pattern, string) -> bool - check if pattern matches string
			"find": {
				Name:    "find",
				Module:  "regex",
				NumArgs: 2,
			}, // find(pattern, string) -> int - return position of first match or -1
			"replace": {
				Name:    "replace",
				Module:  "regex",
				NumArgs: 3,
			}, // replace(pattern, replacement, string) -> string - replace first match
			"replace_all": {
				Name:    "replace_all",
				Module:  "regex",
				NumArgs: 3,
			}, // replace_all(pattern, replacement, string) -> string - replace all matches
			"split": {
				Name:    "split",
				Module:  "regex",
				NumArgs: 2,
			}, // split(pattern, string) -> array - split string by pattern
			"find_all": {
				Name:    "find_all",
				Module:  "regex",
				NumArgs: 2,
			}, // find_all(pattern, string) -> array - find all matches
		},
		Types: map[string]TokenType{},
	}
}

// createJSONModule creates the JSON parsing and serialization module
func createJSONModule() *StdlibModule {
	return &StdlibModule{
		Name: "json",
		Functions: map[string]*StdlibFunction{
			"parse": {
				Name:    "parse",
				Module:  "json",
				NumArgs: 1,
			}, // parse(json_string) -> json_value - parse JSON string
			"stringify": {
				Name:    "stringify",
				Module:  "json",
				NumArgs: 1,
			}, // stringify(value) -> string - convert value to JSON string
			"get": {
				Name:    "get",
				Module:  "json",
				NumArgs: 2,
			}, // get(json, key) -> value - get value by key from JSON object
			"get_int": {
				Name:    "get_int",
				Module:  "json",
				NumArgs: 2,
			}, // get_int(json, key) -> int - get integer value by key
			"get_string": {
				Name:    "get_string",
				Module:  "json",
				NumArgs: 2,
			}, // get_string(json, key) -> string - get string value by key
			"get_bool": {
				Name:    "get_bool",
				Module:  "json",
				NumArgs: 2,
			}, // get_bool(json, key) -> bool - get boolean value by key
			"get_array": {
				Name:    "get_array",
				Module:  "json",
				NumArgs: 2,
			}, // get_array(json, key) -> array - get array by key
			"array_len": {
				Name:    "array_len",
				Module:  "json",
				NumArgs: 1,
			}, // array_len(json_array) -> int - get length of JSON array
			"array_get": {
				Name:    "array_get",
				Module:  "json",
				NumArgs: 2,
			}, // array_get(json_array, index) -> value - get element from JSON array
			"is_null": {
				Name:    "is_null",
				Module:  "json",
				NumArgs: 1,
			}, // is_null(json_value) -> bool - check if value is null
			"is_valid": {
				Name:    "is_valid",
				Module:  "json",
				NumArgs: 1,
			}, // is_valid(json_string) -> bool - check if string is valid JSON
		},
		Types: map[string]TokenType{},
	}
}

// createSDL3Module creates the SDL3 graphics/game development stdlib module
func createSDL3Module() *StdlibModule {
	return &StdlibModule{
		Name: "sdl3",
		Functions: map[string]*StdlibFunction{
			"init": {
				Name:    "init",
				Module:  "sdl3",
				NumArgs: 1,
			}, // init(flags) -> bool - initialize SDL3 subsystems
			"quit": {
				Name:    "quit",
				Module:  "sdl3",
				NumArgs: 0,
			}, // quit() - cleanup SDL3
			"create_window": {
				Name:    "create_window",
				Module:  "sdl3",
				NumArgs: 4,
			}, // create_window(title, w, h, flags) -> window* (SDL3: no x, y)
			"destroy_window": {
				Name:    "destroy_window",
				Module:  "sdl3",
				NumArgs: 1,
			}, // destroy_window(window*) - destroy window
			"create_renderer": {
				Name:    "create_renderer",
				Module:  "sdl3",
				NumArgs: 2,
			}, // create_renderer(window*, flags) -> renderer*
			"destroy_renderer": {
				Name:    "destroy_renderer",
				Module:  "sdl3",
				NumArgs: 1,
			}, // destroy_renderer(renderer*) - destroy renderer
			"render_clear": {
				Name:    "render_clear",
				Module:  "sdl3",
				NumArgs: 1,
			}, // render_clear(renderer*) -> bool
			"render_present": {
				Name:    "render_present",
				Module:  "sdl3",
				NumArgs: 1,
			}, // render_present(renderer*) - present rendered frame
			"set_render_draw_color": {
				Name:    "set_render_draw_color",
				Module:  "sdl3",
				NumArgs: 5,
			}, // set_render_draw_color(renderer*, r, g, b, a) -> bool
			"render_draw_line": {
				Name:    "render_draw_line",
				Module:  "sdl3",
				NumArgs: 5,
			}, // render_draw_line(renderer*, x1, y1, x2, y2) -> bool
			"render_draw_rect": {
				Name:    "render_draw_rect",
				Module:  "sdl3",
				NumArgs: 5,
			}, // render_draw_rect(renderer*, x, y, w, h) -> bool
			"render_fill_rect": {
				Name:    "render_fill_rect",
				Module:  "sdl3",
				NumArgs: 5,
			}, // render_fill_rect(renderer*, x, y, w, h) -> bool
			"poll_event": {
				Name:    "poll_event",
				Module:  "sdl3",
				NumArgs: 1,
			}, // poll_event(event*) -> bool - returns true if event, false if none
			"delay": {
				Name:    "delay",
				Module:  "sdl3",
				NumArgs: 1,
			}, // delay(ms) - delay execution
			"get_ticks": {
				Name:    "get_ticks",
				Module:  "sdl3",
				NumArgs: 0,
			}, // get_ticks() -> int - get milliseconds since init
			"create_texture": {
				Name:    "create_texture",
				Module:  "sdl3",
				NumArgs: 5,
			}, // create_texture(renderer, format, access, w, h) -> texture*
			"destroy_texture": {
				Name:    "destroy_texture",
				Module:  "sdl3",
				NumArgs: 1,
			}, // destroy_texture(texture*)
			"update_texture": {
				Name:    "update_texture",
				Module:  "sdl3",
				NumArgs: 4,
			}, // update_texture(texture*, rect*, pixels*, pitch) -> bool
			"render_texture": {
				Name:    "render_texture",
				Module:  "sdl3",
				NumArgs: 3,
			}, // render_texture(renderer*, texture*, src_rect*, dst_rect*) -> bool
			"lock_texture": {
				Name:    "lock_texture",
				Module:  "sdl3",
				NumArgs: 3,
			}, // lock_texture(texture*, rect*, pixels_out*, pitch_out*) -> bool
			"unlock_texture": {
				Name:    "unlock_texture",
				Module:  "sdl3",
				NumArgs: 1,
			}, // unlock_texture(texture*)
			"get_keyboard_state": {
				Name:    "get_keyboard_state",
				Module:  "sdl3",
				NumArgs: 0,
			}, // get_keyboard_state() -> *bool array (SDL_GetKeyboardState)
			"get_mouse_state": {
				Name:    "get_mouse_state",
				Module:  "sdl3",
				NumArgs: 2,
			}, // get_mouse_state(x_out*, y_out*) -> buttons mask
			"get_relative_mouse_state": {
				Name:    "get_relative_mouse_state",
				Module:  "sdl3",
				NumArgs: 2,
			}, // get_relative_mouse_state(dx_out*, dy_out*) -> buttons mask
			"set_relative_mouse_mode": {
				Name:    "set_relative_mouse_mode",
				Module:  "sdl3",
				NumArgs: 2,
			}, // set_relative_mouse_mode(window*, enabled) -> bool
			"warp_mouse": {
				Name:    "warp_mouse",
				Module:  "sdl3",
				NumArgs: 3,
			}, // warp_mouse(window*, x, y)
			"get_perf_counter": {
				Name:    "get_perf_counter",
				Module:  "sdl3",
				NumArgs: 0,
			}, // get_perf_counter() -> uint64 - high-resolution timer
			"get_perf_freq": {
				Name:    "get_perf_freq",
				Module:  "sdl3",
				NumArgs: 0,
			}, // get_perf_freq() -> uint64 - timer frequency (ticks/sec)
			"get_event_type": {
				Name:    "get_event_type",
				Module:  "sdl3",
				NumArgs: 1,
			}, // get_event_type(event*) -> int - read SDL_Event.type field
			"get_scancode": {
				Name:    "get_scancode",
				Module:  "sdl3",
				NumArgs: 1,
			}, // get_scancode(event*) -> int - read SDL_KeyboardEvent.scancode
			"alloc_event": {
				Name:    "alloc_event",
				Module:  "sdl3",
				NumArgs: 0,
			}, // alloc_event() -> event* - allocate 128-byte SDL_Event
			"show_cursor": {
				Name:    "show_cursor",
				Module:  "sdl3",
				NumArgs: 0,
			}, // show_cursor()
			"hide_cursor": {
				Name:    "hide_cursor",
				Module:  "sdl3",
				NumArgs: 0,
			}, // hide_cursor()
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

// indexOf(s, chOrSubstr): if second arg is string literal of length 1, treat as char
// ============================================================================
// Net module (Linux syscalls)
// ============================================================================

// socket(domain, type, protocol)
// connect_ipv4(fd, ip_u32_host, port_host)
// send(fd, buf, len) -> bytes written
// recv(fd, buf, len) -> bytes read
// close(fd)
// ============================================================================
// UDP Networking Support
// ============================================================================

// bind_ipv4(fd, ip_u32_host, port_host) -> 0 on success, negative on error
// listen(fd, backlog) -> 0 on success, negative on error
// accept(fd) -> client fd or negative on error
// setsockopt(fd, level, optname, optval) -> 0 on success
// sendto_ipv4(fd, buf_ptr, buf_len, dest_ip_u32, dest_port) -> bytes sent
// recvfrom(fd, buf_ptr, buf_len) -> bytes received
// Note: This simplified version doesn't return sender info
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
// generateNetBindIPv6 binds a socket to an IPv6 address
// Args: socket_fd, ipv6_addr_ptr (16 bytes or NULL for any), port
// Returns: 0 on success, negative on error
// generateNetSendtoIPv6 sends data to an IPv6 address (for UDP)
// Args: socket_fd, buf_ptr, buf_len, ipv6_addr_ptr, port
// Returns: bytes sent or negative on error
// ============================================================================
// DNS Resolution (simplified)
// ============================================================================
// These functions parse /etc/hosts for simple resolution
// For real DNS queries, a full resolver would be needed

// generateNetResolve resolves a hostname to IPv4 address via /etc/hosts
// Args: hostname_ptr, out_ipv4_ptr (4 bytes)
// Returns: 1 on success, 0 on failure
// generateNetResolveIPv6 is a stub for IPv6 resolution
// For simplicity, returns 0 (not found) - full implementation would need proper DNS
// ============================================================================
// HTTP module - minimal GET over an existing connected socket
// ============================================================================

// get(fd, host_ptr, host_len, path_ptr, path_len, buf_ptr, buf_len) -> bytes read
// post(fd, host_ptr, host_len, path_ptr, path_len, body_ptr, body_len, buf_ptr, buf_len) -> bytes read
// ============================================================================
// HTTP Response Parsing Functions
// ============================================================================

// generateHTTPParseStatus parses HTTP response status code
// Args: response_buffer, buffer_len -> returns status code (e.g., 200, 404) or 0 on error
// generateHTTPGetHeader extracts a header value from HTTP response
// Args: response_buffer, buffer_len, header_name, out_value_ptr
// Returns: length of value or 0 if not found
// generateHTTPGetBody returns pointer to body (after \r\n\r\n)
// Args: response_buffer, buffer_len -> returns pointer to body or 0 if not found
// generateHTTPParseHeaders populates a buffer with header pointers
// Args: response_buffer, buffer_len, out_headers_array
// This is a simplified version that returns count of headers
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
// generateHTTPPoolGet retrieves a connection from the pool
// Args: pool_ptr, host_ptr, port
// Returns: fd if found, -1 if not found
// generateHTTPPoolPut stores a connection in the pool
// Args: pool_ptr, fd, host_ptr, port
// Returns: 1 if stored, 0 if pool full
// generateHTTPPoolClose closes all connections and frees the pool
// Args: pool_ptr
// Returns: number of connections closed
// printfuncs.go wrappers (will use existing implementations from printfuncs.go)
// Note: Printf and Printf in module system delegate to printfuncs implementations
// The io module functions reference the existing printf/fprintf/sprintf implementations

// ========================= Number conversions =============================
// ========================= Memory implementations ==========================

// generateMemSizeof: sizeof(x) -> returns byte size of variable's type if identifier
// memcpy(dst, src, n): returns dst
// memset(dst, value, n): returns dst (value treated as byte)
// malloc(size): implement via mmap(size) and return pointer
// mmap(size): same as malloc
// free(ptr): currently no-op; recommend mem.munmap
// munmap(ptr, size): rax=11, rdi=addr, rsi=len
// IO module wrapper functions - delegate to printfuncs.go implementations
// ============================================================================
// Collections module - first-pass implementations
// ============================================================================

const collectionsHeaderSize = 40 // len(0), cap(8), head(16), tail(24), data ptr(32)

// allocate header + backing store (cap * elemSize bytes) using mmap; cap is preserved in %rbx
// store common header fields and return base pointer in %rax
// Dynamic array (int)
// generateCollectionsArrayIntCapacity returns the current capacity
// generateCollectionsArrayIntGet gets element at index
// Args: array_ptr, index
// Returns: value at index, or 0 if out of bounds
// generateCollectionsArrayIntSet sets element at index
// Args: array_ptr, index, value
// Returns: 1 if success, 0 if out of bounds
// generateCollectionsArrayIntResize resizes the array to new capacity
// Args: array_ptr, new_capacity
// Returns: new array pointer (may be different if reallocated)
// Uses mmap for new allocation, copies data, munmaps old
// generateCollectionsArrayIntReserve ensures capacity >= new_capacity
// Args: array_ptr, min_capacity
// Returns: array pointer (may be new if resized)
// generateCollectionsArrayIntShrink shrinks capacity to match length
// Args: array_ptr
// Returns: new array pointer
// generateCollectionsArrayIntFree frees the array
// Args: array_ptr
// Returns: 0 on success
// Stack (int) - reuse array layout
// Queue / Deque (int)
// Heap (min-heap, int) backed by array layout
// Hash map (int -> int) with hashing, open addressing, and resize (power-of-two cap)
// Hash set (int) with hashing, open addressing, resize
// ============================================================================
// String-key HashMap and HashSet implementations
// ============================================================================
// These use the djb2 hash algorithm for string keys and strcmp for comparison
// Layout is similar to int variants but stores string pointers as keys

// Helper: compute djb2 hash of null-terminated string
// Input: string ptr in specified register
// Output: hash in %rax
// Helper: compare two null-terminated strings
// Input: str1 in rdi, str2 in rsi
// Output: rax = 1 if equal, 0 if not
// generateCollectionsHashmapStrNew creates a string-keyed hashmap
// Args: initial capacity
// Returns: pointer to hashmap structure
// generateCollectionsHashmapStrPut inserts/updates a key-value pair
// Args: map_ptr, string_key, value
// generateCollectionsHashmapStrGet retrieves value for a string key
// Args: map_ptr, string_key
// Returns: value or 0 if not found
// generateCollectionsHashmapStrContains checks if key exists
// Args: map_ptr, string_key
// Returns: 1 if exists, 0 otherwise
// generateCollectionsHashmapStrRemove removes a key
// Args: map_ptr, string_key
// generateCollectionsHashmapStrLen returns number of entries
// generateCollectionsHashmapStrClear clears all entries
// generateCollectionsHashmapStrFree deallocates the hashmap
// ============================================================================
// String HashSet implementations
// ============================================================================

// generateCollectionsHashsetStrNew creates a string hashset
// generateCollectionsHashsetStrAdd adds a string to the set
// generateCollectionsHashsetStrContains checks if string is in set
// generateCollectionsHashsetStrRemove removes a string from set
// generateCollectionsHashsetStrLen returns size of set
// generateCollectionsHashsetStrClear clears the set
// generateCollectionsHashsetStrFree frees the set
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
// generateCollectionsSortedsetIntNew creates a new sorted set
// Returns: pointer to set structure (root=NULL, count=0)
// generateCollectionsSortedsetIntAdd adds a value to the sorted set
// Args: set_ptr, value
// generateCollectionsSortedsetIntContains checks if value exists
// Args: set_ptr, value -> returns 1 if found, 0 otherwise
// generateCollectionsSortedsetIntRemove removes a value (simplified - marks as deleted)
// For simplicity, we don't actually restructure the tree
// generateCollectionsSortedsetIntMin finds the minimum value
// generateCollectionsSortedsetIntMax finds the maximum value
// generateCollectionsSortedsetIntLen returns the count
// generateCollectionsSortedsetIntFree frees the set (simplified - just frees header)
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
// generateCollectionsSortedmapIntPut inserts or updates a key-value pair
// generateCollectionsSortedmapIntGet retrieves value for key
// generateCollectionsSortedmapIntContains checks if key exists
// generateCollectionsSortedmapIntRemove removes a key (simplified)
// generateCollectionsSortedmapIntMinKey returns the minimum key
// generateCollectionsSortedmapIntMaxKey returns the maximum key
// generateCollectionsSortedmapIntLen returns the count
// generateCollectionsSortedmapIntFree frees the map header
// Binary search helper (int) expects args: base pointer, length, target
// ============================================================================
// Hash module implementations
// ============================================================================

// generateHashCRC32 computes CRC32 checksum (IEEE 802.3 polynomial)
// Args: data_ptr, len -> returns uint32 in rax
// generateHashFNV1a computes FNV-1a hash (64-bit)
// Args: data_ptr, len -> returns uint64 in rax
// generateHashDJB2 computes DJB2 hash (simple string hash)
// Args: string_ptr -> returns uint64 in rax
// generateHashMurmur3 computes MurmurHash3 (32-bit)
// Args: data_ptr, len, seed -> returns uint32 in rax
// generateHashSHA256 computes SHA-256 hash
// Args: data_ptr, len, out_buf (32 bytes) -> void
// Implements the full SHA-256 algorithm per FIPS 180-4
// generateHashMD5 computes MD5 hash
// Args: data_ptr, len, out_buf (16 bytes) -> void
// Implements the full MD5 algorithm per RFC 1321
// ============================================================================
// String Extension Functions (Phase 4)
// ============================================================================

// generateStringSubstring(s, start, len) -> new allocated string
// generateStringSplit(str, delim) -> array ptr
// Returns a pointer to a structure: [count:i64][ptr1][ptr2]...[ptrN]
// Each pointer points to an allocated substring
// generateStringJoin(array, sep) -> joined string
// array format: [count:i64][ptr1][ptr2]...[ptrN]
// generateStringReplace(str, old, new) -> new string with all occurrences replaced
// generateStringToLower(str) -> lowercased copy
// generateStringToUpper(str) -> uppercased copy
// generateStringTrim(str) -> trimmed string (whitespace removed from both ends)
// ============================================================================
// File I/O Functions (Phase 4)
// ============================================================================

// generateFileOpen(path_ptr, flags) -> fd or error
// generateFileClose(fd) -> status
// generateFileRead(fd, buf_ptr, size) -> bytes_read
// generateFileWrite(fd, buf_ptr, size) -> bytes_written
// generateFileSeek(fd, offset, whence) -> new_pos
// generateFileStat(path_ptr, stat_buf) -> status
// generateFileExists(path_ptr) -> 0 or 1
// ============================================================================
// Time Functions (Phase 4)
// ============================================================================

// generateTimeNow() -> unix timestamp
// generateTimeSleep(seconds) -> status
// generateTimeMillis() -> current milliseconds since epoch
// generateTimeNanos() -> current nanoseconds (monotonic for timing)
// generateTimeClock() -> CPU clock ticks (via RDTSC instruction)
// generateTimeGMTime(timestamp, tm_buf) -> void
// Converts Unix timestamp to broken-down time (UTC)
// tm_buf layout (all i64): [sec][min][hour][mday][mon][year][wday][yday][isdst]
// generateTimeLocalTime(timestamp, tm_buf) -> void
// For simplicity, this is the same as gmtime (no timezone support)
// ============================================================================
// RANDOM MODULE - Code Generation Stubs (legacy GCC backend)
// ============================================================================

// generateRandomRand() -> random int
// generateRandomRandRange(min, max) -> random int in [min, max]
// generateRandomSeed(n) -> void
// generateRandomRandFloat() -> random float [0.0, 1.0)
// generateRandomRandBool() -> random boolean
// generateRandomRandBytes(buf, len) -> fills buffer with random bytes
// generateRandomShuffle(arr, len) -> shuffles array in place
// generateRandomChoice(arr, len) -> random element from array
// generateRandomRandN(n) -> random int in [0, n)
// generateRandomRandString(buf, len) -> random alphanumeric string
// ============================================================================
// REGEX MODULE - Code Generation Stubs (legacy GCC backend)
// ============================================================================

// generateRegexMatch(pattern, string) -> bool - check if pattern matches string
// generateRegexFind(pattern, string) -> int - return position of first match or -1
// generateRegexReplace(pattern, replacement, string) -> string - replace first match
// generateRegexReplaceAll(pattern, replacement, string) -> string - replace all matches
// generateRegexSplit(pattern, string) -> array - split string by pattern
// generateRegexFindAll(pattern, string) -> array - find all matches
// ============================================================================
// JSON MODULE - Code Generation Stubs (legacy GCC backend)
// ============================================================================

// generateJSONParse(json_string) -> json_value - parse JSON string
// generateJSONStringify(value) -> string - convert value to JSON string
// generateJSONGet(json, key) -> value - get value by key from JSON object
// generateJSONGetInt(json, key) -> int - get integer value by key
// generateJSONGetString(json, key) -> string - get string value by key
// generateJSONGetBool(json, key) -> bool - get boolean value by key
// generateJSONGetArray(json, key) -> array - get array by key
// generateJSONArrayLen(json_array) -> int - get length of JSON array
// generateJSONArrayGet(json_array, index) -> value - get element from JSON array
// generateJSONIsNull(json_value) -> bool - check if value is null
// generateJSONIsValid(json_string) -> bool - check if string is valid JSON
// ============================================================================
// SDL3 MODULE FUNCTIONS (Assembly backend stubs)
// ============================================================================

// generateSDL3Init(flags) -> bool - initialize SDL3 subsystems
// generateSDL3Quit() - cleanup SDL3
// generateSDL3CreateWindow(title, w, h, flags) -> window* (SDL3: no x, y)
// generateSDL3DestroyWindow(window*) - destroy window
// generateSDL3CreateRenderer(window*, flags) -> renderer*
// generateSDL3DestroyRenderer(renderer*) - destroy renderer
// generateSDL3RenderClear(renderer*) -> bool
// generateSDL3RenderPresent(renderer*) - present rendered frame
// generateSDL3SetRenderDrawColor(renderer*, r, g, b, a) -> bool
// generateSDL3RenderDrawLine(renderer*, x1, y1, x2, y2) -> bool
// generateSDL3RenderDrawRect(renderer*, x, y, w, h) -> bool
// generateSDL3RenderFillRect(renderer*, x, y, w, h) -> bool
// generateSDL3PollEvent(event*) -> bool - returns true if event, false if none
// generateSDL3Delay(ms) - delay execution
// generateSDL3GetTicks() -> int - get milliseconds since init
// ============================================================================
// MATH TRIG/FLOAT FUNCTIONS (Assembly backend: libm-backed fixed-point)
// ============================================================================

// ============================================================================
// SDL3 EXTENDED FUNCTIONS (Assembly backend stubs)
// ============================================================================

// ============================================================================
// SDL_MIXER MODULE
// ============================================================================

// createSDLMixerModule creates the SDL_mixer audio stdlib module
func createSDLMixerModule() *StdlibModule {
	return &StdlibModule{
		Name: "sdl_mixer",
		Functions: map[string]*StdlibFunction{
			"open": {
				Name:    "open",
				Module:  "sdl_mixer",
				NumArgs: 4,
			}, // open(freq, format, channels, chunksize) -> bool
			"close": {
				Name:    "close",
				Module:  "sdl_mixer",
				NumArgs: 0,
			}, // close()
			"load_wav": {
				Name:    "load_wav",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // load_wav(path) -> chunk*
			"free_chunk": {
				Name:    "free_chunk",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // free_chunk(chunk*)
			"play_channel": {
				Name:    "play_channel",
				Module:  "sdl_mixer",
				NumArgs: 3,
			}, // play_channel(channel, chunk*, loops) -> int
			"halt_channel": {
				Name:    "halt_channel",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // halt_channel(channel)
			"volume": {
				Name:    "volume",
				Module:  "sdl_mixer",
				NumArgs: 2,
			}, // volume(channel, vol) -> old_vol
			"load_mus": {
				Name:    "load_mus",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // load_mus(path) -> music*
			"free_music": {
				Name:    "free_music",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // free_music(music*)
			"play_music": {
				Name:    "play_music",
				Module:  "sdl_mixer",
				NumArgs: 2,
			}, // play_music(music*, loops) -> bool
			"halt_music": {
				Name:    "halt_music",
				Module:  "sdl_mixer",
				NumArgs: 0,
			}, // halt_music()
			"volume_music": {
				Name:    "volume_music",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // volume_music(vol) -> old_vol
			"playing": {
				Name:    "playing",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // playing(channel) -> bool
			"playing_music": {
				Name:    "playing_music",
				Module:  "sdl_mixer",
				NumArgs: 0,
			}, // playing_music() -> bool
			"allocate_channels": {
				Name:    "allocate_channels",
				Module:  "sdl_mixer",
				NumArgs: 1,
			}, // allocate_channels(n) -> int
		},
		Types: map[string]TokenType{},
	}
}

// generateOSExec(cmd) -> exit_code
// Calls system(3): runs cmd via /bin/sh, returns the shell exit status.
// generateOSPopen(cmd, mode) -> FILE*
// Calls popen(3): opens a pipe to cmd.  mode must be a pointer to "r" or "w".
// Returns the FILE* as an integer (0 on failure).
// generateOSPread(fp, buf, size) -> bytes_read
// Calls fread(buf, 1, size, fp): reads up to size bytes from a popen pipe.
// generateOSPclose(fp) -> exit_code
// Calls pclose(3): closes the pipe and returns the exit status of the command.
// generateOSGetenv(name) -> ptr
// Calls getenv(3): returns a pointer to the value string, or 0 if unset.
// generateOSSetenv(name, value, overwrite) -> status
// Calls setenv(3): sets env var name=value.  overwrite != 0 replaces existing.
// Returns 0 on success, -1 on error.
// ============================================================================
// FUNC MODULE - Functional programming utilities
// ============================================================================

// createFuncModule creates the functional programming stdlib module.
// These are higher-level combinators that complement the built-in
// Option/Result method syntax (.map, .and_then, etc.).
func createFuncModule() *StdlibModule {
	return &StdlibModule{
		Name: "func",
		Functions: map[string]*StdlibFunction{
			// identity(x) -> x
			"identity": {
				Name:    "identity",
				Module:  "func",
				NumArgs: 1,
			},
			// compose(f, g, x) applies f(x) then g to the result: g(f(x))
			// For ergonomic composition use the |> pipe operator instead.
			"compose": {
				Name:    "compose",
				Module:  "func",
				NumArgs: 3,
			},
			// flip(f, a, b) -> f(b, a)  (swap the first two arguments of f)
			"flip": {
				Name:    "flip",
				Module:  "func",
				NumArgs: 3,
			},
			// const_fn(x, _) -> x  (always returns first argument, ignores second)
			"const_fn": {
				Name:    "const_fn",
				Module:  "func",
				NumArgs: 2,
			},
			// apply(f, x) -> f(x)
			"apply": {
				Name:    "apply",
				Module:  "func",
				NumArgs: 2,
			},
			// ap is an alias for apply(f, x) to mirror applicative terminology.
			"ap": {
				Name:    "ap",
				Module:  "func",
				NumArgs: 2,
			},
			// mempty() -> 0  (additive identity for int monoid)
			"mempty": {
				Name:    "mempty",
				Module:  "func",
				NumArgs: 0,
			},
			// mappend(a, b) -> a + b  (additive int monoid combine)
			"mappend": {
				Name:    "mappend",
				Module:  "func",
				NumArgs: 2,
			},
			// negate(f, x) -> 1 if f(x)==0, else 0
			"negate": {
				Name:    "negate",
				Module:  "func",
				NumArgs: 2,
			},
			// both(f, g, x) -> 1 if f(x)!=0 and g(x)!=0, else 0
			"both": {
				Name:    "both",
				Module:  "func",
				NumArgs: 3,
			},
			// either_fn(f, g, x) -> 1 if f(x)!=0 or g(x)!=0, else 0
			"either_fn": {
				Name:    "either_fn",
				Module:  "func",
				NumArgs: 3,
			},
			// on(f, g, a, b) -> f(g(a), g(b))
			"on": {
				Name:    "on",
				Module:  "func",
				NumArgs: 4,
			},
			// converge(f, g, h, x) -> f(g(x), h(x))
			"converge": {
				Name:    "converge",
				Module:  "func",
				NumArgs: 4,
			},
			// clamp(lo, hi, x) -> max(lo, min(hi, x))
			"clamp": {
				Name:    "clamp",
				Module:  "func",
				NumArgs: 3,
			},
		},
	}
}
