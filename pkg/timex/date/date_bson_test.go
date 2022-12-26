package date

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDate_can_marshal_and_unmarshal(t *testing.T) {
	d := New(1963, 11, 29)

	bytes, err := d.MarshalJSON()

	require.NoError(t, err)

	var d2 Date

	err = d2.UnmarshalJSON(bytes)

	require.NoError(t, err)

	require.Equal(t, d, d2)
}

type sampleStructWithDate struct {
	D Date `bson:"d"`
}

func TestDate_can_marshal_and_unmarshal_in_struct(t *testing.T) {
	d := New(1963, 11, 29)

	s := sampleStructWithDate{d}

	bytes, err := bson.Marshal(s)

	require.NoError(t, err)

	var s2 sampleStructWithDate

	err = bson.Unmarshal(bytes, &s2)

	require.NoError(t, err)

	require.Equal(t, s, s2)
}
