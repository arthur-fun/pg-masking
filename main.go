package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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
			if vRef, ok := rs[i].DataCells[j].(*interface{}); ok {
				fmt.Print(*vRef)
				fmt.Print("\t")
			}
		}
		fmt.Println()
	}
}

func allTableNames(config *maskingConfig) []string {
	for i := range config.Tables {
		if config.Tables[i] == "*" {
			return queryAllTableNames(config.SrcDb.DbURL)
		}
	}
	return config.Tables
}

func processTable(masker *dataMasking, srcDBURL, destDBURL string, tableName string, wg *sync.WaitGroup) {
	fmt.Printf("Starting to mask %s.\n", tableName)
	pgReader := newPostgresTableReader(srcDBURL, tableName)
	pgWriter := newPostgresTableWriter(destDBURL, tableName)
	if pgReader.HasRow() {
		for pgReader.HasRow() {
			rows := pgReader.ReadRows()
			fmt.Println("The data loaded from source via reader...")
			showRows(rows)
			fmt.Println("Done to load data.")
			newRows := (*masker).mask(tableName, rows)
			fmt.Println("The data after converted")
			showRows(newRows)
			fmt.Println("Done to print data converted")
			pgWriter.WriteRows(newRows)
		}
	} else {
		fmt.Printf("Table %s has no records.\n", tableName)
	}
	wg.Done()
}

func main() {
	flag.Parse()
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

	masker := newDataMasker(config.Converters)

	wg := sync.WaitGroup{}
	tableNames := allTableNames(config)
	for i := range tableNames {
		wg.Add(1)
		go processTable(&masker, config.SrcDb.DbURL, config.DestDb.DbURL, tableNames[i], &wg)
	}
	wg.Wait()
}
