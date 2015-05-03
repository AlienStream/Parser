package main

import (
	"github.com/AlienStream/Parser/parser"
	con "github.com/AlienStream/Shared-Go/database"
	models "github.com/AlienStream/Shared-Go/models"
	mysql "github.com/ziutek/mymysql/mysql"
)

func main() {
	updateALLSources()
}

func updateALLSources() {
	db := con.GetDBConnection()
	rows, _, err := db.Query("select * from sources")
	if err != nil {
		panic(err)
	}

	var sources []models.Source = RowsToObjects(rows, "source")
	for _, source := range sources {
		parser.Update(source)
	}

}

func updateExpiredSources(refresh_freq int) {
	db := con.GetDBConnection()
	rows, _, err := db.Query("select * from sources where refresh_frequency >= %d", refresh_freq)
	if err != nil {
		panic(err)
	}

	var sources []models.Source = RowsToObjects(rows, "source")
	for _, source := range sources {
		go parser.Update(source)
	}

}

func RowsToObjects(rows []mysql.Row, object string) []models.Source {
	var sources = []models.Source{}
	if object == "source" {
		for _, row := range rows {
			sources = append(sources, RowToSource(row))
		}

	}
	return sources

}

func RowToSource(row mysql.Row) models.Source {
	var source = models.Source{
		Id:          row.Int(0),
		Title:       row.Str(1),
		Description: row.Str(2),
		Type:        row.Str(3),
		Importance:  row.Int(4),
		Url:         row.Str(5),
		Thumbnail:   row.Str(6),
		Updated_at:  row.Localtime(7),
		Created_at:  row.Localtime(8),
	}
	return source
}
