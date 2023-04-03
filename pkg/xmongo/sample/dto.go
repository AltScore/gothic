package sample

type personDto struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func (p personDto) ToEntity() Person {
	return person{
		id:   p.Id,
		name: p.Name,
		age:  p.Age,
	}
}

func fromEntity(entity Person) personDto {
	return personDto{
		Id:   entity.Id(),
		Name: entity.Name(),
		Age:  entity.Age(),
	}
}
