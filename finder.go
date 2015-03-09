package finder

// Finder is an interface exposing methods to find files
// using multiple storage
type Finder interface {
	Find(path string)
	List(ignorePatterns []string)
}
