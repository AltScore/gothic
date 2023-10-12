package ids

import (
	"bytes"
	"fmt"

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

// SelfOrNew returns a new ID if id is empty, otherwise returns id
// Deprecated: use OrNew instead
func (id ID) SelfOrNew() ID {
	if id.IsNotEmpty() {
		return id
	}
	return New()
}

// OrNew returns a new ID if id is empty, otherwise returns id
func (id ID) OrNew() ID {
	if id.IsEmpty() {
		return New()
	}
	return id
}

// OrDefault returns the defaultId if id is empty, otherwise returns id
func (id ID) OrDefault(defaultId ID) ID {
	if id.IsEmpty() {
		return defaultId
	}
	return id
}

func (id ID) AsUUID() uuid.UUID {
	u, err := uuid.Parse(id.String())

	if err != nil {
		panic(err)
	}
	return u
}

func ParseID(id string) (ID, error) {
	if rawId, err := uuid.Parse(id); err == nil {
		return ID(rawId.String()), nil
	} else {
		return Empty(), err
	}
}

func FromUUID(uuid uuid.UUID) ID {
	return ID(uuid.String())
}

func FromBytes(bytes []byte) (ID, error) {
	if rawId, err := uuid.FromBytes(bytes); err == nil {
		return ID(rawId.String()), nil
	} else {
		return nullId, err
	}
}

// NewID computes a UUID string hashing the provided values
func NewID(values ...any) ID {
	str := fmt.Sprintf("%v", values)
	sha1Uuid := uuid.NewSHA1(uuid.NameSpaceOID, []byte(str))
	return FromUUID(sha1Uuid)
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
