package assets

// ############################################ projectType for asset.json #####################################
type StructureFilesData struct {
    Constraints []string
    From        string
    To          string
    Override    bool
    Update      bool
}

type StructurePathData struct {
    Constraints []string
    Entry       string
    Files       []StructureFilesData
}

type StructureTypeData struct {
    Paths []StructurePathData
}

// Types of data: app level, pkg level and all level
type StructureConfigData struct {
    App StructureTypeData
    Pkg StructureTypeData
    All StructureTypeData
}

// ##################################### Constraints that can be applied to asset.json #########################
type StructureConstraint struct {
    Value bool
}

type StructureConstraints struct {
    DirectoryConstraints map[string]StructureConstraint
    FileConstraints map[string]StructureConstraint
}

// ##################################### Extra information needed by asset.json file ###########################
type StructureExtraInfo struct {
    ProjectDirectory string
    PlatformDirectory string
    Update bool
}
