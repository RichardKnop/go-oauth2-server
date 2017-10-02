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
	rawPkValues        []interface{}
	insertColumns      []string
	updateColumns      []string
	insertValues       []interface{}
	updateValues       []interface{}
	rawInsertValues    []interface{}
	rawUpdateValues    []interface{}
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
	row.rawInsertValues = make([]interface{}, 0)
	row.rawUpdateValues = make([]interface{}, 0)
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
		row.appendValue("pk", row.PK[pkKey])
		row.insertColumns = append(row.insertColumns, pkKey)
		row.updateColumns = append(row.updateColumns, pkKey)
		row.appendValue("insert", row.PK[pkKey])
		row.appendValue("update", row.PK[pkKey])
	}

	// Rest of the fields
	for _, fieldKey := range fieldKeys {
		sv, ok := row.Fields[fieldKey].(string)
		if ok && sv == onInsertNow {
			row.insertColumns = append(row.insertColumns, fieldKey)
			row.appendValue("insert", time.Now())
			row.updateColumnLength--
			continue
		}
		if ok && sv == onUpdateNow {
			row.updateColumns = append(row.updateColumns, fieldKey)
			row.appendValue("update", time.Now())
			row.insertColumnLength--
			continue
		}
		row.insertColumns = append(row.insertColumns, fieldKey)
		row.updateColumns = append(row.updateColumns, fieldKey)
		row.appendValue("insert", row.Fields[fieldKey])
		row.appendValue("update", row.Fields[fieldKey])
	}
}

// GetInsertColumnsLength returns number of columns for INSERT query
func (row *Row) GetInsertColumnsLength() int {
	return row.insertColumnLength
}

// GetInsertValuesLength returns number of values for INSERT query
func (row *Row) GetInsertValuesLength() int {
	return len(row.insertValues)
}

// GetUpdateColumnsLength returns number of columns for UDPATE query
func (row *Row) GetUpdateColumnsLength() int {
	return row.updateColumnLength
}

// GetUpdateValuesLength returns number of values for UDPATE query
func (row *Row) GetUpdateValuesLength() int {
	return len(row.updateValues)
}

// GetInsertColumns returns a slice of column names for INSERT query
func (row *Row) GetInsertColumns() []string {
	escapedColumns := make([]string, len(row.insertColumns))
	for i, insertColumn := range row.insertColumns {
		escapedColumns[i] = fmt.Sprintf("\"%s\"", insertColumn)
	}
	return escapedColumns
}

// GetUpdateColumns returns a slice of column names for UPDATE query
func (row *Row) GetUpdateColumns() []string {
	escapedColumns := make([]string, len(row.updateColumns))
	for i, updateColumn := range row.updateColumns {
		escapedColumns[i] = fmt.Sprintf("\"%s\"", updateColumn)
	}
	return escapedColumns
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
	for i, j := 0, 0; i < row.GetInsertColumnsLength(); i++ {
		val := row.rawInsertValues[i]
		switch v := val.(type) {
		case string:
			if strings.HasPrefix(v, "RAW=") {
				placeholders[i] = strings.TrimPrefix(v, "RAW=")
				continue
			}
		}
		if driver == postgresDriver {
			placeholders[i] = fmt.Sprintf("$%d", j+1)
			j++
		} else {
			placeholders[i] = "?"
		}
	}
	return placeholders
}

// GetUpdatePlaceholders returns a slice of placeholders for UPDATE query
func (row *Row) GetUpdatePlaceholders(driver string) []string {
	placeholders := make([]string, row.GetUpdateColumnsLength())
	j := 0
	for i, c := range row.GetUpdateColumns() {
		val := row.rawUpdateValues[i]
		switch v := val.(type) {
		case string:
			if strings.HasPrefix(v, "RAW=") {
				placeholders[i] = fmt.Sprintf("%s = %s", c, strings.TrimPrefix(v, "RAW="))
				continue
			}
		}
		if driver == postgresDriver {
			placeholders[i] = fmt.Sprintf("%s = $%d", c, j+1)
			j++
		} else {
			placeholders[i] = fmt.Sprintf("%s = ?", c)
		}
	}
	return placeholders
}

// GetWhere returns a where condition based on primary key with placeholders
func (row *Row) GetWhere(driver string, i int) string {
	wheres := make([]string, len(row.PK))
	start, j := i, i
	for _, c := range row.pkColumns {
		val := row.rawPkValues[i-start]
		i++
		switch v := val.(type) {
		case string:
			if strings.HasPrefix(v, "RAW=") {
				wheres[i-1-start] = fmt.Sprintf("%s = %s", c, strings.TrimPrefix(v, "RAW="))

				continue
			}
		}
		if driver == postgresDriver {
			wheres[i-1-start] = fmt.Sprintf("%s = $%d", c, j+1)
			j++
		} else {
			wheres[i-1-start] = fmt.Sprintf("%s = ?", c)
		}

	}
	return strings.Join(wheres, " AND ")
}

// GetPKValues returns a slice of primary key values
func (row *Row) GetPKValues() []interface{} {
	return row.pkValues
}

func (row *Row) appendValue(queryType string, val interface{}) {
	sv, ok := val.(string)
	if !ok || !strings.HasPrefix(sv, "RAW=") {
		switch queryType {
		case "insert":
			row.insertValues = append(row.insertValues, val)
		case "update":
			row.updateValues = append(row.updateValues, val)
		case "pk":
			row.pkValues = append(row.pkValues, val)
		}
	}
	switch queryType {
	case "insert":
		row.rawInsertValues = append(row.rawInsertValues, val)
	case "update":
		row.rawUpdateValues = append(row.rawUpdateValues, val)
	case "pk":
		row.rawPkValues = append(row.rawPkValues, val)
	}
}
