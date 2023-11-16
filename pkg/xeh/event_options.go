package xeh

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/xcontext"
	eh "github.com/looplab/eventhorizon"
)

// WithTenant adds the tenant to the event metadata
func WithTenant(ctx context.Context) eh.EventOption {
	return eh.WithMetadata(map[string]interface{}{
		"tenant": xcontext.GetTenantOrDefault(ctx),
	})
}

// GetEventTenant returns the tenant of the event, or the default tenant if not found
func GetEventTenant(event eh.Event) string {
	metadata := event.Metadata()
	if metadata == nil {
		return xcontext.DefaultTenant
	}

	tenant, ok := metadata["tenant"]
	if !ok {
		return xcontext.DefaultTenant
	}

	return tenant.(string)
}
