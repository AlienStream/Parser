package parser

import (
	//"fmt"
	models "github.com/AlienStream/Shared-Go/models"
)

type DataObject struct {
	Source models.Source
	Posts  []models.Post
}

func Update(source models.Source) {
	source_data := DataObject{
		source,
		[]models.Post{},
	}

	source_data = source_data.getFreshData()

	// update the source meta info
	source_data.Source.Save()
	// update the posts
	for _, post := range source_data.Posts {
		if post.IsNew() {
			post.Insert()
		} else {
			old_post := models.Post{}.Find(post.Source_id, post.Embed_url)
			post.Id = old_post.Id
			post.Save()
		}
	}
}

func (data DataObject) getFreshData() DataObject {
	if data.Source.Type == "reddit/subreddit" {
		data = getRedditSubredditData(data)
	}
	if data.Source.Type == "youtube/channel" {
		data = getYoutubeChannelData(data)
	}
	if data.Source.Type == "youtube/playlist" {
		data = getYoutubePlaylistData(data)
	}
	if data.Source.Type == "soundcloud/channel" {
		data = getSoundcloudChannelData(data)
	}
	if data.Source.Type == "soundcloud/playlist" {
		data = getSoundcloudPlaylistData(data)
	}
	if data.Source.Type == "blog/rss" {
		data = getBlogRSSData(data)
	}
	return data
}
