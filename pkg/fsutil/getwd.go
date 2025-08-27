package fsutil

import (
	"os"
	"runtime"
	"strings"
)

// Getwd returns path of the working directory in a format suitable as the
// prompt.
func Getwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "?"
	}
	return TildeAbbr(pwd)
}

// TildeAbbr abbreviates the user's home directory to ~.
func TildeAbbr(path string) string {
	home, err := GetHome("")
	if home == "" || home == "/" {
		// If home is "" or "/", do not abbreviate because (1) it is likely a
		// problem with the environment and (2) it will make the path actually
		// longer.
		return path
	}
	if err == nil {
		// Normalize both paths to forward slashes for consistent comparison on Windows
		normalizedHome := home
		normalizedPath := path
		if runtime.GOOS == "windows" {
			normalizedHome = strings.ReplaceAll(home, "\\", "/")
			normalizedPath = strings.ReplaceAll(path, "\\", "/")
		}
		
		if normalizedPath == normalizedHome {
			return "~"
		} else if strings.HasPrefix(normalizedPath, normalizedHome+"/") {
			return "~" + normalizedPath[len(normalizedHome):]
		}
	}
	return path
}

// TildeAbbrNative abbreviates the user's home directory to ~ with native path separators.
// This is specifically for UI display where native separators are preferred (like location mode).
func TildeAbbrNative(path string) string {
	home, err := GetHome("")
	if home == "" || home == "/" {
		return path
	}
	if err == nil {
		// Normalize both paths to forward slashes for consistent comparison on Windows
		normalizedHome := home
		normalizedPath := path
		if runtime.GOOS == "windows" {
			normalizedHome = strings.ReplaceAll(home, "\\", "/")
			normalizedPath = strings.ReplaceAll(path, "\\", "/")
		}
		
		if normalizedPath == normalizedHome {
			return "~"
		} else if strings.HasPrefix(normalizedPath, normalizedHome+"/") {
			relativePath := normalizedPath[len(normalizedHome):]
			// On Windows, convert back to backslashes for native display
			if runtime.GOOS == "windows" {
				relativePath = strings.ReplaceAll(relativePath, "/", "\\")
			}
			return "~" + relativePath
		}
	}
	return path
}
