package fs

import (
    "go-utils/errors"
    "path/filepath"

    "bytes"
    "io"
    "os"
    "github.com/spf13/afero"
)

// Checks if the give path is a director and based on the returns
// true or false. If path does not exist, it throws an error
func IsDir(path string) (bool, error) {
    fi, err := Stat(path)
    if err != nil {
        return false, err
    }

    return fi.IsDir(), nil
}

// This checks if the directory is empty or not
func IsDirEmpty(name string) (bool, error) {
    f, err := Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string, override bool) error {
    if PathExists(dst) && !override {
        return nil
    }

    // check directory and throw and error if it is given
    status, err := IsDir(src)
    if err != nil {
        return err
    } else if status {
        return errors.Stringf("Src Path [%s] cannot be a directory", src)
    }

    if !PathExists(src) {
        return errors.Stringf("Path [%s] does not exist", src)
    }

    in, err := Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        if e := out.Close(); e != nil {
            err = e
        }
    }()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }

    err = out.Sync()
    if err != nil {
        return err
    }

    si, err := Stat(src)
    if err != nil {
        return err
    }
    err = Chmod(dst, si.Mode())
    if err != nil {
        return err
    }

    return nil
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string, override bool) (err error) {
    if PathExists(dst) && !override {
        return
    } else {
        if err := RemoveAll(dst); err != nil {
            return err
        }
    }

    if !PathExists(src) {
        return
    }

    src = filepath.Clean(src)
    dst = filepath.Clean(dst)

    si, err := Stat(src)
    if err != nil {
        return err
    }
    if !si.IsDir() {
        return errors.String("source is not a directory")
    }

    _, err = Stat(dst)
    if err != nil && !os.IsNotExist(err) {
        return
    }

    err = MkdirAll(dst, si.Mode())
    if err != nil {
        return
    }

    entries, err := afero.ReadDir(fileConfig.FileSystem, src)
    if err != nil {
        return
    }

    for _, entry := range entries {
        srcPath := filepath.Join(src, entry.Name())
        dstPath := filepath.Join(dst, entry.Name())

        if entry.IsDir() {
            err = CopyDir(srcPath, dstPath, override)
            if err != nil {
                return
            }
        } else {
            // Skip symlinks.
            if entry.Mode()&os.ModeSymlink != 0 {
                continue
            }

            err = CopyFile(srcPath, dstPath, override)
            if err != nil {
                return
            }
        }
    }

    return
}

// Generic copy function that can copy anything from src to destination
func Copy(src string, dst string, override bool) error {
    if PathExists(dst) && !override {
        return nil
    }

    src = filepath.Clean(src)
    dst = filepath.Clean(dst)

    si, err := Stat(src)
    if err != nil {
        return err
    }
    if si.IsDir() {
        return CopyDir(src, dst, override)
    } else {
        return CopyFile(src, dst, override)
    }
}

// Copies multiple files from source to destination. Source files are from filesystem
func CopyMultipleFiles(sources []string, destinations []string, overrides []bool) error {
    if len(sources) != len(destinations) || len(destinations) != len(overrides) {
        return errors.String("length of sources, destinations and overrides is not equal")
    }

    for i := 0; i < len(sources); i++ {
        if err := Copy(sources[i], destinations[i], overrides[i]); err != nil {
            return err
        }
    }

    return nil
}

// Reads the file and provides it's content as a string. From normal filesystem
func ReadFile(fileName string) ([]byte, error) {
    buff, err := afero.ReadFile(fileConfig.FileSystem, fileName)
    return buff, err
}

// Writes text to a file on normal filesystem
func WriteFile(fileName string, data []byte) error {
    return afero.WriteFile(fileConfig.FileSystem, fileName, data, os.ModePerm)
}

// Joins multiple paths together and provides a native path
func Path(values ...string) string {
    var buffer bytes.Buffer
    for _, value := range values {
        buffer.WriteString(value)
        buffer.WriteString(Sep)
    }
    path := buffer.String()
    return filepath.Clean(path[:len(path)-1])
}

// Checks if path exists and returns true and false based on that
func PathExists(path string) bool {
    _, err := Stat(path)
    if os.IsNotExist(err) || err != nil {
        return false
    }

    return true
}

// Deletes all the files from the directory
func RemoveContents(dir string) error {
    d, err := Open(dir)
    if err != nil {
        return err
    }

    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }

    for _, name := range names {
        if err = RemoveAll(filepath.Join(dir, name)); err != nil {
            return err
        }
    }
    return nil
}
