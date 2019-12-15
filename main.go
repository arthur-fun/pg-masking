package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type dbInfo struct {
	DbScheme string `json:"dbms"`
	DbURL    string `json:"dburl"`
}

type converterConfig struct {
	TableName      string `json:"table-name"`
	ColumnName     string `json:"column-name"`
	Converter      string `json:"converter"`
	ConvParameters string `json:"converter-parameters"`
}

type maskingConfig struct {
	Tables     []string          `json:"tables"`
	SrcDb      dbInfo            `json:"source"`
	DestDb     dbInfo            `json:"destination"`
	Converters []converterConfig `json:"column-converter"`
}

var configFileName *string

func init() {
	fmt.Printf("Start to init....\n")
	configFileName = flag.String("f", "config.json", "Configuration file")
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Database error: %v", err)
		panic(err)
	}
}

func showRows(rs []dataRow) {
	for i := range rs {
		fmt.Println(rs[i].ColumnNames)
		fmt.Println(rs[i].ColumnTypes)
		for j := range rs[i].DataCells {
			if vRef, ok := rs[i].DataCells[j].(*int); ok {
				fmt.Print(*vRef)
				fmt.Print("\t")
			} else {
				fmt.Println(rs[i].DataCells[j])
			}
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Hellow world!\n")
	fmt.Println("Configuration file name is", *configFileName)

	config := &maskingConfig{}
	f, err := os.Open(*configFileName)
	defer f.Close()
	checkErr(err)

	contentByte, err := ioutil.ReadAll(f)
	checkErr(err)

	err = json.Unmarshal(contentByte, config)
	checkErr(err)

	fmt.Println(config)

	pgReader := newPostgresTableReader(config.SrcDb.DbURL, "caas_realm")
	if pgReader.HasRow() {
		fmt.Println("Table caas_realm has records.")
		rows := pgReader.ReadRows()
		showRows(rows)
		masker := newDataMasker(config.Converters)
		newRows := masker.mask(rows)
		showRows(newRows)
	} else {
		fmt.Println("Table caas_realm has no records.")
	}
}
