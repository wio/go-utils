package template

import (
    "github.com/valyala/fasttemplate"
    "go-utils/errors"
    "go-utils/fs"
    "io"
)

// Converts normal function for templates into Tag function used by fasttemplate
func TagFunc(function func(io.Writer, string) (int, error)) fasttemplate.TagFunc {
    return fasttemplate.TagFunc(function)
}

// Reads a file, replaces template strings with values provided, and writes the file back with new changes
func IOReplace(path string, start, end string, values map[string]interface{}) error {
    data, err := fs.ReadFile(path)
    if nil != err {
        return errors.ReadFileError{FileName: path, Err: err}
    }

    template := string(data)
    t := fasttemplate.New(template, start, end)

    result := t.ExecuteString(values)
    err = fs.WriteFile(path, []byte(result))
    if nil != err {
        return errors.WriteFileError{FileName: path, Err: err}
    }
    return nil
}

// Replaces template strings from a string give and provides a new string
func Replace(template, start, end string, values map[string]interface{}) string {
    t := fasttemplate.New(template, start, end)

    return t.ExecuteString(values)
}
