package parser

import (
	models "github.com/AlienStream/Shared-Go/models"
	"parser/reddit_parser"
)

type Parser interface {
	UpdateSourceMetaData(*models.Source)
	FetchPostsFromSource(models.Source) []models.Post
}

func Update(source models.Source) {
	updateSourceMetaData(&source)
	source.Save()

	// update the posts
	posts := fetchNewPosts(source)
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
}

func updateSourceMetaData(source *models.Source) {
	var parser Parser;

	switch (source.Type) {
		case "reddit/subreddit":
			parser = reddit_parser.Parser{}
			break;
		// case "youtube/channel":
		// 	getYoutubeChannelData(data)
		// 	break;
		// case "youtube/playlist":
		// 	getYoutubePlaylistData(data)
		// 	break;
		// case "soundcloud/channel":
		// 	getSoundcloudChannelData(data)
		// 	break;
		// case "reddit/subreddit":
		// 	getSoundcloudPlaylistData(data)
		// 	break;
		// case "blog/rss":
		// 	getBlogRSSData(data)
		// 	break;
	}

	parser.UpdateSourceMetaData(source)
}

func fetchNewPosts(source models.Source) []models.Post {
	var parser Parser;

	switch (source.Type) {
		case "reddit/subreddit":
			parser = reddit_parser.Parser{}
		// case "youtube/channel":
		// 	return getRedditSubredditPosts(source)
		// case "youtube/playlist":
		// 	return getRedditSubredditPosts(source)
		// case "soundcloud/channel":
		// 	return getRedditSubredditPosts(source)
		// case "reddit/subreddit":
		// 	return getRedditSubredditPosts(source)
		// case "blog/rss":
		// 	return getRedditSubredditPosts(source)
	}

	return parser.FetchPostsFromSource(source)
}

