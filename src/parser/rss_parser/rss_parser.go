package rss_parser

import (
	"fmt"
	//"net/http"
)

func getBlogRSSData(source_data *DataObject) DataObject {
	fmt.Printf("Updating %s \n", source_data.Source.Title)
	return DataObject{}
}
