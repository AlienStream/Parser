package rss_parser

import (
	"testing"

	models "github.com/AlienStream/Shared-Go/models"
)

func TestParser(t *testing.T) {
	rss_source := SourceFactory("rss/blog")

	update_err := Parser{}.UpdateSourceMetaData(&rss_source)
	if update_err != nil || rss_source.Title != "This Song Is Sick" {
		t.Error("Expected: Title", "This Song Is Sick", "Received:", rss_source.Title)
	}

	rss_posts, err := Parser{}.FetchPostsFromSource(rss_source)
	if err != nil || len(rss_posts) == 0 {
		t.Error("Couldn't fetch new posts for rss source")
	}

	// Broken Sources need to fail gracefully
	rss_source.Url = "http://www.google.com"
	err = Parser{}.UpdateSourceMetaData(&rss_source)
	if err == nil {
		t.Error("Error Not Returned Properly")
	}

	_, err = Parser{}.FetchPostsFromSource(rss_source)
	if err == nil {
		t.Error("Error Not Returned Properly")
	}
}

func SourceFactory(source_type string) models.Source {
	var source models.Source
	switch source_type {
	case "rss/blog":
		source = models.Source{
			Type: source_type,
			Url:  "http://thissongissick.com/feeds/feed",
		}
		return source
	}
	return source
}
