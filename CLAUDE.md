# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Elvish is a modern shell implemented in Go, featuring both an interactive REPL and a scripting language. The project prioritizes readability and maintains comprehensive documentation.

## Development Commands

### Building
```bash
# Build the main elvish binary
make get
# Or manually:
go install ./cmd/elvish

# Build specific variants
go install ./cmd/withpprof/elvish    # With profiling support
go install ./cmd/nodaemon/elvish     # Without daemon functionality

# Build from upstream (if not working from source tree)
go install src.elv.sh/cmd/elvish@latest
```

### Testing
```bash
# Run all unit tests (with race detection if supported)
make test

# Generate and view test coverage report
make cover

# Run specific package tests
go test ./pkg/eval/...
go test ./pkg/parse/...

# Run E2E tests only
go test ./e2e/...

# Scale test timeouts for slower environments
ELVISH_TEST_TIME_SCALE=10 make test
```

### Code Quality
```bash
# Run all checks except gen check
make most-checks

# Run all checks including gen check (requires clean working tree)
make all-checks

# Format code
make fmt

# Individual checks
go vet ./...
staticcheck ./...
codespell
```

### Development Workflow
```bash
# Typical development cycle
make fmt                    # Format code
make test                   # Run tests
make most-checks           # Run quality checks
```

## Architecture

Elvish follows a layered architecture with clear separation of concerns:

### Core Packages (pkg/)
- **pkg/shell**: Main entry point for terminal interface, handles interactive vs non-interactive modes
- **pkg/eval**: The heart of Elvish - language interpreter implementing the evaluator and builtin functions
- **pkg/parse**: Handwritten recursive descent parser for Elvish syntax
- **pkg/edit**: Interactive line editor built on pkg/cli, provides REPL functionality
- **pkg/daemon**: Storage daemon for persistent data (command history, directory history)
- **pkg/cli**: Low-level TUI components and terminal handling

### Supporting Packages
- **pkg/eval/vals**: Standard operations for Elvish values (comparable to a runtime type system)
- **pkg/persistent**: Immutable data structures for lists and maps (inspired by Clojure)
- **pkg/mods/**: Builtin modules (file, os, str, math, etc.)
- **pkg/prog**: Program composition framework used by entrypoints

### Entrypoints (cmd/)
- **cmd/elvish**: Standard build (default)
- **cmd/withpprof/elvish**: With profiling support
- **cmd/nodaemon/elvish**: No daemon functionality
- **cmd/elvmdfmt**: Markdown formatter utility

### Testing Infrastructure
- **pkg/transcript**: Framework for `.elvts` transcript tests that simulate REPL sessions
- **pkg/evaltest**: Testing utilities specific to evaluator testing
- **pkg/testutil**: General testing utilities and helpers

### Language Implementation Details
The evaluator works in three phases:
1. **Parse**: Convert source code to AST using pkg/parse
2. **Compile**: Transform AST into an "operation tree" 
3. **Execute**: Run the operation tree

Language semantics are primarily implemented in:
- `compile_*.go` files (compilation phase)
- `builtin_fn_*.go` files (builtin functions)

## Code Conventions

### Testing Patterns
- Use `.elvts` transcript tests for module functionality - these simulate REPL sessions
- VS Code extension supports Alt-Enter to update transcript outputs
- Use `testutil.Set()` for dependency injection in tests
- Export test dependencies via `testexport_test.go` files when needed
- Respect existing test patterns within each package

### File Organization
- Internal packages follow Go conventions
- `testexport_test.go` files export internal symbols for external test packages  
- `doc.go` files document package architecture and usage

### Module Structure
- Module name is `src.elv.sh` (alias for github.com/elves/elvish)
- All imports start with `src.elv.sh/pkg/...`
- Use last component of package path when referencing symbols (e.g., `eval.Evaler`)

### Build Features
- CGO disabled by default for compatibility (prebuilt binaries)
- Plugin support requires CGO and is available on limited platforms
- Use `CGO_ENABLED=1` to force CGO when building with plugin support

## Development Status & TODO Items

### Current Development Status
- Latest commit: 2025-02-28 (pkg/mods/epm: Add sourcehut to default domain list)
- Active development with regular updates and bug fixes
- Comprehensive TODO tracking via inline code comments

### Known Development Tasks
See [TODO.md](./TODO.md) for a comprehensive list of planned improvements, including:

**High Priority**:
- TUI stack rewrite (pkg/edit, pkg/cli) - major architectural change planned
- Module system enhancements for better package/workspace support
- Platform compatibility improvements, especially Windows support

**Core Features**:
- String module completions (missing Go stdlib bindings)
- Numeric operations improvements
- Enhanced error handling throughout codebase

**Performance**:
- Compilation phase optimizations
- Concurrency safety improvements
- Rendering performance enhancements

### Common Development Areas
When contributing, focus on:
- Windows compatibility (many `*_windows.go` files need work)
- Test coverage improvements (especially transcript tests)
- Error message quality and debugging support
- Performance optimizations in core evaluation loop

## Important Notes

- Project is pre-1.0, expect breaking changes
- Race detector support varies by platform - see `tools/run-race.elv` for supported combinations
- Daemon process is launched on-demand for interactive shells and terminates with the last shell
- Storage backend uses bbolt for persistence (may change in future)
- Time-sensitive tests can be scaled with `ELVISH_TEST_TIME_SCALE` environment variable
- TODO items are tracked as inline comments throughout the codebase rather than centralized tracking