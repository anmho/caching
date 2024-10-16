package cache

// Strategy is the selected caching algorithm.
type Strategy int

const (
	UnsetStrategy Strategy = iota
	WriteThrough
	WriteBack
	ReadThrough
	CacheAside
	WriteAround
)
