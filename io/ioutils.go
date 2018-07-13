package io

import (
    "runtime"
    "os"
    "path/filepath"
    "path"
)

const (
    WINDOWS = "windows"
    LINUX = "linux"
    DARWIN = "darwin"
    UNKNOWN = "unknown"
)

// Returns the root path to the files in terms of this executable
func GetRoot() (string, error) {
    ex, err := os.Executable()
    if err != nil {
        return "", err
    }

    fileInfo, err := os.Lstat(ex)
    if err != nil {
        return "", err
    }

    if fileInfo.Mode()&os.ModeSymlink != 0 {
        newPath, err := os.Readlink(ex)
        if err != nil {
            return "", nil
        }

        newPath = filepath.Dir(newPath)

        // check if the path is relative
        if !filepath.IsAbs(newPath) {
            oldPath := filepath.Dir(ex)
            ex, err = filepath.Abs(path.Join(oldPath, newPath))
            if err != nil {
                return "", err
            }
        } else {
            ex = newPath
        }
    } else {
        ex = filepath.Dir(ex)
    }

    return ex, nil
}

// Returns operating system from three types (windows, darwin, and linux)
func GetOS() string {
    goos := runtime.GOOS

    if goos == WINDOWS {
        return WINDOWS
    } else if goos == DARWIN {
        return DARWIN
    } else if goos == LINUX {
        return LINUX
    } else {
        return UNKNOWN
    }
}

