package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

//connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"

func newPostgresTableReader(connStr string, tableName string) postgresTableReader {
	return postgresTableReader{connStr, tableName, -1, 100, 0}
}

type postgresTableReader struct {
	connStr   string
	tableName string
	rowCount  int
	pageSize  int
	cursor    int
}

func (pgReader *postgresTableReader) HasRow() bool {
	if pgReader.rowCount < 0 {
		db, err := sql.Open("postgres", pgReader.connStr)
		checkErr(err)
		defer db.Close()

		err = db.QueryRow(`select count(*) from caas_realm`).Scan(&pgReader.rowCount)
		checkErr(err)
		fmt.Println("Table ", pgReader.tableName, " has ", pgReader.rowCount, "rows.")
	}
	return pgReader.rowCount > 0
}

func initDataBuffer(dbTypes []*sql.ColumnType) []interface{} {
	cellValues := make([]interface{}, len(dbTypes))
	for i := range cellValues {
		switch dbTypes[i].DatabaseTypeName() {
		case "INT4":
			cellValues[i] = new(int)
		default:
			cellValues[i] = new(string)
		}
	}
	return cellValues
}

func translateColumnDataTypes(dbColumnTypes []*sql.ColumnType) []string {
	columnTypeStr := make([]string, len(dbColumnTypes))
	for i := range dbColumnTypes {
		columnTypeStr[i] = dbColumnTypes[i].DatabaseTypeName()
	}
	return columnTypeStr
}

func copyRowData(src []interface{}, dbColumnTypes []*sql.ColumnType) []interface{} {
	result := initDataBuffer(dbColumnTypes)
	for i := range src {
		result[i] = src[i]
	}
	return result
}

func (pgReader *postgresTableReader) ReadRows() []dataRow {
	db, err := sql.Open("postgres", pgReader.connStr)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(`select * from %s `, pgReader.tableName))
	checkErr(err)
	fmt.Println("Table ", pgReader.tableName, " has ", pgReader.rowCount, "rows.")
	length := pgReader.pageSize
	if pgReader.rowCount < pgReader.pageSize {
		length = pgReader.rowCount
	}
	result := make([]dataRow, length, length)
	colNames, err := rows.Columns()
	checkErr(err)
	colTypes, err := rows.ColumnTypes()
	checkErr(err)

	i := 0
	for rows.Next() {
		cellValues := initDataBuffer(colTypes)
		rows.Scan(cellValues...)
		fmt.Println(cellValues)
		result[i].ColumnNames = colNames
		result[i].ColumnTypes = translateColumnDataTypes(colTypes)
		result[i].DataCells = copyRowData(cellValues, colTypes)
		i++
	}
	return result
}
