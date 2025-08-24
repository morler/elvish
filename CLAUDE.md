# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Elvish is a modern shell implemented in Go, featuring both an interactive REPL and a scripting language. The project prioritizes readability and maintains comprehensive documentation.

## Development Commands

### Building
```bash
# Build the main elvish binary (recommended)
make get
# Or manually:
go install ./cmd/elvish

# Build specific variants
go install ./cmd/withpprof/elvish    # With profiling support
go install ./cmd/nodaemon/elvish     # Without daemon functionality
go install ./cmd/elvmdfmt             # Markdown formatter utility

# Build from upstream (if not working from source tree)
go install src.elv.sh/cmd/elvish@latest

# Note: Requires Go 1.24+ (current project uses Go 1.24)
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
- **Go Version**: Requires Go 1.24+ (project upgraded from earlier versions)
- **CGO**: Disabled by default for compatibility (prebuilt binaries)
- **Plugin Support**: Requires CGO and is available on limited platforms
- **Cross-Platform**: Enhanced Windows compatibility with recent improvements
- **Dependencies**: Updated to latest versions including golang.org/x/sys v0.35.0
- Use `CGO_ENABLED=1` to force CGO when building with plugin support

## Development Status & TODO Items

### Current Development Status
- Latest commit: 656f7fea (feat: Performance optimization - Compilation phase benchmarking infrastructure)
- Recent major improvements:
  - Performance optimization with compilation phase benchmarking infrastructure (656f7fea)
  - LSP Enhancement with variable shadowing support for completions and definitions (0bf3d9b7)
  - Project documentation updates with recent progress and current status (bbbfe951)
  - Comprehensive error handling improvements and reliability enhancements (a97ca688)
  - Complete str module with missing Go stdlib function bindings (cd6cb020)
  - Enhanced module system with concurrency safety and performance optimizations (df9b230c)
  - Windows Platform Compatibility Enhancement - Phase 2 (8a486328)
  - Go version upgraded to 1.24+ with dependency updates (4e0542d6)
- Active development with regular updates and bug fixes
- Comprehensive TODO tracking via inline code comments

### Known Development Tasks
See [TODO.md](./TODO.md) for a comprehensive list of planned improvements, including:

**Recently Completed**:
- ✅ Performance optimization with compilation phase benchmarking infrastructure 
- ✅ LSP Enhancement with variable shadowing support for completions and definitions
- ✅ Comprehensive project documentation updates with progress tracking
- ✅ String module completions - missing Go stdlib function bindings now implemented
- ✅ Comprehensive error handling improvements and reliability enhancements
- ✅ Windows Platform Compatibility Enhancement (Phase 1 & 2)
- ✅ Module system enhanced with concurrency safety and performance optimizations

**High Priority**:
- TUI stack rewrite (pkg/edit, pkg/cli) - major architectural change planned
- Further performance optimizations building on benchmarking infrastructure
- Advanced LSP features beyond variable shadowing
- Further platform compatibility improvements
- Performance optimization in numeric operations

**Core Features**:
- Advanced numeric operations and mathematical functions
- Enhanced testing infrastructure improvements
- Language server protocol enhancements

**Performance**:
- Compilation phase optimizations
- Rendering performance enhancements
- Memory management improvements

### Common Development Areas
When contributing, focus on:
- **Performance optimization**: Leverage new benchmarking infrastructure for systematic improvements
- **LSP development**: Expand on variable shadowing support with additional language server features
- **Windows compatibility**: Significant improvements made, but ongoing work needed
- **Test coverage**: Improvements especially in transcript tests and E2E testing
- **Error handling**: Recently enhanced but continue improving error message quality
- **Performance**: Core evaluation loop optimizations and numerical operations
- **Module system**: Building on recent concurrency safety improvements
- **String operations**: Leverage newly completed Go stdlib function bindings

## Important Notes

- **Version Status**: Project is pre-1.0, expect breaking changes
- **Go Requirements**: Now requires Go 1.24+ (upgraded from earlier versions)
- **Race Detection**: Support varies by platform - see `tools/run-race.elv` for supported combinations
- **Platform Support**: Enhanced Windows compatibility in recent releases
- **Daemon Process**: Launched on-demand for interactive shells, terminates with last shell
- **Storage Backend**: Uses bbolt v1.4.3 for persistence (may change in future)
- **Testing**: Time-sensitive tests can be scaled with `ELVISH_TEST_TIME_SCALE` environment variable
- **TODO Tracking**: Items tracked as inline comments throughout codebase rather than centralized
- **Error Handling**: Recently improved with comprehensive reliability enhancements
- **String Operations**: Now includes complete Go stdlib function bindings
- **Performance Benchmarking**: New infrastructure available for systematic performance measurement and optimization
- **LSP Support**: Enhanced with variable shadowing support for improved development experience