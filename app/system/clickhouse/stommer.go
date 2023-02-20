package clickhouse

import (
	"database/sql/driver"
	"fmt"
	"sort"

	"github.com/elgris/stom"
)

type Stommer struct {
	Mapper  map[string]interface{}
	Columns []string
	Values  []interface{}
}

func NewStommer(o interface{}, omitted ...string) (*Stommer, error) {
	m, err := stom.MustNewStom(o).ToMap(o)
	if err != nil {
		return nil, err
	}

	var (
		columns         []string
		columnsToValues = map[string]interface{}{}
	)
	for columnName, columnValue := range m {
		if isOmitted(columnName, omitted...) {
			continue
		}

		tmpVal := columnValue
		if valuer, ok := columnValue.(driver.Valuer); ok {
			tmpVal, err = valuer.Value()
			if err != nil {
				return nil, fmt.Errorf("could not convert value: %w", err)
			}
		}

		columnsToValues[columnName] = tmpVal
		columns = append(columns, columnName)
	}

	sort.Strings(columns)
	var values []interface{}
	for _, columnName := range columns {
		values = append(values, columnsToValues[columnName])
	}

	return &Stommer{
		Mapper:  m,
		Columns: columns,
		Values:  values,
	}, nil
}

func (s *Stommer) WithPrefix(prefix string) *Stommer {
	var newColumns []string
	for _, v := range s.Columns {
		newColumns = append(newColumns, fmt.Sprintf("%s.%s as %s", prefix, v, v))
	}

	s.Columns = newColumns

	return s
}

func isOmitted(key string, omitted ...string) bool {
	for _, omitted := range omitted {
		if key == omitted {
			return true
		}
	}

	return false
}
