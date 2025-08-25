package vars

import (
	"errors"
	"os"
	"runtime"
	"strings"

	"src.elv.sh/pkg/env"
)

var errEnvMustBeString = errors.New("environment variable can only be set string values")

type envVariable struct {
	name string
}

func (ev envVariable) Set(val any) error {
	if s, ok := val.(string); ok {
		os.Setenv(ev.name, s)
		return nil
	}
	return errEnvMustBeString
}

func (ev envVariable) Get() any {
	value := os.Getenv(ev.name)
	// Normalize HOME environment variable on Windows for cross-platform consistency
	if ev.name == env.HOME && runtime.GOOS == "windows" && value != "" {
		return strings.ReplaceAll(value, "\\", "/")
	}
	return value
}

func (ev envVariable) Unset() error {
	return os.Unsetenv(ev.name)
}

func (ev envVariable) IsSet() bool {
	_, ok := os.LookupEnv(ev.name)
	return ok
}

// FromEnv returns a Var corresponding to the named environment variable.
func FromEnv(name string) UnsettableVar {
	return envVariable{name}
}
