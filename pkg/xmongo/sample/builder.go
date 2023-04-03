package sample

type Builder struct {
	person person
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithId(id string) *Builder {
	b.person.id = id
	return b
}

func (b *Builder) WithName(name string) *Builder {
	b.person.name = name
	return b

}

func (b *Builder) WithAge(age int) *Builder {
	b.person.age = age
	return b
}

func (b *Builder) Build() Person {
	return b.person
}
