package ids

import (
	"bytes"

	"github.com/google/uuid"
	"github.com/totemcaf/gollections/slices"
)

// ID is a generic identifier for entities
type ID string

var nullId = ID(uuid.MustParse("00000000-0000-0000-0000-000000000000").String())

func New() ID {
	return ID(uuid.New().String())
}

func Empty() ID { return nullId }

func (id ID) IsEmpty() bool {
	return string(id) == "" || id == nullId
}

func (id ID) IsNotEmpty() bool {
	return !id.IsEmpty()
}

func (id ID) SelfOrNew() ID {
	if id.IsNotEmpty() {
		return id
	}
	return New()
}

func ParseId(id string) (ID, error) {
	if rawId, err := uuid.Parse(id); err == nil {
		return ID(rawId.String()), nil
	} else {
		return Empty(), err
	}
}

func FromBytes(bytes []byte) (ID, error) {
	if rawId, err := uuid.FromBytes(bytes); err == nil {
		return ID(rawId.String()), nil
	} else {
		return nullId, err
	}
}

func AllToString(ids []ID) []string {
	return slices.Map(ids, func(id ID) string { return id.String() })
}

func (id ID) String() string {
	return string(id)
}

func (id ID) ToBytes() []byte {
	rawId, _ := uuid.Parse(string(id))
	idBytes := [16]byte(rawId)

	return idBytes[:]
}

func (id ID) Compare(other ID) int {
	return bytes.Compare(id.ToBytes(), other.ToBytes())
}
