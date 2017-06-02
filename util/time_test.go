package util_test

import (
	"testing"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func TestFormatTime(t *testing.T) {
	var (
		timestamp        time.Time
		expected, actual string
	)

	// UTC
	timestamp = time.Date(2012, 12, 11, 8, 52, 31, 493729031, time.UTC)
	expected = "2012-12-11T08:52:31.493Z"
	actual = util.FormatTime(&timestamp)
	assert.Equal(t, expected, actual)

	// UTC
	timestamp = time.Date(2012, 12, 11, 8, 52, 31, 493729031, time.FixedZone("HKT", 8*3600))
	expected = "2012-12-11T00:52:31.493Z"
	actual = util.FormatTime(&timestamp)
	assert.Equal(t, expected, actual)
}

func TestParseTimestamp(t *testing.T) {
	var (
		parsedTime *time.Time
		err        error
	)

	parsedTime, err = util.ParseTimestamp("")
	assert.Nil(t, err)
	assert.Nil(t, parsedTime)

	parsedTime, err = util.ParseTimestamp("bogus")
	assert.NotNil(t, err)
	assert.Nil(t, parsedTime)

	parsedTime, err = util.ParseTimestamp("2016-05-04T12:08:35Z")
	assert.Nil(t, err)
	assert.NotNil(t, parsedTime)
	assert.Equal(t, 2016, parsedTime.UTC().Year())
	assert.Equal(t, time.May, parsedTime.UTC().Month())
	assert.Equal(t, 4, parsedTime.UTC().Day())
	assert.Equal(t, 12, parsedTime.UTC().Hour())
	assert.Equal(t, 8, parsedTime.UTC().Minute())
	assert.Equal(t, 35, parsedTime.UTC().Second())

	parsedTime, err = util.ParseTimestamp("2016-05-04T12:08:35+07:00")
	assert.Nil(t, err)
	assert.NotNil(t, parsedTime)
	assert.Equal(t, 2016, parsedTime.UTC().Year())
	assert.Equal(t, time.May, parsedTime.UTC().Month())
	assert.Equal(t, 4, parsedTime.UTC().Day())
	assert.Equal(t, 5, parsedTime.UTC().Hour())
	assert.Equal(t, 8, parsedTime.UTC().Minute())
	assert.Equal(t, 35, parsedTime.UTC().Second())

	// RFC3339 timestamp generated in Django: '{}{}'.format(timezone.now().utcnow(), 'Z')
	parsedTime, err = util.ParseTimestamp("2016-07-19T09:28:03.285195Z")
	assert.Nil(t, err)
	assert.NotNil(t, parsedTime)
	assert.Equal(t, 2016, parsedTime.UTC().Year())
	assert.Equal(t, time.July, parsedTime.UTC().Month())
	assert.Equal(t, 19, parsedTime.UTC().Day())
	assert.Equal(t, 9, parsedTime.UTC().Hour())
	assert.Equal(t, 28, parsedTime.UTC().Minute())
	assert.Equal(t, 3, parsedTime.UTC().Second())
}

func TestFormatDate(t *testing.T) {
	var (
		timestamp        time.Time
		expected, actual string
	)

	// UTC
	timestamp = time.Date(2012, 12, 11, 8, 52, 31, 493729031, time.UTC)
	expected = "2012-12-11"
	actual = util.FormatDate(&timestamp)
	assert.Equal(t, expected, actual)
}

func TestParseDate(t *testing.T) {
	var (
		parsedTime *time.Time
		err        error
	)

	parsedTime, err = util.ParseDate("")
	assert.Nil(t, err)
	assert.Nil(t, parsedTime)

	parsedTime, err = util.ParseDate("bogus")
	assert.NotNil(t, err)
	assert.Nil(t, parsedTime)

	parsedTime, err = util.ParseDate("2016-05-04")
	assert.Nil(t, err)
	assert.NotNil(t, parsedTime)
	assert.Equal(t, 2016, parsedTime.UTC().Year())
	assert.Equal(t, time.May, parsedTime.UTC().Month())
	assert.Equal(t, 4, parsedTime.UTC().Day())
	assert.Equal(t, 0, parsedTime.UTC().Hour())
	assert.Equal(t, 0, parsedTime.UTC().Minute())
	assert.Equal(t, 0, parsedTime.UTC().Second())
}
