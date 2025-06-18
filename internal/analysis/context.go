package analysis

type ProjectContext struct {
	Files []FileInfo
	Deps  []Dependency
}

func ExtractContext(files []FileInfo, deps []Dependency) *ProjectContext {
	return &ProjectContext{
		Files: files,
		Deps:  deps,
	}
}
