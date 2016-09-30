package reddit_parser

import (
	"testing"

	models "github.com/AlienStream/Shared-Go/models"
)

func TestParser(t *testing.T) {
	reddit_source := SourceFactory("reddit/subreddit")

	update_err := Parser{}.UpdateSourceMetaData(&reddit_source)
	if update_err != nil || reddit_source.Title != "HipHopHeads" {
		t.Error("Expected: Title", "HipHopHeads", "Received:", reddit_source.Title)
	}

	reddit_posts, err := Parser{}.FetchPostsFromSource(reddit_source)
	if err != nil || len(reddit_posts) == 0 {
		t.Error("Couldn't fetch new posts for reddit source")
	}

	// Broken Sources need to fail gracefully
	reddit_source.Url = "http://www.google.com"
	err = Parser{}.UpdateSourceMetaData(&reddit_source)
	if err == nil {
		t.Error("Error Not Returned Properly")
	}

	_, err = Parser{}.FetchPostsFromSource(reddit_source)
	if err == nil {
		t.Error("Error Not Returned Properly")
	}

	reddit_source.Url = "http://www.foo"
	err = Parser{}.UpdateSourceMetaData(&reddit_source)
	if err == nil {
		t.Error("Error Not Returned Properly")
	}
	_, err = Parser{}.FetchPostsFromSource(reddit_source)
	if err == nil {
		t.Error("Error Not Returned Properly")
	}
}

func SourceFactory(source_type string) models.Source {
	var source models.Source
	switch source_type {
	case "reddit/subreddit":
		source = models.Source{
			Type: source_type,
			Url:  "https://www.reddit.com/r/hiphopheads",
		}
		return source
	}
	return source
}
