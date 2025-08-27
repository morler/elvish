// Package str exposes functionality from Go's strings package as an Elvish
// module.
package str

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/eval/errs"
	"src.elv.sh/pkg/eval/vals"
)

var Ns = eval.BuildNsNamed("str").
	AddGoFns(map[string]any{
		"compare":         strings.Compare,
		"contains":        strings.Contains,
		"contains-any":    strings.ContainsAny,
		"count":           strings.Count,
		"equal-fold":      strings.EqualFold,
		"fields":          strings.Fields,
		"fields-func":     fieldsFunc,
		"from-codepoints": fromCodepoints,
		"from-utf8-bytes": fromUtf8Bytes,
		"has-prefix":      strings.HasPrefix,
		"has-suffix":      strings.HasSuffix,
		"index":           strings.Index,
		"index-any":       strings.IndexAny,
		"index-func":      indexFunc,
		"join":            join,
		"last-index":      strings.LastIndex,
		"last-index-func": lastIndexFunc,
		"map":             mapFunc,
		"repeat":          repeat,
		"replace":         replace,
		"split":           split,
		"split-after":     splitAfter,
		//lint:ignore SA1019 Elvish builtins need to be formally deprecated
		// before removal
		"title":            strings.Title,
		"to-codepoints":    toCodepoints,
		"to-lower":         strings.ToLower,
		"to-title":         strings.ToTitle,
		"to-upper":         strings.ToUpper,
		"to-utf8-bytes":    toUtf8Bytes,
		"to-lower-special": toLowerSpecial,
		"to-title-special": toTitleSpecial,
		"to-upper-special": toUpperSpecial,
		"trim":             strings.Trim,
		"trim-left":        strings.TrimLeft,
		"trim-right":       strings.TrimRight,
		"trim-left-func":   trimLeftFunc,
		"trim-right-func":  trimRightFunc,
		"trim-space":       strings.TrimSpace,
		"trim-prefix":      strings.TrimPrefix,
		"trim-suffix":      strings.TrimSuffix,
	}).Ns()

func fromCodepoints(nums ...int) (string, error) {
	var b bytes.Buffer
	for _, num := range nums {
		if num < 0 || num > unicode.MaxRune {
			return "", errs.OutOfRange{
				What:     "codepoint",
				ValidLow: "0", ValidHigh: strconv.Itoa(unicode.MaxRune),
				Actual: hex(num),
			}
		}
		if !utf8.ValidRune(rune(num)) {
			return "", errs.BadValue{
				What:   "argument to str:from-codepoints",
				Valid:  "valid Unicode codepoint",
				Actual: hex(num),
			}
		}
		b.WriteRune(rune(num))
	}
	return b.String(), nil
}

func hex(i int) string {
	if i < 0 {
		return "-0x" + strconv.FormatInt(-int64(i), 16)
	}
	return "0x" + strconv.FormatInt(int64(i), 16)
}

func fromUtf8Bytes(nums ...int) (string, error) {
	var b bytes.Buffer
	for _, num := range nums {
		if num < 0 || num > 255 {
			return "", errs.OutOfRange{
				What:     "byte",
				ValidLow: "0", ValidHigh: "255",
				Actual: strconv.Itoa(num),
			}
		}
		b.WriteByte(byte(num))
	}
	if !utf8.Valid(b.Bytes()) {
		return "", errs.BadValue{
			What:   "arguments to str:from-utf8-bytes",
			Valid:  "valid UTF-8 sequence",
			Actual: fmt.Sprint(b.Bytes()),
		}
	}
	return b.String(), nil
}

func join(sep string, inputs eval.Inputs) (string, error) {
	var buf bytes.Buffer
	var errJoin error
	first := true
	inputs(func(v any) {
		if errJoin != nil {
			return
		}
		if s, ok := v.(string); ok {
			if first {
				first = false
			} else {
				buf.WriteString(sep)
			}
			buf.WriteString(s)
		} else {
			errJoin = errs.BadValue{
				What: "input to str:join", Valid: "string", Actual: vals.Kind(v),
			}
		}
	})
	return buf.String(), errJoin
}

func repeat(s string, n int) (string, error) {
	if n < 0 {
		return "", errs.BadValue{What: "n", Valid: "non-negative number", Actual: vals.ToString(n)}
	}
	if len(s)*n < 0 {
		return "", errs.BadValue{What: "n", Valid: "small enough not to overflow result", Actual: vals.ToString(n)}
	}
	return strings.Repeat(s, n), nil
}

type maxOpt struct{ Max int }

func (o *maxOpt) SetDefaultOptions() { o.Max = -1 }

func replace(opts maxOpt, old, repl, s string) string {
	return strings.Replace(s, old, repl, opts.Max)
}

