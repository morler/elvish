# TODO.md

This document summarizes the development tasks and improvement opportunities identified throughout the Elvish codebase.

## High Priority Items

### TUI Stack Rewrite
**Status**: Planned major refactoring  
**Location**: `pkg/edit`, `pkg/cli`  
**Description**: The entire TUI (Terminal User Interface) stack is due for a rewrite. The current editor implementation built on `cli.App` needs modernization.

### Module System Enhancements
**Location**: `pkg/eval/builtin_special.go`  
- Add support for module specs relative to a package/workspace
- Improve module access concurrency safety
- For non-relative imports, use the spec instead of full path

## Core Language Features

### String Module Completions
**Location**: `pkg/mods/str/str.go`  
Missing Go standard library function bindings:
- `FieldsFunc`
- `IndexFunc`, `LastIndexFunc`
- `Map`
- `SplitAfter`
- `ToLowerSpecial`, `ToTitleSpecial`, `ToUpperSpecial`
- `TrimLeft`, `TrimRight`, `TrimLeftFunc`, `TrimRightFunc`

### Numeric Operations
**Location**: `pkg/eval/builtin_fn_num.go`  
- Fix range function default value handling (currently can only be used implicitly)
- Improve numeric type conversion and mixed argument handling

### Error Handling Improvements
Multiple locations need better error handling:
- `pkg/eval/builtin_fn_io.go`: Don't ignore JSON formatting errors
- `pkg/eval/builtin_fn_flow.go`: Add proper multi-error documentation
- `pkg/cli/modes/location.go`: Surface file system errors properly

## Platform Compatibility

### Windows Support
**Priority**: Medium  
**Locations**: Multiple files with `*_windows.go` suffix

#### File System Operations
- `pkg/mods/os/stat_windows.go`: Implement CreationTime, LastAccessTime, LastWriteTime
- `pkg/cli/lscolors/stat_windows.go`: Implement file feature detection
- `pkg/eval/compile_value.go`: Fix path handling correctness on Windows

#### Terminal Support
- `pkg/cli/term/reader_windows.go`: Improve key sequence normalization (currently Unix-centric)

### Cross-Platform Path Handling
**Location**: `pkg/glob/glob.go`  
- Preserve original path separator (/ or \) on Windows
- Handle Windows UNC paths properly
- Improve glob pattern matching on Windows

## Performance Optimizations

### Compilation Phase
**Location**: `pkg/eval/compile_*.go`  
- `compile_effect.go`: Avoid unnecessary type conversions
- `compile_value.go`: Optimize tilde expansion logic
- Improve overall compilation performance (currently not very performant)

### Concurrency Safety
**Location**: `pkg/eval/builtin_special.go`  
- Make access to `fm.Evaler.modules` concurrency-safe
- Improve variable access thread safety

### Rendering Optimizations
**Location**: `pkg/cli/tk/`  
- `label.go`: Optimize rendering by stopping after height rows are written
- `listbox.go`: Fix multi-line item rendering issues
- Improve overall TUI rendering performance

## User Experience Enhancements

### LSP (Language Server Protocol)
**Location**: `pkg/lsp/server.go`  
- Take variable shadowing into consideration for completions and definitions
- Support more completion item kinds beyond current basic set
- Improve overall language server feature set

### Editor Features
**Location**: `pkg/edit/`  
- `completion.go`: Add completion display improvements
- `highlight/regions.go`: Extend highlighting to cover more command types beyond barewords
- `filter/highlight.go`: Add error highlighting support

### Command Completion
**Location**: `pkg/edit/complete_getopt.go`  
- Make Config field configurable
- Improve argument completer notifications
- Better handling of chained options

## Testing and Quality

### Test Coverage Gaps
- `pkg/mods/daemon/daemon_test.go`: Empty test file needs implementation
- `pkg/edit/store_api_test.go`: Add session history testing
- `pkg/glob/glob_test.go`: Add more Lstat failure test cases and dotfile tests
- Various transcript test files need Windows compatibility

### Test Infrastructure
- `pkg/cli/term/read_rune_test.go`: Remove Unix dependency
- `pkg/eval/compiler_test.go`: Convert deterministic tests to fuzz tests
- Add more comprehensive error condition testing

## Documentation and Usability

### Function Documentation
**Location**: Various `builtin_fn_*.go` files  
- `builtin_fn_cmd.go`: Document "fg" command
- `builtin_fn_str.go`: Document `-override-wcswidth` option
- `builtin_fn_flow.go`: Document "multi-error" properly

### Error Messages
- `pkg/eval/builtin_special_test.elvts`: Improve stack traces to point to correct locations
- `pkg/mods/os/os_test.elvts`: Make error messages more informative
- Better error reporting for type mismatches in various modules

## Low Priority / Nice to Have

### Code Organization
- `pkg/eval/`: Move `builtin_fn_*.go` files to a separate package
- `pkg/glob/parse.go`: Eliminate duplicate code with `parse/parser.go`
- Clean up various TODO comments for code structure improvements

### Feature Enhancements
- `pkg/eval/builtin_fn_cmd_unix.go`: Find and display command names for processes
- `pkg/cli/modes/histlist.go`: Improve index alignment for >10000 entries
- `pkg/ui/text_segment.go`: Make string conversion environment-aware (e.g., HTML output)

### Database Migration
**Status**: Under evaluation  
**Location**: `pkg/daemon/`  
The current bbolt-based storage daemon might be replaced with a custom database solution in the future.

## Notes

- Most TODO items are tracked as inline comments in the source code
- Priority should be given to items that affect core functionality or user experience
- Windows support improvements would significantly expand platform compatibility
- The TUI rewrite is likely the largest single undertaking on this list

## Contributing

When working on these items:
1. Check if the TODO is still relevant (some may have been addressed)
2. Consider the impact on existing functionality
3. Add appropriate tests for new features
4. Update documentation as needed
5. Follow established code patterns within each package