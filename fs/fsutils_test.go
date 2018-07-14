package fs

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
)

const (
    helloFile = "/home/hello.txt"
    emptyDirectory = "/home/emptyDirectory"
    oneFileDirectory = "/home/oneFileDirectory"
    allFilesDirectory = "/home"
    copyFilesDirectory = "/copy/home"
    invalidPath = "/jojojojojo/noway/noway"
)

func SetupFunction() {
    // for testing we will use memory filesystem
    SetFileSystem(MemFs)

    if _, err := Create(helloFile); err != nil {
        panic(err)
    }

    if err := Mkdir(emptyDirectory, os.ModePerm); err != nil {
        panic(err)
    }

    if err := Mkdir(oneFileDirectory, os.ModePerm); err != nil {
        panic(err)
    }

    if err := Mkdir(copyFilesDirectory, os.ModePerm); err != nil {
        panic(err)
    }

    if _, err := Create(oneFileDirectory + "/randomFile.txt"); err != nil {
        panic(err)
    }
}

func TearDownFunction() {
    if err := RemoveAll(allFilesDirectory); err != nil {
        panic(err)
    }

    if err := RemoveAll(copyFilesDirectory); err != nil {
        panic(err)
    }

    if err := RemoveAll(invalidPath); err != nil {
        panic(err)
    }
}

func TestMain(m *testing.M) {
    SetupFunction()
    retCode := m.Run()
    TearDownFunction()
    os.Exit(retCode)
}

func TestPathExistsProvideValidPathExpectTrue(t *testing.T) {
    a := assert.New(t)

    a.True(PathExists(helloFile), "created file path should exist")
}

func TestPathExistsProvideInValidPathExpectFalse(t *testing.T) {
    a := assert.New(t)

    a.False(PathExists(invalidPath), "this is a random file path and should not exist")
}

func TestIsDirProvideFileDoesNotExistExpectError(t *testing.T) {
    a := assert.New(t)

    _, err := IsDir(invalidPath)

    a.NotNil(err, "error thrown for path not existing")
}

func TestIsDirProvideDirectoryExpectTrue(t *testing.T) {
    a := assert.New(t)

    // true for directory
    dirValue, err := IsDir(emptyDirectory)
    if a.Nil(err, "no error should occur") {
        a.True(dirValue, "directory created should be recognized as a directory")
    }
}

func TestIsDirProvideFileExpectFalse(t *testing.T) {
    a := assert.New(t)

    // false for file
    fileValue, err := IsDir(helloFile)
    if a.Nil(err, "no error should occur") {
        a.False(fileValue, "file created should not be recognized as a directory")
    }
}

func TestIsEmptyProvideFileDoesNotExistExpectError(t *testing.T) {
    a := assert.New(t)

    dirPath := "./kjhashdjkhask"

    _, err := IsDirEmpty(dirPath)

    a.NotNil(err, "error thrown for path not existing")
}

func TestIsEmptyProvideEmptyFolderExpectTrue(t *testing.T) {
    a := assert.New(t)

    dirValue, err := IsDirEmpty(emptyDirectory)
    if a.Nil(err, "no error should occur") {
        a.True(dirValue, "created directory must be empty")
    }
}

func TestIsEmptyProvideFolderWithFilesExpectFalse(t *testing.T) {
    a := assert.New(t)

    dirValue, err := IsDirEmpty(allFilesDirectory)
    if a.Nil(err, "no error should occur") {
        a.False(dirValue, "created directory with a file must not be empty")
    }
}

func TestRemoveContentsProvideFolderWithFiles(t *testing.T) {
    a := assert.New(t)

    if err := Copy(allFilesDirectory, copyFilesDirectory, true); err != nil {
        t.Fatal(err)
    }

    err := RemoveContents(copyFilesDirectory)
    if a.Nil(err,"no error should occur") {
        isEmpty, err := IsDirEmpty(copyFilesDirectory)
        if err != nil {
            t.Fatal(err)
        }

        a.True(isEmpty, "since content has been deleted, there should not be any files")
        a.True(PathExists(copyFilesDirectory), "directory should not be deleted")
    }
}

func TestRemoveContentsProvideEmptyFolderExpectNoError(t *testing.T) {
    a := assert.New(t)

    err := RemoveContents(emptyDirectory)
    if a.Nil(err,"no error should occur") {
        isEmpty, err := IsDirEmpty(emptyDirectory)
        if err != nil {
            t.Fatal(err)
        }

        a.True(isEmpty, "since content has been deleted, there should not be any files")
        a.True(PathExists(emptyDirectory), "directory should not be deleted")
    }
}

func TestRemoveContentsProvideInvalidPathExpectError(t *testing.T) {
    a := assert.New(t)

    err := RemoveContents(invalidPath)
    a.NotNil(err, "path is invalid so an error must be thrown")
}

func TestCopyFileProvideValidFileExpectNoError(t *testing.T) {
    a := assert.New(t)

    fileName := "/helloFile.md"

    if PathExists(copyFilesDirectory + fileName) {
       if err := Remove(copyFilesDirectory + fileName); err != nil {
           t.Fatal(err)
       }
    }

    err := CopyFile(helloFile, copyFilesDirectory + fileName, false)
    defer func() {
        if err := Remove(copyFilesDirectory + fileName); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "no error should occur") {
        a.True(PathExists(copyFilesDirectory + fileName), "file should have been copied successfully")
    }
}

func TestCopyFileProvideInValidFileExpectError(t *testing.T) {
    a := assert.New(t)

    fileName := "/helloFile.md"

    err := CopyFile(invalidPath, copyFilesDirectory + fileName, false)

    a.NotNil(err, "invalid file cannot be copied")
}

func TestCopyFileProvideDirectoryExpectError(t *testing.T) {
    a := assert.New(t)

    fileName := "/helloFile.md"

    err := CopyFile(emptyDirectory, copyFilesDirectory + fileName, false)

    a.NotNil(err, "directory cannot be copied as a file")
}


func TestCopyFileProvideFileExistsNoOverrideExpectCopyNoError(t *testing.T) {
    a := assert.New(t)

    err := CopyFile(helloFile, helloFile, false)

    a.Nil(err, "destination exists so copy should have happened")
}

func TestCopyFileProvideFileOverrideExistsExpectCopyNoError(t *testing.T) {
    assert.New(t)
}
