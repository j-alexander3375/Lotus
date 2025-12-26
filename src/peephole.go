package main

import (
	"regexp"
	"strings"
)

// peephole.go - Assembly-level peephole optimizations
// This file contains optimizations that run on the generated assembly output.
// Peephole optimizations look at small windows of instructions and replace
// inefficient patterns with more efficient equivalents.

// Pre-compiled regex patterns for better performance
var (
	redundantMovePattern = regexp.MustCompile(`^\s+movq\s+%(\w+),\s+%(\w+)$`)
	storePattern         = regexp.MustCompile(`^\s+movq\s+%(\w+),\s+(-?\d+\(%\w+\))$`)
	loadPattern          = regexp.MustCompile(`^\s+movq\s+(-?\d+\(%\w+\)),\s+%(\w+)$`)
	pushPattern          = regexp.MustCompile(`^\s+pushq\s+%(\w+)$`)
	popPattern           = regexp.MustCompile(`^\s+popq\s+%(\w+)$`)
	zeroLoadPattern      = regexp.MustCompile(`^\s+movq\s+\$0,\s+%(\w+)$`)
	addOneRegPattern     = regexp.MustCompile(`^\s+addq\s+\$1,\s+%(\w+)$`)
	subOneRegPattern     = regexp.MustCompile(`^\s+subq\s+\$1,\s+%(\w+)$`)
	addOneMemPattern     = regexp.MustCompile(`^\s+addq\s+\$1,\s+(-?\d+\(%\w+\))$`)
	subOneMemPattern     = regexp.MustCompile(`^\s+subq\s+\$1,\s+(-?\d+\(%\w+\))$`)
)

// PeepholeOptimizer applies local optimizations to assembly output.
type PeepholeOptimizer struct {
	lines []string
}

// NewPeepholeOptimizer creates a new peephole optimizer.
func NewPeepholeOptimizer(assembly string) *PeepholeOptimizer {
	return &PeepholeOptimizer{
		lines: strings.Split(assembly, "\n"),
	}
}

// Optimize applies all peephole optimization passes.
func (po *PeepholeOptimizer) Optimize() string {
	// Apply passes multiple times until no more changes
	changed := true
	for changed {
		changed = false
		changed = po.removeRedundantMoves() || changed
		changed = po.removeDeadStores() || changed
		changed = po.optimizeStackPushPop() || changed
		changed = po.optimizeZeroLoading() || changed
		changed = po.optimizeAddSubByOne() || changed
	}
	return strings.Join(po.lines, "\n")
}

// removeRedundantMoves removes mov instructions that move a register to itself.
// Example: movq %rax, %rax → (removed)
func (po *PeepholeOptimizer) removeRedundantMoves() bool {
	changed := false
	newLines := make([]string, 0, len(po.lines))
	for _, line := range po.lines {
		if matches := redundantMovePattern.FindStringSubmatch(line); matches != nil {
			if matches[1] == matches[2] {
				changed = true
				continue
			}
		}
		newLines = append(newLines, line)
	}
	po.lines = newLines
	return changed
}

// removeDeadStores removes store-then-load sequences to the same location.
// Example:
//
//	movq %rax, -8(%rbp)
//	movq -8(%rbp), %rax
//
// The second instruction can be removed if %rax hasn't changed.
func (po *PeepholeOptimizer) removeDeadStores() bool {
	changed := false
	newLines := make([]string, 0, len(po.lines))
	for i := 0; i < len(po.lines); i++ {
		line := po.lines[i]

		// Check if this is a store followed by a load from the same location to the same register
		if i+1 < len(po.lines) {
			storeMatch := storePattern.FindStringSubmatch(line)
			loadMatch := loadPattern.FindStringSubmatch(po.lines[i+1])

			if storeMatch != nil && loadMatch != nil {
				storeReg := storeMatch[1]
				storeLoc := storeMatch[2]
				loadLoc := loadMatch[1]
				loadReg := loadMatch[2]

				// If storing to and loading from the same location, and same register
				if storeLoc == loadLoc && storeReg == loadReg {
					// Keep the store, skip the load
					newLines = append(newLines, line)
					i++ // Skip next line
					changed = true
					continue
				}
			}
		}

		newLines = append(newLines, line)
	}
	po.lines = newLines
	return changed
}

// optimizeStackPushPop removes push-pop sequences that cancel out.
// Example:
//
//	pushq %rax
//	popq %rax
//
// Both can be removed.
func (po *PeepholeOptimizer) optimizeStackPushPop() bool {
	changed := false
	newLines := make([]string, 0, len(po.lines))
	for i := 0; i < len(po.lines); i++ {
		line := po.lines[i]

		// Check for immediate push-pop of the same register
		if i+1 < len(po.lines) {
			pushMatch := pushPattern.FindStringSubmatch(line)
			popMatch := popPattern.FindStringSubmatch(po.lines[i+1])

			if pushMatch != nil && popMatch != nil && pushMatch[1] == popMatch[1] {
				i++
				changed = true
				continue
			}
		}

		newLines = append(newLines, line)
	}
	po.lines = newLines
	return changed
}

// optimizeZeroLoading replaces movq $0, %reg with xorq %reg, %reg.
// xorq is smaller and faster for zeroing registers.
func (po *PeepholeOptimizer) optimizeZeroLoading() bool {
	changed := false
	for i, line := range po.lines {
		if matches := zeroLoadPattern.FindStringSubmatch(line); matches != nil {
			reg := matches[1]
			indent := getIndent(line)
			po.lines[i] = indent + "xorq %" + reg + ", %" + reg
			changed = true
		}
	}
	return changed
}

// getIndent extracts the leading whitespace from a line.
func getIndent(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " \t"))]
}

// optimizeAddSubByOne replaces addq $1, %reg with incq %reg (and similar for sub).
// incq/decq are more compact.
func (po *PeepholeOptimizer) optimizeAddSubByOne() bool {
	changed := false
	for i, line := range po.lines {
		indent := getIndent(line)

		if matches := addOneRegPattern.FindStringSubmatch(line); matches != nil {
			po.lines[i] = indent + "incq %" + matches[1]
			changed = true
		} else if matches := subOneRegPattern.FindStringSubmatch(line); matches != nil {
			po.lines[i] = indent + "decq %" + matches[1]
			changed = true
		} else if matches := addOneMemPattern.FindStringSubmatch(line); matches != nil {
			po.lines[i] = indent + "incq " + matches[1]
			changed = true
		} else if matches := subOneMemPattern.FindStringSubmatch(line); matches != nil {
			po.lines[i] = indent + "decq " + matches[1]
			changed = true
		}
	}
	return changed
}

// ApplyPeepholeOptimizations is the main entry point for peephole optimization.
func ApplyPeepholeOptimizations(assembly string) string {
	optimizer := NewPeepholeOptimizer(assembly)
	return optimizer.Optimize()
}
