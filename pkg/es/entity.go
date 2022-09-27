package es

type EntityID[ID EntityID[ID]] interface {
	New() ID
	Empty() bool
	Eq(ID) bool
	String() string
}
