package io

import (
    "github.com/stretchr/testify/assert"
    "go-utils/fs"
    "math"
    "testing"
)

func TestWriteJsonProvideJsonExpectJsonWritten(t *testing.T) {
    a := assert.New(t)

    fs.SetFileSystem(fs.MemFs)

    type SampleJson struct {
        Name string
        Des  string
    }

    sampleJson := SampleJson{
        Name: "Waterloop",
        Des:  "Cool",
    }

    err := WriteJson("jsonFile.json", &sampleJson)
    if a.Nil(err) {
        a.True(fs.PathExists("jsonFile.json"))
    }

    text, err := fs.ReadFile("jsonFile.json")
    if err != nil {
        t.Fatal(err)
    }
    a.Equal(string(text), `{
  "Name": "Waterloop",
  "Des": "Cool"
}`)
}

func TestWriteJsonProvideInvalidJsonExpectError(t *testing.T) {
    a := assert.New(t)

    fs.SetFileSystem(fs.MemFs)

    type SampleJson struct {
        Name string
        Des  float64
    }

    err := WriteJson("jsonFile.json", &SampleJson{Des: math.NaN()})
    a.NotNil(err)

}

func TestWriteYamlProvideYamlExpectYamlWritten(t *testing.T) {
    a := assert.New(t)

    fs.SetFileSystem(fs.MemFs)

    type SampleJson struct {
        Name  string
        Des   string
        Slice []string          `yaml:"slice,omitempty,flow"`
        Map   map[string]string `yaml:"map,omitempty"`
    }

    sampleJson := SampleJson{
        Name:  "Waterloop",
        Des:   "Cool",
        Slice: []string{"Hello"},
        Map: map[string]string{
            "1": "One",
            "2": "Two",
        },
    }

    err := WriteYaml("yamlFile.yaml", &sampleJson)
    if a.Nil(err) {
        a.True(fs.PathExists("yamlFile.yaml"))
    }

    text, err := fs.ReadFile("yamlFile.yaml")
    if err != nil {
        t.Fatal(err)
    }
    a.Equal(string(text), `name: Waterloop
des: Cool
slice: [Hello]
map:
  "1": One
  "2": Two
`)
}

func TestWriteJsonProvideInvalidYamlExpectPanic(t *testing.T) {
    a := assert.New(t)

    fs.SetFileSystem(fs.MemFs)

    type SampleJson struct {
        Name string  `yaml:"name"`
        Des  float64 `yaml:"name"`
    }

    panicFunc := func() {
        WriteYaml("yamlFile.yaml", &SampleJson{Des: math.NaN()})
    }

    a.Panics(panicFunc)
}
