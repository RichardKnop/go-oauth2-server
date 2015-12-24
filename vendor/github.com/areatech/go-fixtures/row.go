package fixtures

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	onInsertNow    = "ON_INSERT_NOW()"
	onUpdateNow    = "ON_UPDATE_NOW()"
	postgresDriver = "postgres"
)

// Row represents a single database row
type Row struct {
	Table              string
	PK                 map[string]interface{}
	Fields             map[string]interface{}
	insertColumnLength int
	updateColumnLength int
	pkColumns          []string
	pkValues           []interface{}
	insertColumns      []string
	updateColumns      []string
	insertValues       []interface{}
	updateValues       []interface{}
}

// Init loads internal struct variables
func (row *Row) Init() {
	// Initial values
	row.insertColumnLength = len(row.PK) + len(row.Fields)
	row.updateColumnLength = len(row.PK) + len(row.Fields)
	row.pkColumns = make([]string, 0)
	row.pkValues = make([]interface{}, 0)
	row.insertColumns = make([]string, 0)
	row.updateColumns = make([]string, 0)
	row.insertValues = make([]interface{}, 0)
	row.updateValues = make([]interface{}, 0)

	// Get and sort map keys
	var i int
	pkKeys := make([]string, len(row.PK))
	i = 0
	for pkKey := range row.PK {
		pkKeys[i] = pkKey
		i++
	}
	sort.Strings(pkKeys)
	fieldKeys := make([]string, len(row.Fields))
	i = 0
	for fieldKey := range row.Fields {
		fieldKeys[i] = fieldKey
		i++
	}
	sort.Strings(fieldKeys)

	// Primary keys
	for _, pkKey := range pkKeys {
		row.pkColumns = append(row.pkColumns, pkKey)
		row.pkValues = append(row.pkValues, row.PK[pkKey])
		row.insertColumns = append(row.insertColumns, pkKey)
		row.updateColumns = append(row.updateColumns, pkKey)
		row.insertValues = append(row.insertValues, row.PK[pkKey])
		row.updateValues = append(row.updateValues, row.PK[pkKey])
	}

	// Rest of the fields
	for _, fieldKey := range fieldKeys {
		sv, ok := row.Fields[fieldKey].(string)
		if ok && sv == onInsertNow {
			row.insertColumns = append(row.insertColumns, fieldKey)
			row.insertValues = append(row.insertValues, time.Now())
			row.updateColumnLength--
			continue
		}
		if ok && sv == onUpdateNow {
			row.updateColumns = append(row.updateColumns, fieldKey)
			row.updateValues = append(row.updateValues, time.Now())
			row.insertColumnLength--
			continue
		}
		row.insertColumns = append(row.insertColumns, fieldKey)
		row.updateColumns = append(row.updateColumns, fieldKey)
		row.insertValues = append(row.insertValues, row.Fields[fieldKey])
		row.updateValues = append(row.updateValues, row.Fields[fieldKey])
	}
}

// GetInsertColumnsLength returns number of columns for INSERT query
func (row *Row) GetInsertColumnsLength() int {
	return row.insertColumnLength
}

// GetUpdateColumnsLength returns number of columns for UDPATE query
func (row *Row) GetUpdateColumnsLength() int {
	return row.updateColumnLength
}

// GetInsertColumns returns a slice of column names for INSERT query
func (row *Row) GetInsertColumns() []string {
	return row.insertColumns
}

// GetUpdateColumns returns a slice of column names for UPDATE query
func (row *Row) GetUpdateColumns() []string {
	return row.updateColumns
}

// GetInsertValues returns a slice of values for INSERT query
func (row *Row) GetInsertValues() []interface{} {
	return row.insertValues
}

// GetUpdateValues returns a slice of values for UPDATE query
func (row *Row) GetUpdateValues() []interface{} {
	return row.updateValues
}

// GetInsertPlaceholders returns a slice of placeholders for INSERT query
func (row *Row) GetInsertPlaceholders(driver string) []string {
	placeholders := make([]string, row.GetInsertColumnsLength())
	for i := range row.insertColumns {
		if driver == postgresDriver {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		} else {
			placeholders[i] = "?"
		}
	}
	return placeholders
}

// GetUpdatePlaceholders returns a slice of placeholders for UPDATE query
func (row *Row) GetUpdatePlaceholders(driver string) []string {
	placeholders := make([]string, row.GetInsertColumnsLength())
	for i, c := range row.updateColumns {
		if driver == postgresDriver {
			placeholders[i] = fmt.Sprintf("%s = $%d", c, i+1)
		} else {
			placeholders[i] = fmt.Sprintf("%s = ?", c)
		}
	}
	return placeholders
}

// GetWhere returns a where condition based on primary key with placeholders
func (row *Row) GetWhere(driver string, i int) string {
	wheres := make([]string, len(row.PK))
	j := i
	for _, c := range row.pkColumns {
		if driver == postgresDriver {
			wheres[i-j] = fmt.Sprintf("%s = $%d", c, i+1)
		} else {
			wheres[i-j] = fmt.Sprintf("%s = ?", c)
		}
		i++
	}
	return strings.Join(wheres, " AND ")
}

// GetPKValues returns a slice of primary key values
func (row *Row) GetPKValues() []interface{} {
	return row.pkValues
}
