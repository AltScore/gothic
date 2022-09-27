package es

type EntityID interface {
	Empty() bool
	Eq(id EntityID) bool
	String() string
}
