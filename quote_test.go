package postgres

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuote_Panic(t *testing.T) {
	quoter := Quote{}
	assert.PanicsWithValue(t, "unsupported value", func() {
		quoter.Value(1)
	})
}

func TestQuote_ID(t *testing.T) {
	quoter := Quote{}

	cases := []struct {
		input string
		want  string
	}{
		{`foo`, `"foo"`},
		{`foo bar baz`, `"foo bar baz"`},
		{`foo"bar`, `"foo""bar"`},
		{"foo\x00bar", `"foo"`},
		{"\x00foo", `""`},
	}

	for _, test := range cases {
		assert.Equal(t, test.want, quoter.ID(test.input))
	}
}

func TestQuote_Value(t *testing.T) {
	quoter := Quote{}

	cases := []struct {
		input string
		want  string
	}{
		{`foo`, `'foo'`},
		{`foo bar baz`, `'foo bar baz'`},
		{`foo'bar`, `'foo''bar'`},
		{`foo\bar`, ` E'foo\\bar'`},
		{`foo\ba'r`, ` E'foo\\ba''r'`},
		{`foo"bar`, `'foo"bar'`},
		{`foo\x00bar`, ` E'foo\\x00bar'`},
		{`\x00foo`, ` E'\\x00foo'`},
		{`'`, `''''`},
		{`''`, `''''''`},
		{`\`, ` E'\\'`},
		{`'abc'; DROP TABLE users;`, `'''abc''; DROP TABLE users;'`},
		{`\'`, ` E'\\'''`},
		{`E'\''`, ` E'E''\\'''''`},
		{`e'\''`, ` E'e''\\'''''`},
		{`E'\'abc\'; DROP TABLE users;'`, ` E'E''\\''abc\\''; DROP TABLE users;'''`},
		{`e'\'abc\'; DROP TABLE users;'`, ` E'e''\\''abc\\''; DROP TABLE users;'''`},
	}

	for _, test := range cases {
		assert.Equal(t, test.want, quoter.Value(test.input))
	}
}

type customType int

func (c customType) Value() (driver.Value, error) {
	return int(c), nil
}

func TestValueConvert_CustomType(t *testing.T) {
	valuer := ValueConvert{}
	v, err := valuer.ConvertValue(customType(1))
	assert.EqualError(t, err, "non-Value type int returned from Value")
	assert.Nil(t, v)
}
