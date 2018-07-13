package template

import (
    "strings"
    "go-utils/fs"
    "go-utils/errors"
)

func IOReplace(path string, values map[string]string) error {
    data, err := fs.ReadFile(path)
    if nil != err {
        return errors.ReadFileError{FileName: path, Err: err}
    }
    result := Replace(string(data), values)
    err = fs.WriteFile(path, []byte(result))
    if nil != err {
        return errors.WriteFileError{FileName: path, Err: err}
    }
    return nil
}

// Replaces template strings from a string give and provides a new string
func Replace(template string, values map[string]string) string {
    for match, replace := range values {
        template = strings.Replace(template, "{{"+match+"}}", replace, -1)
    }
    return template
}
