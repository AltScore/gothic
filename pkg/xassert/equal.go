package xassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/AltScore/gothic/v2/pkg/ids"
	"github.com/AltScore/gothic/v2/pkg/xtime/date"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// EqualTime asserts the two times are same with milliseconds precision.
// Used to compare times serialized with bson
func EqualTime(t *testing.T, expected any, actual any) {
	assert.Equal(t, toTime(t, expected).UnixMilli(), toTime(t, actual).UnixMilli())
}

func toTime(t *testing.T, x any) time.Time {
	switch a := x.(type) {
	case bson.DateTime:
		return a.Time()
	case time.Time:
		return a
	case *time.Time:
		return *a
	case date.Date:
		return a.Time()
	case *date.Date:
		return a.Time()
	default:
		msg := fmt.Sprintf("expected bson.DateTime, got %T", x)
		assert.Fail(t, "", msg)
		return time.Time{}
	}
}

func EqualId(t *testing.T, expected ids.Id, actual any) {
	switch a := actual.(type) {
	case bson.Binary:
		bytes, err := uuid.FromBytes(a.Data)
		require.NoError(t, err)
		assert.Equal(t, expected, bytes)
	case string:
		parsed, err := uuid.Parse(a)
		require.NoError(t, err)
		assert.Equal(t, expected, parsed)
	default:
		assert.Fail(t, "expected bson.Binary, got %T", actual)
	}
}

// V returns the value of a field navigating a hierarchy of maps
func V(m any, keys ...string) any {
	for _, k := range keys {
		m = m.(bson.M)[k]
	}

	return m
}
