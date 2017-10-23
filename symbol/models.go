package symbol

// Branch ... information
//
type Branch struct {
	BuildName   string `json:"buildName"`
	StoreName   string `json:"storeName"`
	BuildPath   string `json:"buildPath"`
	StorePath   string `json:"storePath"`
	CreateDate  string `json:"createDate"`
	LatestBuild string `json:"latestBuild"`
}

// Build ... analyze from server.txt
//
type Build struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Branch  string `json:"branch"`
	Version string `json:"version"`
	Comment string `json:"comment"`
}

// Symbol represent each symbol file's detail
//
type Symbol struct {
	Arch    string `json:"arch"` // x64 or x86
	Hash    string `json:"hash"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

// Builder interface
//
type Builder interface {
	// Name return builder name
	Name() string
	// Init builder
	Init() error
	// Delete current branch
	Delete() error
	// Persist builder information
	Persist() error
	// CanUpdate check if current branch is valid on build server.
	CanUpdate() bool
	// CanBrowse check if current branch is valid on local symbol store.
	CanBrowse() bool

	// SetSubpath change the subpath on build server and local store.
	SetSubpath(buildserver, localstore string) error

	// Add an given build pdb to symbol server.
	// if `buildVersion` is empty, it will try to add the latest build on build server if exist.
	AddBuild(buildVerion string) error

	// GetSymbolPath get given symbol file full path on symbol server.
	// The path can be used to serve download.
	GetSymbolPath(hash, name string) string

	// ParseBuilds parse all version of builds that already in the symbol server of curent branch.
	//
	ParseBuilds(handler func(b *Build) error) (int, error)

	// ParseSymbols parse all the symbols of given build vesrion
	//
	ParseSymbols(buildID string, handler func(sym *Symbol) error) (int, error)
}
