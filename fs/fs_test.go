package fs

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestSetFileSystem(t *testing.T) {
    a := assert.New(t)

    // set file system to os fs
    SetFileSystem(OsFs)

    a.Equal(OsFs, fileConfig.FileSystem)

    // set file system to memory fs
    SetFileSystem(MemFs)
    a.Equal(MemFs, fileConfig.FileSystem)
}

func TestName(t *testing.T) {
    a := assert.New(t)

    // set file system to os fs
    SetFileSystem(OsFs)
    a.Equal("OsFs", Name())

    // set file system to memory fs
    SetFileSystem(MemFs)
    a.Equal("MemMapFS", Name())

}

func TestChmod(t *testing.T) {
    a := assert.New(t)

    file, err := OpenFile("randomFile.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
    defer func() {
        if err := file.Close(); err != nil {
            t.Fatal(err)
        }

        if err := Remove("randomFile.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if err != nil {
        t.Fatal(err)
    }

    statValue, err := file.Stat()
    if err != nil {
        t.Fatal(err)
    }

    a.Equal(statValue.Mode(), os.ModePerm)

    err = Chmod("randomFile.txt", os.ModeDevice)

    if a.Nil(err) {
        a.Equal(statValue.Mode(), os.ModeDevice)
    }
}

func TestRename(t *testing.T) {
    a := assert.New(t)

    file, err := Create("randomFile.txt")
    defer func() {
        if err := file.Close(); err != nil {
            t.Fatal(err)
        }

        if err := Remove("newFile.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if err != nil {
        t.Fatal(err)
    }

    err = Rename("randomFile.txt", "newFile.txt")

    if a.Nil(err) {
        a.True(PathExists("newFile.txt"))
    }
}

func TestSymlink(t *testing.T) {
    a := assert.New(t)

    /// Operating system File system
    SetFileSystem(OsFs)

    _, err := Create("fileOs.txt")
    defer func() {
        SetFileSystem(OsFs)

        if err := Remove("fileOsSym.txt"); err != nil {
            t.Fatal(err)
        }

        if err := Remove("fileOs.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if err != nil {
        t.Fatal(err)
    }

    err = Symlink("fileOs.txt", "fileOsSym.txt")

    if a.Nil(err) {
        a.True(PathExists("fileOsSym.txt"))
    }

    /// Memory File system
    SetFileSystem(MemFs)

    _, err = Create("fileMem.txt")
    defer func() {
        SetFileSystem(MemFs)

        if err := Remove("fileMem.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if err != nil {
        t.Fatal(err)
    }

    err = Symlink("fileMem.txt", "fileMemSym.txt")

    a.NotNil(err)
}

func TestLink(t *testing.T) {
    defer func() {
        SetFileSystem(MemFs)
    }()

    a := assert.New(t)

    /// Operating system File system
    SetFileSystem(OsFs)

    _, err := Create("fileOs.txt")
    defer func() {
        SetFileSystem(OsFs)

        if err := Remove("fileOsLink.txt"); err != nil {
            t.Fatal(err)
        }

        if err := Remove("fileOs.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if err != nil {
        t.Fatal(err)
    }

    err = Link("fileOs.txt", "fileOsLink.txt")

    if a.Nil(err) {
        a.True(PathExists("fileOsLink.txt"))
    }

    /// Memory File system
    SetFileSystem(MemFs)

    _, err = Create("fileMem.txt")
    defer func() {
        SetFileSystem(MemFs)

        if err := Remove("fileMem.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if err != nil {
        t.Fatal(err)
    }

    err = Link("fileMem.txt", "fileMemLink.txt")

    a.NotNil(err)
}
