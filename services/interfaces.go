package services

type Lister interface {
	// List now accepts jsonOutput (bool) and userPageSize (int)
	List(bool, int) error
}

type ListerFor interface {
	// ListFor now accepts target (string), user (bool), jsonOutput (bool), and userPageSize (int)
	ListFor(string, bool, bool, int) error
}

type Shower interface {
	// Show remains unchanged: name (string), jsonOutput (bool)
	Show(string, bool) error
}
