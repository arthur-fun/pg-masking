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
			if len(params) != 2 {
				panic("Failed to parse the parameters of replaceConverter")
			}
			start, err := strconv.Atoi(strings.Trim(params[1], " "))
			checkErr(err)
			converter = replaceConverter{params[0], start}
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
			gConverter, gOk := dc.globalConverters[cellName]
			tConverter, tOk := dc.tableConverters[tableName+"___"+cellName]
			if gOk {
				newRows[i].DataCells[j] = gConverter.mask(rows[i].DataCells[j])
			}
			if tOk {
				newRows[i].DataCells[j] = tConverter.mask(rows[i].DataCells[j])
			}
			if !gOk && !tOk {
				newRows[i].DataCells[j] = rows[i].DataCells[j]
			}
		}
	}
	return newRows
}
