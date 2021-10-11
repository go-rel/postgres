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
