package rss_parser

import (
	"regexp"
	"strings"

	models "github.com/AlienStream/Shared-Go/models"
	"github.com/SlyMarbo/rss"
)

func (Parser) UpdateSourceMetaData(source *models.Source) error {
	rss.CacheParsedItemIDs(false)
	feed, err := rss.Fetch(source.Url)

	if err == nil {
		source.Title = feed.Title
		source.Description = feed.Description
		source.Thumbnail = feed.Image.Url
	}

	return err
}

func (Parser) FetchPostsFromSource(source models.Source) ([]models.Post, error) {
	posts := []models.Post{}

	rss.CacheParsedItemIDs(false)
	feed, err := rss.Fetch(source.Url)
	if err != nil {
		return posts, err
	}

	iframe_regex := regexp.MustCompile("<iframe src=\"(?P<url>.*?)\".*?>.*?</iframe>")
	for _, item := range feed.Items {
		embed_url := iframe_regex.FindStringSubmatch(item.Summary)
		if embed_url == nil || !urlIsEmbeddable(embed_url[1]) {
			embed_url = iframe_regex.FindStringSubmatch(item.Content)
		}
		if embed_url != nil && urlIsEmbeddable(embed_url[1]) {
			post := models.Post{
				Id:                 0,
				Source_id:          source.Id,
				Title:              item.Title,
				Number_of_comments: 0,
				Permalink:          item.Link,
				Thumbnail:          feed.Image.Url,
				Embed_url:          embed_url[1],
				Likes:              0,
				Dislikes:           0,
				Submitter:          feed.Title,
				Posted_at:          item.Date,
			}
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func urlIsEmbeddable(url string) bool {
	return strings.Contains(url, "soundcloud.com") ||
		strings.Contains(url, "youtube.com") ||
		strings.Contains(url, "youtu.be")
}

type Parser struct {
}
