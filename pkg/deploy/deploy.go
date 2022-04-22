package deploy

const (
	// ResourceModified means an existing resource got changed
	ResourceModified = iota
	// ResourceAdded means a resource got added
	ResourceAdded
	// ResourceDeleted means a resource got deleted
	ResourceDeleted
)

const (
	// CookbookType denotes a cookbook resource
	CookbookType = iota
	// RoleType denotes a role resource
	RoleType
	// DataBagType denotes a data bag resource
	DataBagType
)

// Change represents a chef change
type Change struct {
	Name string
	Type int
	File string
}

// Deployer is an interface for deploying (or previewing) chef changes
type Deployer interface {
	Deploy([]Change) error
	Preview([]Change) error
}

// Differ is an interface to return chef diffs
type Differ interface {
	Diff() ([]Change, error)
}
