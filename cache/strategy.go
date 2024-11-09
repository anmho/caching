package cache

// Strategy is the selected caching algorithm.
type Strategy int

// TODO Describe the pros and cons here and in README.md
const (
	UnsetStrategy Strategy = iota
	WriteAround
	WriteThrough
	WriteBack
	ReadThrough
	CacheAside
)
