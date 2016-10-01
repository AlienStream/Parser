package main

import (
	"parser"
	"time"

	"github.com/bugsnag/bugsnag-go"

	db "github.com/AlienStream/Shared-Go/database"
	models "github.com/AlienStream/Shared-Go/models"
)

func main() {
	source_types := []string{
		"reddit",
		"soundcloud",
	}

	update_interval := time.Minute * 15

	bugsnag.Configure(bugsnag.Configuration{
		APIKey: "14d6cbfd67182668a12d97372c33f13e",
	})

	for _, source_type := range source_types {
		// select sources by type, then put each queue on it's own thread
		go func(source_type string) {
			rows, _, err := db.Con.Query("select * from sources where type LIKE %s", "%"+source_type+"%")
			if err != nil {
				panic(err)
			}

			sources := models.RowsToSources(rows)
			for _, source := range sources {
				parser.Update(source) // Refresh automatically on boot
			}

			refreshinterval := time.NewTicker(update_interval).C
			for {
				select {
				case <-refreshinterval:
					for _, source := range sources {
						parser.Update(source)
					}
					break
				}
			}
		}()
	}
}
