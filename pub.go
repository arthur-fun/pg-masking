package main

import (
	_ "github.com/lib/pq"
)

type maskingType int

const (
	replace = iota
	shuffling
	random
)

type cellConverter interface {
	mask(interface{}) interface{}
}

type dataRow struct {
	TableName   string
	ColumnNames []string
	ColumnTypes []string
	DataCells   []interface{}
}

type tableReader interface {
	TableName() string
	HasRow() bool
	ReadRows() []dataRow
}

type dataMasking interface {
	mask(tableName string, rows []dataRow) []dataRow
}

type tableWriter interface {
	WriteRows([]dataRow)
}
