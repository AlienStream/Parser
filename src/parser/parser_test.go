package parser

import (
   "testing"
   models "github.com/AlienStream/Shared-Go/models"
)

func SourceFactory(source_type string) models.Source {
	var source models.Source

	switch(source_type) {
		case "reddit/subreddit":
			source = models.Source{ 
				Type: "reddit/subreddit",
				Url:  "https://www.reddit.com/r/hiphopheads",
			}
			return source;
	}

	return source;
}

func TestParser(t *testing.T) {
	reddit_source := SourceFactory("reddit/subreddit");

	updateSourceMetaData(&reddit_source);
	if (reddit_source.Title != "HipHopHeads") {
		t.Error("Expected: Title", "HipHopHeads", "Received:", reddit_source.Title)
  	}

  	reddit_posts := fetchNewPosts(reddit_source)
  	if (len(reddit_posts) == 0) {
  		t.Error("Couldn't fetch new posts for reddit source")
  	}
}
