package ehmocks

import (
	"context"
	"reflect"

	eh "github.com/looplab/eventhorizon"
	"github.com/stretchr/testify/mock"
)

// CommandHandlerMock is a mock implementation of eh.CommandHandler
type CommandHandlerMock struct {
	mock.Mock
}

func (c *CommandHandlerMock) HandleCommand(ctx context.Context, command eh.Command) error {
	value := reflect.ValueOf(command)

	if value.Kind() != reflect.Ptr {
		panic("command must be a pointer, check code to see if it is a pointer to a struct")
	}

	args := c.Called(ctx, command)
	return args.Error(0)
}
