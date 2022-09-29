package loans

import (
	"fmt"
	"time"
)

type ID = string

func NewId() ID {
	return fmt.Sprintf("flow-%d", time.Now().UnixNano())
}
