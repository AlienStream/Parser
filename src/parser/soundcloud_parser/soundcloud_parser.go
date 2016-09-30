package soundcloud_parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	models "github.com/AlienStream/Shared-Go/models"
)

const SOUNDCLOUD_CLIENT_ID = "ff43d208510d35ce49ed972b01f116ab"

func (Parser) UpdateSourceMetaData(source *models.Source) error {
	info, err := getRawSoundcloudChannelMeta(source.Url)

	if err == nil {
		source.Title = info.Username
		source.Description = info.Description
		source.Thumbnail = info.Thumbnail
	}

	return err
}

func (Parser) FetchPostsFromSource(source models.Source) ([]models.Post, error) {
	posts := []models.Post{}

	// TODO: Multithread this into a queueable worker that respects the soundcloud limits
	raw_posts, err := getRawSoundcloudTracks(source.Url)
	if err != nil {
		return posts, err
	}
	for _, raw_post := range raw_posts {
		post := models.Post{
			Id:                 0,
			Source_id:          source.Id,
			Title:              raw_post.Title,
			Number_of_comments: raw_post.Comments,
			Permalink:          raw_post.Permalink,
			Thumbnail:          raw_post.Thumbnail,
			Embed_url:          raw_post.Permalink,
			Likes:              raw_post.Favorites,
			Dislikes:           0,
			Submitter:          raw_post.User.Username,
		}
		post.Posted_at, _ = time.Parse("2006/01/02 15:04:05 -0700", raw_post.Created_at)

		posts = append(posts, post)
	}

	return posts, nil
}

func getRawSoundcloudChannelMeta(source_url string) (SoundcloudUser, error) {
	data := SoundcloudUser{}

	request_url := strings.Replace(source_url, "://soundcloud.com/", "://api.soundcloud.com/users/", -1)
	request_url += "?client_id=" + SOUNDCLOUD_CLIENT_ID

	// Make Our Request
	request, _ := http.NewRequest("GET", request_url, nil)
	request.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")
	response, request_err := (&http.Client{}).Do(request)
	if request_err != nil {
		return data, request_err
	}

	if response.StatusCode != 200 {
		return data, errors.New("Response Was Not 200, Instead found:" + (string)(response.StatusCode))
	}

	// Parse the Response
	defer response.Body.Close()
	temp, _ := ioutil.ReadAll(response.Body)
	parse_err := json.Unmarshal(temp, &data)
	if parse_err != nil {
		return data, parse_err
	}

	return data, nil
}

func getRawSoundcloudTracks(source_url string) ([]SoundcloudTrack, error) {
	data := []SoundcloudTrack{}

	request_url := strings.Replace(source_url, "://soundcloud.com/", "://api.soundcloud.com/users/", -1)
	request_url += "/tracks?client_id=" + SOUNDCLOUD_CLIENT_ID

	// Make Our Request
	request, _ := http.NewRequest("GET", request_url, nil)
	request.Header.Set("User-Agent", "AlienStream Master Server v. 2.0")
	response, request_err := (&http.Client{}).Do(request)
	if request_err != nil {
		return data, request_err
	}

	// Parse the response
	defer response.Body.Close()
	temp, _ := ioutil.ReadAll(response.Body)
	parse_err := json.Unmarshal(temp, &data)
	if parse_err != nil {
		return data, parse_err
	}

	return data, nil
}

type Parser struct {
}

type SoundcloudUser struct {
	Permalink   string `json:"permalink_url"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Thumbnail   string `json:"avatar_url"`
}

type SoundcloudTrack struct {
	Title        string         `json:"title"`
	Permalink    string         `json:"permalink_url"`
	Thumbnail    string         `json:"artwork_url"`
	Submitted_by string         `json:"username"`
	Created_at   string         `json:"created_at"`
	Comments     int            `json:"comment_count"`
	Favorites    int            `json:"favoritings_count"`
	User         SoundcloudUser `json:"user"`
}
