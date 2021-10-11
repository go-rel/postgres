package postgres

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

// Quote PostgreSQL identifiers and literals.
type Quote struct{}

func (q Quote) ID(name string) string {
	return pq.QuoteIdentifier(name)
}

func (q Quote) Value(v interface{}) string {
	switch v := v.(type) {
	default:
		panic("unsupported value")
	case string:
		return pq.QuoteLiteral(v)
	}
}

// ValueConvert converts values to PostgreSQL literals.
type ValueConvert struct{}

func (c ValueConvert) ConvertValue(v interface{}) (driver.Value, error) {
	v, err := driver.DefaultParameterConverter.ConvertValue(v)
	if err != nil {
		return nil, err
	}
	switch v := v.(type) {
	default:
		return v, nil
	case time.Time:
		return v.Truncate(time.Microsecond).Format("2006-01-02 15:04:05.999999999Z07:00:00"), nil
	}
}
