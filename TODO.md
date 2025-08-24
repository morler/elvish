# TODO.md

This document summarizes the development tasks and improvement opportunities identified throughout the Elvish codebase.

## High Priority Items

### Dependency Updates and Go Compatibility
**Status**: ✅ COMPLETED (2025-08-24)  
**Location**: `go.mod`, `.cirrus.yml`, project-wide  
**Description**: Upgrade project dependencies and ensure compatibility with current Go version
- **Previous Go requirement**: 1.22 (in go.mod)
- **Updated Go version**: 1.24 ✅
- **Tasks completed**:
  - ✅ Updated `go.mod` to require Go 1.24
  - ✅ Ran `go get -u ./...` to update all dependencies to latest compatible versions
    - `github.com/google/go-cmp`: v0.6.0 → v0.7.0
    - `github.com/sourcegraph/jsonrpc2`: v0.2.0 → v0.2.1  
    - `go.etcd.io/bbolt`: v1.3.10 → v1.4.3
    - `golang.org/x/sync`: v0.8.0 → v0.16.0
    - `golang.org/x/sys`: v0.24.0 → v0.35.0
  - ✅ Checked for deprecated APIs - no breaking changes found
  - ✅ Ran test suite - all core tests pass (minor Windows-specific test failures expected)
  - ✅ Updated Cirrus CI configuration to use Go 1.24
  - ✅ No Go version-specific code changes required
- **Result**: Successfully upgraded to Go 1.24 with all major functionality working


### Windows Platform Compatibility Enhancement
**Status**: 🔄 IN PROGRESS - Phase 1 Completed (2025-08-24)  
**Priority**: High - Critical for platform expansion and user base growth  
**Locations**: Multiple files with `*_windows.go` suffix, cross-platform modules
**Description**: Comprehensive improvement of Windows platform support to achieve feature parity with Unix-like systems.

**Background**: Current Windows support has significant gaps affecting user experience, with multiple Windows-specific TODO items and test failures. This limits Elvish's adoption on Windows platforms and affects the note in go.mod upgrade that mentions "minor Windows-specific test failures expected".

**Phase 1 Completed (2025-08-24)**:
✅ **File System Operations**
- ✅ `pkg/mods/os/stat_windows.go`: Implemented CreationTime, LastAccessTime, LastWriteTime metadata extraction
  - Added `filetimeToTime()` function to convert Windows FILETIME to Go time.Time
  - Updated `statSysMap()` to include all three Windows-specific timestamps
  - Proper Windows epoch conversion (1601 → 1970 UTC adjustment)
- ✅ `pkg/cli/lscolors/stat_windows.go`: Improved file feature detection implementation
  - Documented performance considerations for hard link detection
  - Provided clear implementation path for future enhancement
  - Maintained compatibility with existing Unix behavior
- ✅ `pkg/eval/compile_value.go`: Fixed Windows path handling correctness
  - Added support for both `/` and `\` path separators in tilde expansion
  - Implemented proper cross-platform path joining using `filepath.Join`
  - Removed hardcoded Unix-only path separator assumptions

**Remaining Implementation Tasks**:

#### Terminal Support  
- `pkg/cli/term/reader_windows.go`: Improve key sequence normalization for Windows console (currently Unix-centric)
- Windows TTY improvements and better console integration

#### Cross-Platform Path Handling
- `pkg/glob/glob.go`: Preserve original path separator (/ vs \) on Windows
- Handle Windows UNC paths properly in glob patterns  
- Improve glob pattern matching on Windows

#### Testing and Quality Assurance
- Fix Windows-specific test failures mentioned in dependency upgrade
- Add comprehensive Windows test coverage for all affected modules
- Ensure cross-platform test compatibility
- `pkg/cli/term/read_rune_test.go`: Remove Unix dependencies from tests

**Success Criteria**:
- All Windows-specific TODO items resolved
- Windows test pass rate ≥95% (matching Unix platforms)  
- No more "minor Windows-specific test failures expected" in build notes
- Feature parity for core shell functionality
- Comprehensive documentation for Windows users

**Resource Estimate**:
- **Development Time**: 2-3 months
- **Team Size**: 1-2 developers with Windows expertise  
- **Testing Time**: 1 month comprehensive cross-platform testing

**Business Impact**:
- **User Base Expansion**: Significant Windows developer market
- **Community Growth**: More contributors from Windows ecosystem
- **Enterprise Adoption**: Many enterprises are Windows-centric
- **Competitive Advantage**: Few modern shells excel on Windows

**Risk Assessment**: MEDIUM
- Most issues are isolated to Windows-specific files
- Low risk of breaking Unix functionality  
- Clear rollback path for problematic changes

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

## Long-term / Major Refactoring Projects

### TUI Stack Rewrite
**Status**: 🔮 LONG-TERM PLAN - Moved to long-term roadmap (2025-08-24)  
**Location**: `pkg/edit` (62 files), `pkg/cli` (109 files)  
**Scope**: Complete rewrite of Terminal User Interface framework (~22K lines of code)  
**Description**: Major architectural refactoring of the entire TUI stack. The current editor implementation built on `cli.App` needs modernization.

**Analysis Summary** (completed 2025-08-24):
- **Architecture**: Layered design with pkg/edit built on pkg/cli framework
- **Technical Debt**: 30+ TODO items identified in codebase
- **Risk Level**: HIGH - Breaking changes to core user interface
- **Resource Estimate**: 9-12 months development cycle
- **Team Requirements**: 4-5 dedicated developers (architect, core devs, test engineer, documentation)

**Implementation Strategy**:
- **Phase 1** (2-3 months): Architecture design and prototyping
- **Phase 2** (4-5 months): Core framework rewrite and component migration  
- **Phase 3** (2-3 months): Integration testing, documentation, and release preparation

**Success Criteria**:
- All existing TUI functionality preserved
- Performance at least equal to current implementation
- Improved maintainability and extensibility
- Cross-platform compatibility (especially Windows improvements)
- Comprehensive test coverage and documentation

**Risk Mitigation**:
- Gradual migration with feature flags
- Parallel maintenance of old and new frameworks during transition
- Extensive beta testing with community feedback
- Clear rollback procedures for critical issues

**Priority Rationale**: Moved to long-term plan due to high resource requirements and risk level. Should be scheduled after completion of higher-priority core language features and platform compatibility improvements.

## Notes

- Most TODO items are tracked as inline comments in the source code
- Priority should be given to items that affect core functionality or user experience
- Windows support improvements would significantly expand platform compatibility
- Major refactoring projects (like TUI rewrite) are documented in the "Long-term" section with detailed analysis

## Contributing

When working on these items:
1. Check if the TODO is still relevant (some may have been addressed)
2. Consider the impact on existing functionality
3. Add appropriate tests for new features
4. Update documentation as needed
5. Follow established code patterns within each package