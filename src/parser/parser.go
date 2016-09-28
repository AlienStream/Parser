package parser

import (
	models "github.com/AlienStream/Shared-Go/models"
	"parser/reddit_parser"
)

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
	switch (source.Type) {
		case "reddit/subreddit":
			reddit_parser.UpdateSourceMetaData(source)
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
}

func fetchNewPosts(source models.Source) []models.Post {
	switch (source.Type) {
		case "reddit/subreddit":
			return reddit_parser.FetchPostsFromSource(source)
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

	panic("Invalid Source Type");
}

