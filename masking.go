package main

import (
	"strconv"
	"strings"
)

func newDataMasker(config []converterConfig) dataMasking {
	gConvs := make(map[string]cellConverter)
	tConvs := make(map[string]cellConverter)

	for i := range config {
		var converter cellConverter
		switch config[i].Converter {
		case "ReplaceAll":
			converter = replaceAllConverter{config[i].ConvParameters}
		case "Replace":
			params := strings.Split(config[i].ConvParameters, ",")
			if len(params) != 3 {
				panic("Failed to parse the parameters of replaceConverter")
			}
			start, err := strconv.Atoi(params[1])
			checkErr(err)
			length, err := strconv.Atoi(params[2])
			checkErr(err)
			converter = replaceConverter{params[0], start, length}
		case "Random":
			converter = randomConverter{}
		default:
			continue
		}
		if "*" == config[i].TableName {
			gConvs[config[i].ColumnName] = converter
		} else {
			tConvs[config[i].TableName+"___"+config[i].ColumnName] = converter
		}
	}
	return dataConveters{globalConverters: gConvs, tableConverters: tConvs}
}

type dataConveters struct {
	globalConverters map[string]cellConverter
	tableConverters  map[string]cellConverter
}

func (dc dataConveters) mask(tableName string, rows []dataRow) []dataRow {
	newRows := make([]dataRow, len(rows), len(rows))
	for i := range rows {
		newRows[i] = rows[i]
		newRows[i].DataCells = make([]interface{}, len(rows[i].DataCells))
		for j := range rows[i].DataCells {
			cellName := rows[i].ColumnNames[j]
			if converter, ok := dc.globalConverters[cellName]; ok {
				newRows[i].DataCells[j] = converter.mask(rows[i].DataCells[j])
			} else if converter, ok := dc.tableConverters[tableName+"___"+cellName]; ok {
				newRows[i].DataCells[j] = converter.mask(rows[i].DataCells[j])
			} else {
				newRows[i].DataCells[j] = rows[i].DataCells[j]
			}
		}
	}
	return newRows
}
