package pubsub

import (
	"time"
)

const (
	AggregateIDMessageAttributeKey      = "aggID"
	AggregateNameMessageAttributeKey    = "aggName"
	AggregateVersionMessageAttributeKey = "aggVer"
	EventIDMessageAttributeKey          = "id"
	EventNameMessageAttributeKey        = "name"
	EventTimeMessageAttributeKey        = "time"

	EventTimeFormat = time.RFC3339
)
