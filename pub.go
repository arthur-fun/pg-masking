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
	ColumnNames []string
	ColumnTypes []string
	DataCells   []interface{}
}

type tableReader interface {
	HasRow() bool
	ReadRows() []dataRow
}

type dataMasking interface {
	mask([]dataRow) []dataRow
}

type tableWriter interface {
	WriteRows([]dataRow)
}
