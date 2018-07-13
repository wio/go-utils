package assets

import (
    "path/filepath"
    "os"
    "go-utils/fs"
    "github.com/gobuffalo/packr"
    "go-utils/errors"
)

type boxStruct struct {
    Box packr.Box
    Initialized bool
}

var box = boxStruct {
    Initialized: false,
}

// Creates a pool of files that are static assets. When binary is compiled, these files will
// be included in that. Path provided is the path to the folder where assets are stored
func CreateStaticAssetsBox(packrBox packr.Box) {
    if box.Initialized {
        errors.String("static assets already created and initialized")
    } else {
        box = boxStruct{
            Box: packrBox,
            Initialized: true,
        }
    }
}

// Writes file stored in binary to actual file system.
// assetFilePath is a relative path to assets folder
func WriteFile(assetFilePath string, newPath string) error {
    bytes, err := box.Box.MustBytes(assetFilePath)
    if err != nil {
        return err
    } else {
        return fs.WriteFile(newPath, bytes)
    }
}

// Reads a file and returns bytes
// path is a relative path to assets folder
func ReadFile(path string) ([]byte, error) {
    return box.Box.MustBytes(path)
}

// This allows for copying asset files by creating an asset.json file. You can check the format of the file
// down below. You can provide constraints and extra info and this function will do everything
func CopyProjectAssets(structureData *StructureTypeData, constraintsProvided StructureConstraints,
    extra StructureExtraInfo) error {
    for _, path := range structureData.Paths {
        directoryPath := fs.Path(extra.ProjectDirectory, path.Entry)
        skipDir := false

        // handle directory constraints
        for _, constraint := range path.Constraints {
            value, exists := constraintsProvided.DirectoryConstraints[constraint]
            if exists && !value.Value {
                skipDir = true
                break
            }
        }
        if skipDir {
            continue
        }

        if !fs.PathExists(directoryPath) {
            if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
                return err
            }
        }

        for _, file := range path.Files {
            toPath := filepath.Clean(directoryPath + fs.Sep + file.To)
            skipFile := false
            // handle file constraints
            for _, constraint := range file.Constraints {
                value, exists := constraintsProvided.FileConstraints[constraint]
                if exists && !value.Value {
                    skipFile = true
                    break
                }
            }
            if skipFile {
                continue
            }

            // handle updates
            if !file.Update && extra.Update {
                continue
            }

            // copy assets
            if err := fs.CopyFile(fs.Path(extra.PlatformDirectory, file.From), toPath, file.Override); err != nil {
                return nil
            }
        }
    }

    return nil
}

// Sample asset.json file
/*
{
  "app": {
    "paths": [
      {
        "constraints": [],
        "entry": "/src",
        "files": [
          {
            "constraints": ["example", "cosa"],
            "from": "assets/example/cosa/app/main.cpp",
            "to": "main.cpp",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "arduino"],
            "from": "assets/example/arduino/app/main.cpp",
            "to": "main.cpp",
            "override": false,
            "update": false
          }
        ]
      }
    ]
  },
  "pkg": {
    "paths": [
      {
        "constraints": ["!header-only"],
        "entry": "/src",
        "files": [
          {
            "constraints": ["example", "cosa"],
            "from": "assets/example/cosa/pkg/output.cpp",
            "to": "output.cpp",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "arduino"],
            "from": "assets/example/arduino/pkg/output.cpp",
            "to": "output.cpp",
            "override": false,
            "update": false
          }
        ]
      },
      {
        "constraints": [],
        "entry": "/include",
        "files": [
          {
            "constraints": ["example", "!header-only", "cosa"],
            "from": "assets/example/cosa/pkg/output.h",
            "to": "output.h",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "!header-only", "arduino"],
            "from": "assets/example/arduino/pkg/output.h",
            "to": "output.h",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "header-only", "cosa"],
            "from": "assets/example/cosa/pkg-header-only/printer.h",
            "to": "printer.h",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "header-only", "arduino"],
            "from": "assets/example/arduino/pkg-header-only/printer.h",
            "to": "printer.h",
            "override": false,
            "update": false
          }
        ]
      },
      {
        "constraints": [],
        "entry": "/tests",
        "files": [
          {
            "constraints": ["example", "!header-only", "cosa"],
            "from": "assets/example/cosa/pkg/main.cpp",
            "to": "main.cpp",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "!header-only", "arduino"],
            "from": "assets/example/arduino/pkg/main.cpp",
            "to": "main.cpp",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "header-only", "cosa"],
            "from": "assets/example/cosa/pkg-header-only/main.cpp",
            "to": "main.cpp",
            "override": false,
            "update": false
          },
          {
            "constraints": ["example", "header-only", "arduino"],
            "from": "assets/example/arduino/pkg-header-only/main.cpp",
            "to": "main.cpp",
            "override": false,
            "update": false
          }
        ]
      }
    ]
  }
}
 */
