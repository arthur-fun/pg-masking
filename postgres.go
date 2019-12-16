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

func (pgReader *postgresTableReader) TableName() string {
	return pgReader.tableName
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
	return pgReader.rowCount > pgReader.cursor*pgReader.pageSize
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

	rows, err := db.Query(fmt.Sprintf(`select * from %s limit %s offset %s`, pgReader.tableName, pgReader.pageSize, pgReader.cursor))
	checkErr(err)
	fmt.Println("Table ", pgReader.tableName, " has ", pgReader.rowCount, "rows.")
	length := pgReader.pageSize
	if pgReader.rowCount < pgReader.pageSize*(1+pgReader.cursor) {
		length = pgReader.rowCount - pgReader.pageSize*pgReader.cursor
	}
	pgReader.cursor++
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

func queryAllTableNames(dbURL string) []string {
	db, err := sql.Open("postgres", dbURL)
	checkErr(err)
	defer db.Close()

	var cnt int
	err = db.QueryRow(`SELECT count(tablename) FROM pg_tables WHERE tablename NOT LIKE 'pg%' AND tablename NOT LIKE 'sql_%'`).Scan(&cnt)
	checkErr(err)

	rows, err := db.Query("SELECT tablename FROM pg_tables WHERE tablename NOT LIKE 'pg%' AND tablename NOT LIKE 'sql_%' ORDER BY tablename")
	checkErr(err)

	tableNames := make([]string, cnt)
	var tName string
	i := 0
	for rows.Next() {
		rows.Scan(&tName)
		tableNames[i] = tName
		i++
	}
	return tableNames
}
