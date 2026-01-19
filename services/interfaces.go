package services

type Lister interface {
	List(bool) error //true for json output
}

type ListerFor interface {
	ListFor(string, bool, bool) error //target, user, json
}

type Shower interface {
	Show(string, bool) error //name what should be displayed, true for json output
}
