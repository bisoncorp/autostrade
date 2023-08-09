package sampledata

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"math/rand"
)

//go:embed db_city.csv
var dbCityCsv []byte
var dbCity [][]string

const (
	Istat = iota
	Comune
	Provincia
	Regione
	Prefisso
	CAP
	CodFisco
	Abitanti
	Link
)

func init() {
	reader := csv.NewReader(bytes.NewReader(dbCityCsv))
	reader.Comma = ';'
	all, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	all = all[1:]
	dbCity = all
}

func RandomCityName() string {
	index := rand.Intn(len(dbCity))
	return dbCity[index][Comune]
}
