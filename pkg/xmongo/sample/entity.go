package sample

type Person interface {
	Id() string
	Name() string
	Age() int
}

// let compiler check if person implements interface Person
var _ Person = (*person)(nil)

type person struct {
	id   string
	name string
	age  int
}

func (p person) Id() string {
	return p.id
}

func (p person) Name() string {
	return p.name
}

func (p person) Age() int {
	return p.age
}
