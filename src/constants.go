package main

// constants.go - Compiler-wide constants and configuration values

const (
	// CompilerVersion is the current version of the Lotus compiler
	CompilerVersion = "1.2.3"

	// DefaultStackAlignment is the default stack alignment in bytes (16-byte alignment for x86-64)
	DefaultStackAlignment = 16

	// PointerSize is the size of a pointer in bytes (64-bit architecture)
	PointerSize = 8

	// DefaultStringCapacity is the initial capacity for string builders
	DefaultStringCapacity = 4096
)

// Assembly generation constants
const (
	// DataSectionDirective is the assembly directive for the data section
	DataSectionDirective = ".section .data"

	// TextSectionDirective is the assembly directive for the text/code section
	TextSectionDirective = ".text"

	// GlobalDirective is the assembly directive for global symbols (GNU as)
	GlobalDirective = ".globl"

	// EntryPointLabel is the standard entry point for x86-64 programs
	EntryPointLabel = "_start"
)

// System call numbers for x86-64 Linux
const (
	SyscallWrite = 1  // write(fd, buf, count)
	SyscallExit  = 60 // exit(status)
)

// Standard file descriptors
const (
	StdoutFD = 1 // Standard output file descriptor
	StderrFD = 2 // Standard error file descriptor
)

// Type sizes in bytes
const (
	Int8Size    = 1
	Int16Size   = 2
	Int32Size   = 4
	Int64Size   = 8
	Float32Size = 4
	Float64Size = 8
	BoolSize    = 1
)

// Label prefixes for generated assembly labels
const (
	StringLabelPrefix   = ".str"
	IfLabelPrefix       = ".if"
	ElseLabelPrefix     = ".else"
	EndIfLabelPrefix    = ".endif"
	WhileLabelPrefix    = ".while"
	EndWhileLabelPrefix = ".endwhile"
	ForLabelPrefix      = ".for"
	EndForLabelPrefix   = ".endfor"
	FuncLabelPrefix     = ".func"
	CatchLabelPrefix    = ".catch"
	FinallyLabelPrefix  = ".finally"
	TryEndLabelPrefix   = ".try_end"
)
