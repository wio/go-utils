package errors

import "fmt"
import "errors"

const (
    Spaces = " "
)

type Error interface {
    error
}

type Generic struct {
    message string
}

func (err Generic) Error() string {
    return err.message
}

func String(message string) error {
    return errors.New(message)
}

func Stringf(format string, a ...interface{}) error {
    msg := fmt.Sprintf(format, a...)
    return String(msg)
}

type ReadFileError struct {
    FileName string
    Err      error
}

func (err ReadFileError) Error() string {
    str := fmt.Sprintf(`"%s" file read failed`, err.FileName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type WriteFileError struct {
    FileName string
    Err      error
}

func (err WriteFileError) Error() string {
    str := fmt.Sprintf(`"%s" file write failed`, err.FileName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type YamlMarshallError struct {
    Err error
}

func (err YamlMarshallError) Error() string {
    str := fmt.Sprintf(`yaml data could not be marshalled`)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type PathDoesNotExist struct {
    Path string
    Err  error
}

func (err PathDoesNotExist) Error() string {
    str := fmt.Sprintf(`path does not exist: %s`, err.Path)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type DeleteDirectoryError struct {
    DirName string
    Err     error
}

func (err DeleteDirectoryError) Error() string {
    str := fmt.Sprintf(`"%s" directory failed to be deleted`, err.DirName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type DeleteFileError struct {
    FileName string
    Err      error
}

func (err DeleteFileError) Error() string {
    str := fmt.Sprintf(`"%s" file failed to be deleted`, err.FileName)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}

type FatalError struct {
    Log interface{}
    Err error
}

func (err FatalError) Error() string {
    str := fmt.Sprintf("a fatal error occured. Contact developers for a fix")
    str += "\n" + Spaces + fmt.Sprintln(err.Log)

    if err.Err != nil {
        str += fmt.Sprintf("\n%s%s", Spaces, err.Err.Error())
    }

    return str
}
