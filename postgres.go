package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

//connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"

func newPostgresTableReader(connStr string, tableName string) tableReader {
	return &postgresTableReader{connStr, tableName, -1, 100, 0}
}

func newPostgresTableWriter(connStr string, tableName string) tableWriter {
	return &postgresTableWriter{connStr, tableName}
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
		sqlStr := fmt.Sprintf("select count(*) from %s", pgReader.tableName)
		err = db.QueryRow(sqlStr).Scan(&pgReader.rowCount)
		checkErr(err)
	}
	return pgReader.rowCount > pgReader.cursor*pgReader.pageSize
}

func (pgReader *postgresTableReader) ReadRows() []dataRow {
	db, err := sql.Open("postgres", pgReader.connStr)
	checkErr(err)
	defer db.Close()

	sqlStr := fmt.Sprintf("select * from %s limit %d offset %d", pgReader.tableName, pgReader.pageSize, pgReader.cursor)
	rows, err := db.Query(sqlStr)
	checkErr(err)
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
		result[i].ColumnNames = colNames
		result[i].ColumnTypes = translateColumnDataTypes(colTypes)
		result[i].DataCells = copyRowData(cellValues, colTypes)
		i++
	}
	return result
}

type postgresTableWriter struct {
	connStr   string
	tableName string
}

func (pgWriter *postgresTableWriter) WriteRows(rows []dataRow) {
	db, err := sql.Open("postgres", pgWriter.connStr)
	checkErr(err)
	defer db.Close()

	sqlStr := makeInsertSQL(pgWriter.tableName, rows)
	//fmt.Println(sqlStr)
	result, err := db.Exec(sqlStr)
	checkErr(err)
	affected, err := result.RowsAffected()
	checkErr(err)
	if int(affected) != len(rows) {
		panic(fmt.Sprintf("The inserted rows(%d) are not equal to the input rows(%d)", affected, len(rows)))
	}
}

func makeInsertSQL(tableName string, rows []dataRow) string {
	if len(rows) < 1 {
		return ""
	}
	buf := bytes.Buffer{}
	buf.WriteString("INSERT INTO ")
	buf.WriteString(tableName)
	buf.WriteString(" ( ")
	for i := range rows[0].ColumnNames {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(rows[0].ColumnNames[i])
	}
	buf.WriteString(") VALUES\n")
	for i := range rows {
		if i > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString("(")
		for j := range rows[i].DataCells {
			if j > 0 {
				buf.WriteString(", ")
			}
			if rows[i].ColumnTypes[j] == "BPCHAR" ||
				rows[i].ColumnTypes[j] == "CHAR" ||
				rows[i].ColumnTypes[j] == "VARCHAR" ||
				rows[i].ColumnTypes[j] == "DATE" ||
				rows[i].ColumnTypes[j] == "TEXT" ||
				rows[i].ColumnTypes[j] == "TIME" ||
				rows[i].ColumnTypes[j] == "TIMETZ" ||
				rows[i].ColumnTypes[j] == "TIMESTAMPTZ" ||
				rows[i].ColumnTypes[j] == "TIMESTAMP" {
				buf.WriteString("'")
			}
			if ptr, ok := (rows[i].DataCells[j]).(*interface{}); ok {
				v := *ptr
				if v != nil {
					switch v.(type) {
					case string:
						buf.WriteString(v.(string))
					case int:
						buf.WriteString(strconv.Itoa(v.(int)))
					case int64:
						buf.WriteString(strconv.FormatInt(v.(int64), 10))
					case float32:
						buf.WriteString(strconv.FormatFloat(v.(float64), 'e', -1, 32))
					case float64:
						buf.WriteString(strconv.FormatFloat(v.(float64), 'e', -1, 64))
					case bool:
						buf.WriteString(strconv.FormatBool(v.(bool)))
					case []uint8:
						buf.WriteString(string(v.([]uint8)))
					case time.Time:
						tv := v.(time.Time)
						switch rows[i].ColumnTypes[j] {
						case "DATE":
							buf.WriteString(tv.Format("2006-01-02"))
						case "TIME":
							buf.WriteString(tv.Format("15:04:05"))
						case "TIMETZ":
							buf.WriteString(tv.Format("15:04:05+8"))
						case "TIMESTAMPTZ":
							buf.WriteString(tv.Format("2006-01-02 15:04:05.123+8"))
						case "TIMESTAMP":
							buf.WriteString(tv.Format("2006-01-02 15:04:05.123"))
						}
					default:
						panic("Unsupported data type")
					}
				}
			} else {
				panic("Failed to cast to *interface{}")
			}
			if rows[i].ColumnTypes[j] == "BPCHAR" ||
				rows[i].ColumnTypes[j] == "CHAR" ||
				rows[i].ColumnTypes[j] == "VARCHAR" ||
				rows[i].ColumnTypes[j] == "DATE" ||
				rows[i].ColumnTypes[j] == "TEXT" ||
				rows[i].ColumnTypes[j] == "TIME" ||
				rows[i].ColumnTypes[j] == "TIMETZ" ||
				rows[i].ColumnTypes[j] == "TIMESTAMPTZ" ||
				rows[i].ColumnTypes[j] == "TIMESTAMP" {
				buf.WriteString("'")
			}
		}
		buf.WriteString(")")
	}
	return buf.String()
}

func initDataBuffer(dbTypes []*sql.ColumnType) []interface{} {
	cellValues := make([]interface{}, len(dbTypes))
	for i := range cellValues {
		var v interface{}
		cellValues[i] = &v
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
