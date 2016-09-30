package parser

import (
	"parser/reddit_parser"
	"parser/soundcloud_parser"

	models "github.com/AlienStream/Shared-Go/models"
	"github.com/bugsnag/bugsnag-go"
)

type Parser interface {
	UpdateSourceMetaData(*models.Source) error
	FetchPostsFromSource(models.Source) ([]models.Post, error)
}

func Update(source models.Source) {
	update_err := updateSourceMetaData(&source)
	if update_err == nil {
		source.Save()
	} else {
		bugsnag.Notify(update_err)
	}

	// update the posts
	posts, err := fetchNewPosts(source)
	if err == nil {
		for _, post := range posts {
			if post.IsNew() {
				post.Insert()
			} else {
				old_post := models.Post{}.Find(post.Source_id, post.Embed_url)
				post.Id = old_post.Id
				post.Is_new = true
				post.Save()
			}
		}
	} else {
		bugsnag.Notify(err)
	}
}

// TODO: each parser should be a singleton that can block on a per source basis
// so we can Multithread requests to multiple sources at the same time
func updateSourceMetaData(source *models.Source) error {
	var parser Parser

	switch source.Type {
	case "reddit/subreddit":
		parser = reddit_parser.Parser{}
		break
		// case "youtube/channel":
		// 	getYoutubeChannelData(data)
		// 	break;
		// case "youtube/playlist":
		// 	getYoutubePlaylistData(data)
		// 	break;
	case "soundcloud/channel":
		parser = soundcloud_parser.Parser{}
		break
		// case "reddit/subreddit":
		// 	getSoundcloudPlaylistData(data)
		// 	break;
		// case "blog/rss":
		// 	getBlogRSSData(data)
		// 	break;
	}

	return parser.UpdateSourceMetaData(source)
}

func fetchNewPosts(source models.Source) ([]models.Post, error) {
	var parser Parser

	switch source.Type {
	case "reddit/subreddit":
		parser = reddit_parser.Parser{}
		break
	// case "youtube/channel":
	// 	return getRedditSubredditPosts(source)
	// case "youtube/playlist":
	// 	return getRedditSubredditPosts(source)
	case "soundcloud/channel":
		parser = soundcloud_parser.Parser{}
		// case "reddit/subreddit":
		// 	return getRedditSubredditPosts(source)
		// case "blog/rss":
		// 	return getRedditSubredditPosts(source)
	}

	return parser.FetchPostsFromSource(source)
}
