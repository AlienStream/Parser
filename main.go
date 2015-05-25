package main

import (
	"github.com/AlienStream/Parser/parser"
	db "github.com/AlienStream/Shared-Go/database"
	models "github.com/AlienStream/Shared-Go/models"
)

func main() {
	updateALLSources()
}

func updateALLSources() {
	var sources []models.Source = models.AllSources()
	for _, source := range sources {
		parser.Update(source)
	}

}

func updateExpiredSources(refresh_freq int) {
	rows, _, err := db.Con.Query("select * from sources where refresh_frequency >= %d", refresh_freq)
	if err != nil {
		panic(err)
	}

	var sources []models.Source = models.RowsToSources(rows)
	for _, source := range sources {
		go parser.Update(source)
	}

}