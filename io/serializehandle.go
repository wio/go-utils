package io

import (
    "encoding/json"
    "go-utils/fs"
    "gopkg.in/yaml.v2"
)

// Parses JSON from the file on filesystem
func ParseJson(fileName string, out interface{}) (err error) {
    text, err := fs.ReadFile(fileName)
    if err != nil {
        return err
    }

    err = json.Unmarshal([]byte(text), out)
    return err
}

// Parses YML from the file on filesystem
func ParseYaml(fileName string, out interface{}) error {
    text, err := fs.ReadFile(fileName)
    if err != nil {
        return err
    }

    return yaml.Unmarshal(text, out)
}

// Writes JSON data to a file on filesystem
func WriteJson(fileName string, in interface{}) error {
    data, err := json.MarshalIndent(in, "", "  ")
    if err != nil {
        return err
    }

    return fs.WriteFile(fileName, data)
}

// Writes YML data to a file on filesystem
func WriteYaml(fileName string, in interface{}) error {
    data, err := yaml.Marshal(in)
    if err != nil {
        return err
    }

    return fs.WriteFile(fileName, data)
}
