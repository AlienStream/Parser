package soundcloud_parser

import (
	"testing"

	models "github.com/AlienStream/Shared-Go/models"
)

func TestChannels(t *testing.T) {
	soundcloud_channel_source := SourceFactory("soundcloud/channel")

	update_err := Parser{}.UpdateSourceMetaData(&soundcloud_channel_source)
	if update_err != nil || soundcloud_channel_source.Title != "LIL UZI VERT" {
		t.Error("Expected: Title", "LIL UZI VERT", "Received:", soundcloud_channel_source.Title)
	}

	soundcloud_channel_posts, err := Parser{}.FetchPostsFromSource(soundcloud_channel_source)
	if err != nil || len(soundcloud_channel_posts) == 0 {
		t.Error("Couldn't fetch new posts for soundcloud source")
	}

	// Try a restricted one
	soundcloud_channel_source.Url = "https://soundcloud.com/bassnectar"
	soundcloud_channel_posts, _ = Parser{}.FetchPostsFromSource(soundcloud_channel_source)
	if len(soundcloud_channel_posts) != 0 {
		t.Error("Restricted soundcloud source returning tracks")
	}

	// Try Some Broken Ones
	soundcloud_channel_source.Url = "https://soundcloud.com/bassnectarasdsa"
	update_err = Parser{}.UpdateSourceMetaData(&soundcloud_channel_source)
	if update_err == nil {
		t.Error("Error Not Properly being returned")
	}

	soundcloud_channel_source.Url = "https://soundcloud.com/bassnectarasdsa"
	soundcloud_channel_posts, _ = Parser{}.FetchPostsFromSource(soundcloud_channel_source)
	if len(soundcloud_channel_posts) != 0 {
		t.Error("Broken soundcloud source returning tracks")
	}
}

func TestPlaylists(t *testing.T) {
	soundcloud_playlist_source := SourceFactory("soundcloud/playlist")

	update_err := Parser{}.UpdateSourceMetaData(&soundcloud_playlist_source)
	if update_err != nil || soundcloud_playlist_source.Title != "THE PERFECT LUV TAPE®️" {
		t.Error("Expected: Title", "THE PERFECT LUV TAPE®️", "Received:", soundcloud_playlist_source.Title)
	}

	soundcloud_playlist_posts, err := Parser{}.FetchPostsFromSource(soundcloud_playlist_source)
	if err != nil || len(soundcloud_playlist_posts) == 0 {
		t.Error("Couldn't fetch new posts for soundcloud source")
	}

	// Try Some Broken Ones
	soundcloud_playlist_source.Url = "https://soundcloud.com/sets/bassnectarasdsa"
	update_err = Parser{}.UpdateSourceMetaData(&soundcloud_playlist_source)
	if update_err == nil {
		t.Error("Error Not Properly being returned")
	}

}

func SourceFactory(source_type string) models.Source {
	var source models.Source

	switch source_type {
	case "soundcloud/channel":
		source = models.Source{
			Type: source_type,
			Url:  "https://soundcloud.com/liluzivert",
		}
		return source
	case "soundcloud/playlist":
		source = models.Source{
			Type: source_type,
			Url:  "https://soundcloud.com/sets/the-perfect-luv-tape-r",
		}
		return source
	}

	return source
}
