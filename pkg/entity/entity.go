package entity

import (
	"github.com/AltScore/gothic/pkg/ids"
)

type Entity interface {
	ID() ids.ID
	Version() uint
}
