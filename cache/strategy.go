package cache


type Strategy int

const (
	UnsetStrategy Strategy = iota 
	WriteThrough
	WriteBack
	ReadThrough
	CacheAside
	WriteAround
)