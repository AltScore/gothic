package ids

import (
	"bytes"
	"fmt"
	guuid "github.com/google/uuid"
	"github.com/looplab/eventhorizon/uuid"

	"github.com/totemcaf/gollections/slices"
)

// Id is a generic identifier for entities
type Id = uuid.UUID

var nullId = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func New() Id {
	return uuid.New()
}

func Empty() Id { return nullId }

func IsEmpty(id Id) bool {
	return id == nullId
}

func IsNotEmpty(id Id) bool {
	return !IsEmpty(id)
}

// SelfOrNew returns a new Id if id is empty, otherwise returns id
// Deprecated: use OrNew instead
func SelfOrNew(id Id) Id {
	if IsNotEmpty(id) {
		return id
	}
	return New()
}

// OrNew returns a new Id if id is empty, otherwise returns id
func OrNew(id Id) Id {
	if IsEmpty(id) {
		return New()
	}
	return id
}

// OrDefault returns the defaultId if id is empty, otherwise returns id
func OrDefault(id Id, defaultId Id) Id {
	if IsEmpty(id) {
		return defaultId
	}
	return id
}

// Parse parses a string into an Id.
// If the string is empty, returns an empty Id.
func Parse(id string) (Id, error) {
	if len(id) == 0 {
		return Empty(), nil
	}
	if parsedId, err := uuid.Parse(id); err == nil {
		return parsedId, nil
	} else {
		return Empty(), err
	}
}

// MustParse parses a string into an Id.
// If the string is empty, returns an empty Id.
// If the string is not a valid Id, panics.
func MustParse(id string) Id {
	if parsedId, err := Parse(id); err == nil {
		return parsedId
	} else {
		panic(err)
	}
}

// FromBytes converts a byte slice into an Id.
func FromBytes(bytes []byte) (Id, error) {
	if rawId, err := guuid.FromBytes(bytes); err == nil {
		return rawId, nil
	} else {
		return nullId, err
	}
}

// NewID computes a UUID string hashing the provided values
func NewID(values ...any) Id {
	str := fmt.Sprintf("%v", values)
	return guuid.NewSHA1(guuid.NameSpaceOID, []byte(str))
}

// AllToString converts a slice of Ids into a slice of strings
func AllToString(ids []Id) []string {
	return slices.Map(ids, func(id Id) string { return id.String() })
}

// ToBytes converts an Id into a byte slice
func ToBytes(id Id) []byte {
	idBytes := [16]byte(id)

	return idBytes[:]
}

// Compare compares two Ids
func Compare(id Id, other Id) int {
	return bytes.Compare(ToBytes(id), ToBytes(other))
}
