package soundcloud_parser

import (
	"testing"

	models "github.com/AlienStream/Shared-Go/models"
)

func SourceFactory(source_type string) models.Source {
	var source models.Source

	switch source_type {
	case "soundcloud/channel":
		source = models.Source{
			Type: source_type,
			Url:  "https://soundcloud.com/liluzivert",
		}
		return source

	}

	return source
}

func TestParser(t *testing.T) {
	reddit_source := SourceFactory("reddit/subreddit")

	updateSourceMetaData(&reddit_source)
	if reddit_source.Title != "HipHopHeads" {
		t.Error("Expected: Title", "HipHopHeads", "Received:", reddit_source.Title)
	}

	reddit_posts := fetchNewPosts(reddit_source)
	if len(reddit_posts) == 0 {
		t.Error("Couldn't fetch new posts for reddit source")
	}

	soundcloud_channel_source := SourceFactory("soundcloud/channel")

	updateSourceMetaData(&soundcloud_channel_source)
	if soundcloud_channel_source.Title != "LIL UZI VERT" {
		t.Error("Expected: Title", "LIL UZI VERT", "Received:", soundcloud_channel_source.Title)
	}

	soundcloud_channel_posts := fetchNewPosts(soundcloud_channel_source)
	if len(soundcloud_channel_posts) == 0 {
		t.Error("Couldn't fetch new posts for soundcloud source")
	}

	// Try a restricted one
	soundcloud_channel_source.Url = "https://soundcloud.com/bassnectar"
	soundcloud_channel_posts = fetchNewPosts(soundcloud_channel_source)
	if len(soundcloud_channel_posts) != 0 {
		t.Error("Restricted soundcloud source returning tracks")
	}

	// Try Some Broken Ones
	reddit_source.Url = "https://soundcloud.com/bassnectarasdsa"
	updateSourceMetaData(&reddit_source)

	soundcloud_channel_source.Url = "https://soundcloud.com/bassnectarasdsa"
	soundcloud_channel_posts = fetchNewPosts(soundcloud_channel_source)
	if len(soundcloud_channel_posts) != 0 {
		t.Error("Restricted soundcloud source returning tracks")
	}
}
