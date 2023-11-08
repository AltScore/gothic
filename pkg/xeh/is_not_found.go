package xeh

import (
	"errors"

	eh "github.com/looplab/eventhorizon"
)

func IsEHNotFound(err error) bool { return errors.Is(err, eh.ErrEntityNotFound) }
