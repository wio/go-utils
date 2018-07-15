package fs

import (
    "github.com/spf13/afero"
    "go-utils/errors"
    "os"
    "path/filepath"
    "time"
)

var (
    OsFs  = afero.NewOsFs()
    MemFs = afero.NewMemMapFs()
)

var Sep = string(filepath.Separator)

type fileConfigStruct struct {
    FileSystem afero.Fs
}

// default file system configuration
var fileConfig = fileConfigStruct{
    FileSystem: afero.NewOsFs(),
}

// Allows changing of filesystem
func SetFileSystem(fs afero.Fs) {
    fileConfig.FileSystem = fs
}

// Chmod changes the mode of the named file to mode.
func Chmod(name string, mode os.FileMode) error {
    return fileConfig.FileSystem.Chmod(name, mode)
}

// Chtimes changes the access and modification times of the named file
func Chtimes(name string, atime time.Time, mtime time.Time) error {
    return fileConfig.FileSystem.Chtimes(name, atime, mtime)
}

// Create creates a file in the filesystem, returning the file and an
// error, if any happens.
func Create(name string) (afero.File, error) {
    return fileConfig.FileSystem.Create(name)
}

// Mkdir creates a directory in the filesystem, return an error if any
// happens.
func Mkdir(name string, perm os.FileMode) error {
    return fileConfig.FileSystem.Mkdir(name, perm)
}

// MkdirAll creates a directory path and all parents that does not exist
// yet.
func MkdirAll(path string, perm os.FileMode) error {
    return fileConfig.FileSystem.MkdirAll(path, perm)
}

// The name of this FileSystem
func Name() string {
    return fileConfig.FileSystem.Name()
}

// Open opens a file, returning it or an error, if anything happens.
func Open(name string) (afero.File, error) {
    return fileConfig.FileSystem.Open(name)
}

// OpenFile opens a file using the given flags and the given mode.
func OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
    return fileConfig.FileSystem.OpenFile(name, flag, perm)
}

// Remove removes a file identified by name, returning an error, if any
// happens.
func Remove(name string) error {
    return fileConfig.FileSystem.Remove(name)
}

// RemoveAll removes a directory path and any children it contains. It
// does not fail if the path does not exist (return nil).
func RemoveAll(path string) error {
    return fileConfig.FileSystem.RemoveAll(path)
}

// Rename renames a file.
func Rename(oldname, newname string) error {
    return fileConfig.FileSystem.Rename(oldname, newname)
}

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func Stat(name string) (os.FileInfo, error) {
    return fileConfig.FileSystem.Stat(name)
}

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
func Link(oldname, newname string) error {
    if fileConfig.FileSystem == MemFs {
        return &os.LinkError{Err: errors.String("link only available for OS filesystem")}
    } else {
        return os.Link(oldname, newname)
    }
}

// Symlink creates newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.
func Symlink(oldName string, newName string) error {
    if fileConfig.FileSystem == MemFs {
        return &os.LinkError{Err: errors.String("symlink only available for OS filesystem")}
    } else {
        return os.Symlink(oldName, newName)
    }
}
