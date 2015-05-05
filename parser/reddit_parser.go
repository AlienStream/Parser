package parser

import (
	"encoding/json"
	"fmt"
	models "github.com/AlienStream/Shared-Go/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

////////////////////////////////
// REDDIT INTERMEDIATE OBJECT //
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

func getRedditSubredditData(source_data DataObject) DataObject {
	info := getRawSubredditMeta(source_data.Source.Url)
	source_data.Source.Title = info.Title
	source_data.Source.Description = info.Description
	source_data.Source.Thumbnail = info.Thumbnail

	// TODO: Multithread this into a queueable worker that respects the reddit limits
	raw_posts := getRawSubredditPosts(source_data.Source.Url, "sort=hot")
	raw_posts = append(raw_posts, getRawSubredditPosts(source_data.Source.Url+"/top/", "sort=top&t=day")...)
	raw_posts = append(raw_posts, getRawSubredditPosts(source_data.Source.Url+"/top/", "sort=top&t=week")...)
	raw_posts = append(raw_posts, getRawSubredditPosts(source_data.Source.Url+"/top/", "sort=top&t=month")...)
	raw_posts = append(raw_posts, getRawSubredditPosts(source_data.Source.Url+"/top/", "sort=top&t=year")...)
	raw_posts = append(raw_posts, getRawSubredditPosts(source_data.Source.Url+"/top/", "sort=top&t=all")...)

	for _, raw_post := range raw_posts {
		if strings.Contains(raw_post.Url, "soundcloud.com") || strings.Contains(raw_post.Url, "youtube.com") || strings.Contains(raw_post.Url, "youtu.b e") {

			post := models.Post{
				Id:                 0,
				Source_id:          source_data.Source.Id,
				Title:              raw_post.Title,
				Number_of_comments: raw_post.Num_Comments,
				Permalink:          raw_post.Permalink,
				Thumbnail:          raw_post.Thumbnail,
				Embed_url:          raw_post.Url,
				Likes:              raw_post.Upvotes,
				Dislikes:           raw_post.Downvotes,
				Submitter:          raw_post.Submitted_by,
				Posted_at:          time.Unix(int64(raw_post.Submitted_time), 0),
			}
			source_data.Posts = append(source_data.Posts, post)
		}
	}

	return source_data
}

func getRawSubredditPosts(source_url string, sort string) []redditPost {
	client := &http.Client{}

	// TODO: Paginated Results Gaunteeing at least 200 playable tracks
	target_url := fmt.Sprintf("%s.json?%s&limit=1000", source_url, sort)
	req, _ := http.NewRequest("GET", target_url, nil)
	req.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")

	resp, request_err := client.Do(req)
	defer resp.Body.Close()
	if request_err != nil {
		panic(request_err)
	}

	var data RedditRoot
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	decode_err := decoder.Decode(&data)

	if decode_err != nil {
		panic(decode_err)
	}

	var posts []redditPost
	for _, child := range data.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts
}

func getRawSubredditMeta(source_url string) SubredditInfo {
	client := &http.Client{}
	request_url := fmt.Sprintf("%s/about.json", source_url)
	request, _ := http.NewRequest("GET", request_url, nil)
	request.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")

	response, request_err := client.Do(request)
	if request_err != nil {
		panic(request_err)
	}

	defer response.Body.Close()
	data := SubredditAbout{}
	temp, _ := ioutil.ReadAll(response.Body)

	parse_err := json.Unmarshal(temp, &data)

	if parse_err != nil {
		panic("Requester failed to fetch from source")
	}

	return data.Data
}