func split(fm *eval.Frame, opts maxOpt, sep, s string) error {
	out := fm.ValueOutput()
	parts := strings.SplitN(s, sep, opts.Max)
	for _, p := range parts {
		err := out.Put(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func toCodepoints(fm *eval.Frame, s string) error {
	out := fm.ValueOutput()
	for _, r := range s {
		err := out.Put("0x" + strconv.FormatInt(int64(r), 16))
		if err != nil {
			return err
		}
	}
	return nil
}

func toUtf8Bytes(fm *eval.Frame, s string) error {
	out := fm.ValueOutput()
	for _, r := range []byte(s) {
		err := out.Put("0x" + strconv.FormatInt(int64(r), 16))
		if err != nil {
			return err
		}
	}
	return nil
}

// fieldsFunc splits the string s at each run of Unicode code points c satisfying f(c) and returns an array of slices of s.
// If all code points in s satisfy f(c) or the string is empty, an empty slice is returned.
func fieldsFunc(fm *eval.Frame, f eval.Callable, s string) error {
	out := fm.ValueOutput()
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return callFuncForBool(fm, f, string(r))
	})
	for _, field := range fields {
		err := out.Put(field)
		if err != nil {
			return err
		}
	}
	return nil
}

// indexFunc returns the index into s of the first Unicode code point satisfying f(c), or -1 if none do.
func indexFunc(fm *eval.Frame, f eval.Callable, s string) (int, error) {
	index := strings.IndexFunc(s, func(r rune) bool {
		return callFuncForBool(fm, f, string(r))
	})
	return index, nil
}

// lastIndexFunc returns the index into s of the last Unicode code point satisfying f(c), or -1 if none do.
func lastIndexFunc(fm *eval.Frame, f eval.Callable, s string) (int, error) {
	index := strings.LastIndexFunc(s, func(r rune) bool {
		return callFuncForBool(fm, f, string(r))
	})
	return index, nil
}

// mapFunc returns a copy of the string s with all its characters modified according to the mapping function.
// If mapping returns a negative value, the character is dropped from the string with no replacement.
func mapFunc(fm *eval.Frame, f eval.Callable, s string) (string, error) {
	result := strings.Map(func(r rune) rune {
		return callFuncForRune(fm, f, string(r), r)
	}, s)
	return result, nil
}

// splitAfter slices s into all substrings after each instance of sep and returns a slice of those substrings.
func splitAfter(fm *eval.Frame, opts maxOpt, sep, s string) error {
	out := fm.ValueOutput()
	parts := strings.SplitAfterN(s, sep, opts.Max)
	// Filter out empty strings at the end (same behavior as strings.Split)
	for _, p := range parts {
		if p != "" {
			err := out.Put(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// toLowerSpecial returns a copy of the string s with all Unicode letters mapped to their lower case using the case mapping specified by c.
func toLowerSpecial(c, s string) (string, error) {
	casing, err := parseCasing(c)
	if err != nil {
		return "", err
	}
	return strings.ToLowerSpecial(casing, s), nil
}

// toTitleSpecial returns a copy of the string s with all Unicode letters mapped to their title case, giving priority to the special casing rules.
func toTitleSpecial(c, s string) (string, error) {
	casing, err := parseCasing(c)
	if err != nil {
		return "", err
	}
	return strings.ToTitleSpecial(casing, s), nil
}

// toUpperSpecial returns a copy of the string s with all Unicode letters mapped to their upper case using the case mapping specified by c.
func toUpperSpecial(c, s string) (string, error) {
	casing, err := parseCasing(c)
	if err != nil {
		return "", err
	}
	return strings.ToUpperSpecial(casing, s), nil
}

// parseCasing converts casing names to unicode.SpecialCase constants
func parseCasing(c string) (unicode.SpecialCase, error) {
	switch c {
	case "turkish", "tr":
		return unicode.TurkishCase, nil
	case "azeri", "az":
		return unicode.AzeriCase, nil
	default:
		return nil, errs.BadValue{
			What:   "casing",
			Valid:  "turkish, tr, azeri, or az",
			Actual: c,
		}
	}
}

// trimLeftFunc returns a slice of the string s with all leading Unicode code points c satisfying f(c) removed.
func trimLeftFunc(fm *eval.Frame, f eval.Callable, s string) (string, error) {
	result := strings.TrimLeftFunc(s, func(r rune) bool {
		return callFuncForBool(fm, f, string(r))
	})
	return result, nil
}

// trimRightFunc returns a slice of the string s with all trailing Unicode code points c satisfying f(c) removed.
func trimRightFunc(fm *eval.Frame, f eval.Callable, s string) (string, error) {
	result := strings.TrimRightFunc(s, func(r rune) bool {
		return callFuncForBool(fm, f, string(r))
	})
	return result, nil
}

// Helper function to call Elvish function and get a boolean result
func callFuncForBool(fm *eval.Frame, f eval.Callable, arg string) bool {
	// Use CaptureOutput to capture the function's output
	outputs, err := fm.CaptureOutput(func(subFm *eval.Frame) error {
		return f.Call(subFm, []any{arg}, eval.NoOpts)
	})

	if err != nil || len(outputs) == 0 {
		return false
	}

	// Convert the first output to boolean
	return vals.Bool(outputs[0])
}

// Helper function to call Elvish function and get a rune result
func callFuncForRune(fm *eval.Frame, f eval.Callable, arg string, originalRune rune) rune {
	// Use CaptureOutput to capture the function's output
	outputs, err := fm.CaptureOutput(func(subFm *eval.Frame) error {
		return f.Call(subFm, []any{arg}, eval.NoOpts)
	})

	if err != nil || len(outputs) == 0 {
		return originalRune // keep original character on error
	}

	// Handle different return types
	switch val := outputs[0].(type) {
	case string:
		if len(val) == 0 {
			return -1 // drop character
		}
		runes := []rune(val)
		return runes[0] // return first rune
	case int:
		if val < 0 || val > unicode.MaxRune {
			return -1 // drop character for invalid runes
		}
		return rune(val)
	default:
		return originalRune // keep original character for other types
	}
}
