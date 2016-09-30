package reddit_parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	models "github.com/AlienStream/Shared-Go/models"
)

type Parser struct {
}

func (Parser) UpdateSourceMetaData(source *models.Source) error {
	info, err := getRawSubredditMeta(source.Url)

	if err == nil {
		source.Title = info.Title
		source.Description = info.Description
		source.Thumbnail = info.Thumbnail
	}

	return err
}

func (Parser) FetchPostsFromSource(source models.Source) ([]models.Post, error) {
	var posts []models.Post = []models.Post{}
	var error_occured error

	var request_urls = make(map[string]string)
	request_urls["sort=hot"] = source.Url
	request_urls["sort=top&t=day"] = source.Url + "/top/"
	request_urls["sort=top&t=week"] = source.Url + "/top/"
	request_urls["sort=top&t=month"] = source.Url + "/top/"
	request_urls["sort=top&t=year"] = source.Url + "/top/"
	request_urls["sort=top&t=all"] = source.Url + "/top/"

	for sort, url := range request_urls {
		var raw_posts []redditPost
		// TODO: Put this request onto a worker that respects the reddit limits
		raw_posts, err := getRawSubredditPosts(url, sort)
		if err != nil {
			error_occured = err
		}

		for _, raw_post := range raw_posts {
			if postIsEmbeddable(&raw_post) {
				post := models.Post{
					Id:                 0,
					Source_id:          source.Id,
					Title:              raw_post.Title,
					Number_of_comments: raw_post.Num_Comments,
					Permalink:          "https://reddit.com" + raw_post.Permalink,
					Thumbnail:          raw_post.Thumbnail,
					Embed_url:          raw_post.Url,
					Likes:              raw_post.Upvotes,
					Dislikes:           raw_post.Downvotes,
					Submitter:          raw_post.Submitted_by,
					Posted_at:          time.Unix(int64(raw_post.Submitted_time), 0),
				}
				posts = append(posts, post)
			}
		}
	}

	return posts, error_occured
}

func postIsEmbeddable(raw_post *redditPost) bool {
	return strings.Contains(raw_post.Url, "soundcloud.com") ||
		strings.Contains(raw_post.Url, "youtube.com") ||
		strings.Contains(raw_post.Url, "youtu.be")
}

func getRawSubredditPosts(source_url string, sort string) ([]redditPost, error) {
	// Setup Our Default Return Value
	var posts []redditPost

	// Make Our Request
	// TODO: Paginated Results Guaranteeing at least 200 playable tracks
	request_url := fmt.Sprintf("%s.json?%s&limit=1000", source_url, sort)
	req, _ := http.NewRequest("GET", request_url, nil)
	req.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")
	resp, request_err := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	if request_err != nil {
		return posts, request_err
	}

	var data RedditRoot
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	decode_err := decoder.Decode(&data)

	if decode_err != nil {
		return posts, decode_err
	}

	// Pull Out of Reddit Structure
	for _, child := range data.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts, nil
}

func getRawSubredditMeta(source_url string) (SubredditInfo, error) {
	// Setup Our Default Return Value
	data := SubredditAbout{
		Data: SubredditInfo{},
	}

	// Make our Request
	request_url := fmt.Sprintf("%s/about.json", source_url)
	request, _ := http.NewRequest("GET", request_url, nil)
	request.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")
	response, request_err := (&http.Client{}).Do(request)
	defer response.Body.Close()
	if request_err != nil {
		return data.Data, request_err
	}

	// Parse our Data
	temp, _ := ioutil.ReadAll(response.Body)
	parse_err := json.Unmarshal(temp, &data)
	if parse_err != nil {
		return data.Data, parse_err
	}

	// return just the info we need, reddit data nesting is dumb
	return data.Data, nil
}

////////////////////////////////
// REDDIT JSON OBJECTS        //
////////////////////////////////
type RedditRoot struct {
	Kind string     `json:"kind"`
	Data RedditData `json:"data"`
}

type RedditData struct {
	Children []RedditDataChild `json:"children"`
}

type RedditDataChild struct {
	Data redditPost `json:"data"`
}

type redditPost struct {
	Url            string  `json:"url"`
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	Thumbnail      string  `json:"thumbnail"`
	Submitted_by   string  `json:"author"`
	Submitted_time float64 `json:"created_utc"`
	Upvotes        int     `json:"ups"`
	Downvotes      int     `json:"downs"`
	Num_Comments   int     `json:"num_comments"`
	Permalink      string  `json:"permalink"`
}

type SubredditAbout struct {
	Data SubredditInfo `json:"data"`
}

type SubredditInfo struct {
	Title       string `json:"title"`
	Id          string `json:"name"`
	Thumbnail   string `json:"header_img"`
	Description string `json:"header_title"`
	Info        string `json:"description_html"`
	Subscribers int    `json:"subscribers"`
}
