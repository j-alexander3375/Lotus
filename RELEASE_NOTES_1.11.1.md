# Release Notes for Lotus v1.11.1

## Highlights

Minor fix to how flags are parsed in compiler.

## Fixes

### CLI
- `ParseFlags` (`src/flags.go`) now supports interspersed flags and positional arguments.
  Previously, Go's `flag.Parse` stopped at the first non-flag token, so a flag placed after
  the source file (e.g. `lotus program.lts -o out`) was silently ignored and `-o`/`out` were
  treated as leftover positional arguments — the compiler fell back to its default output
  path (`a.out`) instead of the requested one.
- Flags may now appear before, after, or interleaved with the input file:
  `lotus -o out program.lts` and `lotus program.lts -o out` behave identically.

## Upgrade Notes
- No source-level changes from 1.11.0. The 1.11.0 upgrade note "Compiler flags must precede
  the input file" no longer applies as of this release.
