package symbol

// Build ... analyze from server.txt
//
type Build struct {
	ID      string
	Date    string
	Branch  string
	Version string
	Comment string
}

// Symbol represent each symbol file's detail
//
type Symbol struct {
	Arch    string // x64 or x86
	Hash    string
	Name    string
	Path    string
	Version string
}

// Builder interface
//
type Builder interface {
	// Name return builder name
	Name() string

	// Init builder
	Init() error

	// Persist builder information
	Persist() error

	// Add an given build pdb to symbol server.
	// if `buildVersion` is empty, it will try to add the latest build on build server if exist.
	Add(buildVerion string) error

	// GetSymbolPath get given symbol file full path on symbol server.
	// The path can be used to serve download.
	GetSymbolPath(hash, name string) string

	// ParseBuilds parse all version of builds that already in the symbol server of curent branch.
	//
	ParseBuilds(handler func(b *Build) error) (int, error)

	// ParseSymbols parse all the symbols of given build vesrion
	//
	ParseSymbols(buildID string, handler func(sym *Symbol) error) (int, error)

	// SetSubpath change the subpath on build server and local store.
	SetSubpath(buildserver, localstore string) error

	// CanUpdate check if current branch is valid on build server.
	CanUpdate() bool
	// CanBrowse check if current branch is valid on local symbol store.
	CanBrowse() bool
}
