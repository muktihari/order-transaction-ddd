package postgresql

import (
	"database/sql"
	"encoding/json"
)

// Map maps current cursor of sqlRows into v
func Map(sqlRows *sql.Rows, v interface{}) error {
	columnNames, err := sqlRows.Columns()
	if err != nil {
		return err
	}

	var n = len(columnNames)
	var buff = make([]interface{}, n)
	var buffPtr = make([]interface{}, n)
	for i := range columnNames {
		buffPtr[i] = &buff[i]
	}

	if err := sqlRows.Scan(buffPtr...); err != nil {
		return err
	}

	m := map[string]interface{}{}
	for i, columnName := range columnNames {
		m[columnName] = buff[i]
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
