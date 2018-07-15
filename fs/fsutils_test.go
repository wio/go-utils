package fs

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

const (
    helloFile          = "/home/hello.txt"
    emptyDirectory     = "/home/emptyDirectory"
    oneFileDirectory   = "/home/oneFileDirectory"
    allFilesDirectory  = "/home"
    copyFilesDirectory = "/copy/home"
    invalidPath        = "/jojojojojo/noway/noway"
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
    if a.Nil(err, "no error should occur") {
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
    if a.Nil(err, "no error should occur") {
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

func TestCopyDirProvideValidDirectoryExpectNoError(t *testing.T) {
    a := assert.New(t)

    dirName := "/helloDir"

    if PathExists(copyFilesDirectory + dirName) {
        if err := Remove(copyFilesDirectory + dirName); err != nil {
            t.Fatal(err)
        }
    }

    err := CopyDir(emptyDirectory, copyFilesDirectory+dirName, false)
    defer func() {
        if err := Remove(copyFilesDirectory + dirName); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "no error should occur") {
        a.True(PathExists(copyFilesDirectory+dirName), "directory should have been copied successfully")
    }
}

func TestCopyDirProvideInValidDirExpectError(t *testing.T) {
    a := assert.New(t)

    dirName := "/helloDir"

    err := CopyDir(invalidPath, copyFilesDirectory+dirName, false)

    a.NotNil(err, "invalid directory cannot be copied")
}

func TestCopyDirProvideFileExpectError(t *testing.T) {
    a := assert.New(t)

    dirName := "/helloDir"

    err := CopyDir(helloFile, copyFilesDirectory+dirName, false)

    a.NotNil(err, "file cannot be copied as a directory")
}

func TestCopyDirProvideDirExistsNoOverrideExpectCopyNoError(t *testing.T) {
    a := assert.New(t)

    dirName := "/helloDir"

    // create a dir with some files
    err := MkdirAll(dirName, os.ModePerm)
    if err != nil {
        t.Fatal(err)
    }

    _, err = Create(dirName + "/file1.txt")
    if err != nil {
        t.Fatal(err)
    }
    _, err = Create(dirName + "/file2.txt")
    if err != nil {
        t.Fatal(err)
    }

    status, err := IsDirEmpty(emptyDirectory)
    if err != nil {
        t.Fatal(err)
    }

    a.True(status, "empty directory must be empty")

    err = CopyDir(dirName, emptyDirectory, false)
    defer func() {
        if err := Remove(dirName); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "destination exists so copy should not have happened") {
        // confirm that empty directory is empty
        status, err := IsDirEmpty(emptyDirectory)
        if err != nil {
            t.Fatal(err)
        }

        a.True(status, "since empty directory is not replaced, empty directory must be empty")
    }
}

func TestCopyDirProvideDirOverrideExistsExpectCopyNoError(t *testing.T) {
    a := assert.New(t)

    dirName := "/helloDir"

    // create a dir with some files
    err := MkdirAll(dirName, os.ModePerm)
    if err != nil {
        t.Fatal(err)
    }

    _, err = Create(dirName + "/file1.txt")
    if err != nil {
        t.Fatal(err)
    }
    _, err = Create(dirName + "/file2.txt")
    if err != nil {
        t.Fatal(err)
    }

    status, err := IsDirEmpty(emptyDirectory)
    if err != nil {
        t.Fatal(err)
    }

    a.True(status, "empty directory must be empty")

    err = CopyDir(dirName, emptyDirectory, true)
    defer func() {
        if err := Remove(dirName); err != nil {
            t.Fatal(err)
        }

        if err := RemoveContents(emptyDirectory); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "destination exists but override so copy should have happened") {
        // confirm that empty directory is empty
        status, err := IsDirEmpty(emptyDirectory)
        if err != nil {
            t.Fatal(err)
        }

        a.False(status, "since empty directory is replaced, empty directory must not be empty")
    }
}

func TestCopyFileProvideValidFileExpectNoError(t *testing.T) {
    a := assert.New(t)

    fileName := "/helloFile.md"

    if PathExists(copyFilesDirectory + fileName) {
        if err := Remove(copyFilesDirectory + fileName); err != nil {
            t.Fatal(err)
        }
    }

    err := CopyFile(helloFile, copyFilesDirectory+fileName, false)
    defer func() {
        if err := Remove(copyFilesDirectory + fileName); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "no error should occur") {
        a.True(PathExists(copyFilesDirectory+fileName), "file should have been copied successfully")
    }
}

func TestCopyFileProvideInValidFileExpectError(t *testing.T) {
    a := assert.New(t)

    fileName := "/helloFile.md"

    err := CopyFile(invalidPath, copyFilesDirectory+fileName, false)

    a.NotNil(err, "invalid file cannot be copied")
}

func TestCopyFileProvideDirectoryExpectError(t *testing.T) {
    a := assert.New(t)

    fileName := "/helloFile.md"

    err := CopyFile(emptyDirectory, copyFilesDirectory+fileName, false)

    a.NotNil(err, "directory cannot be copied as a file")
}

func TestCopyFileProvideFileExistsNoOverrideExpectCopyNoError(t *testing.T) {
    a := assert.New(t)

    // create a file with some text
    err := WriteFile("somefile.txt", []byte("Hello World"))
    if err != nil {
        t.Fatal(err)
    }

    err = CopyFile("somefile.txt", helloFile, false)
    defer func() {
        if err := Remove("somefile.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "destination exists so copy should not have happened") {
        // confirm that content of hello file is still empty
        data, err := ReadFile(helloFile)
        if err != nil {
            t.Fatal(err)
        }

        a.Equal("", string(data), "since hello file is not modified, content is same")
    }
}

func TestCopyFileProvideFileOverrideExistsExpectCopyNoError(t *testing.T) {
    a := assert.New(t)

    // create a file with some text
    err := WriteFile("somefile.txt", []byte("Hello World"))
    if err != nil {
        t.Fatal(err)
    }

    err = CopyFile("somefile.txt", helloFile, true)
    defer func() {
        // make hello file empty again
        err := WriteFile(helloFile, []byte(""))
        if err != nil {
            t.Fatal(err)
        }

        if err := Remove("somefile.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "destination exists so copy should have happened") {
        // confirm that content of hello file have been set to hello world
        data, err := ReadFile(helloFile)
        if err != nil {
            t.Fatal(err)
        }

        a.Equal("Hello World", string(data), "since hello file is modified, content is changed")
    }
}

func TestCopy(t *testing.T) {
    a := assert.New(t)

    fileName := "helloFile.md"
    dirName := "sampleDir"

    // test copying file
    err := Copy(helloFile, fileName, false)
    defer func() {
        if err := Remove(fileName); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "no error should occur") {
        a.True(PathExists(fileName), "file should have been copied successfully")
    }

    // test copying directory
    Copy(allFilesDirectory, dirName, false)
    defer func() {
        if err := Remove(dirName); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "no error should occur") {
        a.True(PathExists(dirName), "dir should have been copied successfully")

        status, err := IsDirEmpty(emptyDirectory)
        if err != nil {
            t.Fatal(err)
        }

        a.True(status, "copied directory should have couple of files")
    }

    // create a file with some text
    err = WriteFile("somefile.txt", []byte("Hello World"))
    if err != nil {
        t.Fatal(err)
    }

    err = Copy("somefile.txt", helloFile, false)
    defer func() {
        if err := Remove("somefile.txt"); err != nil {
            t.Fatal(err)
        }
    }()

    if a.Nil(err, "destination exists so copy should not have happened") {
        // confirm that content of hello file is still empty
        data, err := ReadFile(helloFile)
        if err != nil {
            t.Fatal(err)
        }

        a.Equal("", string(data), "since hello file is not modified, content is same")
    }
}

func TestPath(t *testing.T) {
    a := assert.New(t)

    base1Path := "hello/deep"
    base1Extra1 := "jojo"
    base1Extra2 := "/lonzo"
    base1Extra3 := "lebronzo/"

    a.Equal("hello/deep/jojo", Path(base1Path, base1Extra1))
    a.Equal("hello/deep/lonzo", Path(base1Path, base1Extra2))
    a.Equal("hello/deep/lebronzo", Path(base1Path, base1Extra3))
    a.Equal("hello/deep/jojo/lebronzo/lonzo", Path(base1Path, base1Extra1, base1Extra3, base1Extra2))
}

func TestCopyMultipleFiles(t *testing.T) {
    a := assert.New(t)

    dirName1 := "testingDir1"
    dirName2 := "testingDir2"

    // basic copy
    err := CopyMultipleFiles([]string{
        helloFile, emptyDirectory, allFilesDirectory},
        []string{
            Path(dirName1, "helloFile.txt"),
            Path(dirName1, "emptyDir"),
            Path(dirName2, "filesDir")}, []bool{false, false, false})

    if a.Nil(err) {
        // check for helloFile.txt
        a.True(PathExists(Path(dirName1, "helloFile.txt")))
        a.True(PathExists(Path(dirName1, "emptyDir")))
        a.True(PathExists(Path(dirName1, "helloFile.txt")))
        a.True(PathExists(Path(dirName2, "filesDir")))
    }

    // fields not of the same size
    err = CopyMultipleFiles([]string{
        helloFile},
        []string{
            Path(dirName1, "helloFile.txt"),
            Path(dirName1, "emptyDir"),
            Path(dirName2, "filesDir")}, []bool{false, false})

    a.NotNil(err)
    a.EqualError(err, "length of sources, destinations and overrides is not equal")

    // error while copying (path does not exist)
    err = CopyMultipleFiles([]string{
        helloFile, emptyDirectory, invalidPath},
        []string{
            Path(dirName1, "helloFile.txt"),
            Path(dirName1, "emptyDir"),
            Path(dirName2, "filesDirRandom")}, []bool{false, false, false})

    a.NotNil(err)
    a.EqualErrorf(err, fmt.Sprintf("open %s: file does not exist", invalidPath), "")
}
