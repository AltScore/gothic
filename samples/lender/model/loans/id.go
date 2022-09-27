package loans

import (
	"fmt"
	"time"

	"github.com/AltScore/gothic/pkg/es"
)

type ID string

func NewId() ID {
	return ID(fmt.Sprintf("flow-%d", time.Now().UnixNano()))
}

func (id ID) Empty() bool {
	return id == ""
}

func (id ID) Eq(id2 es.EntityID) bool {
	return id == id2
}

func (id ID) String() string {
	return string(id)
}
