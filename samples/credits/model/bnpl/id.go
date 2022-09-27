package bnpl

import (
	"fmt"
	"time"
)

type ID string

func (id ID) New() ID {
	return ID(fmt.Sprintf("flow-%d", time.Now().UnixNano()))
}

func (id ID) Empty() bool {
	return id == ""
}

func (id ID) Eq(id2 ID) bool {
	return id == id2
}

func (id ID) String() string {
	return string(id)
}
