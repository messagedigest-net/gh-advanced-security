package services

type Lister interface {
	// List now accepts jsonOutput (bool) , userPageSize (int), and fetchAll (bool)
	List(bool, int, bool) error
}

type ListerFor interface {
	// ListFor now accepts target (string), user (bool), jsonOutput (bool),  userPageSize (int), and fetchAll (bool)
	ListFor(string, bool, bool, int, bool) error
}

type Shower interface {
	// Show remains unchanged: name (string), jsonOutput (bool)
	Show(string, bool) error
}
